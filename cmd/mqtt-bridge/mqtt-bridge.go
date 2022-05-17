//--------------------------------------------------------------------------
// Copyright 2018-2022 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cskr/pubsub"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/infinimesh/pkg/mqtt"
	mqttps "github.com/infinimesh/infinimesh/pkg/mqtt/pubsub"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	pb "github.com/infinimesh/proto/node"
	devpb "github.com/infinimesh/proto/node/devices"
	stpb "github.com/infinimesh/proto/shadow"
	"github.com/slntopp/mqtt-go/packet"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	inflog "github.com/infinimesh/infinimesh/pkg/log"
)

func verifyBasicAuth(p *packet.ConnectControlPacket) (fingerprint []byte, err error) {
	if p.ConnectPayload.Password == "" {
		return nil, errors.New("payload Password is Empty")
	}
	return base64.StdEncoding.DecodeString(p.ConnectPayload.Password)
}

func getFingerprint(c []byte) []byte {
	s := sha256.New()
	_, _ = s.Write(c) // nolint: gosec
	return s.Sum(nil)
}

var (
	conn   *grpc.ClientConn
	client pb.DevicesServiceClient
	shadow stpb.ShadowServiceClient
	debug  bool

	devicesHost string
	shadowHost  string

	RabbitMQConn string
	tlsCertFile  string
	tlsKeyFile   string

	ps *pubsub.PubSub

	log          *zap.Logger
	internal_ctx context.Context
)

func init() {
	var err error
	log, err = inflog.NewProdOrDev()
	if err != nil {
		panic(err)
	}

	viper.AutomaticEnv()

	viper.SetDefault("DEVICES_HOST", "api.infinimesh.local")
	viper.SetDefault("SHADOW_HOST", "shadow:8080")
	viper.SetDefault("DB_ADDR2", ":6379")
	viper.SetDefault("RABBITMQ_CONN", "amqp://infinimesh:infinimesh@localhost:5672/")
	viper.SetDefault("TLS_CERT_FILE", "/cert/tls.crt")
	viper.SetDefault("TLS_KEY_FILE", "/cert/tls.key")
	viper.SetDefault("DEBUG", false)
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	devicesHost = viper.GetString("DEVICES_HOST")
	shadowHost = viper.GetString("SHADOW_HOST")
	RabbitMQConn = viper.GetString("RABBITMQ_CONN")
	tlsCertFile = viper.GetString("TLS_CERT_FILE")
	tlsKeyFile = viper.GetString("TLS_KEY_FILE")
	debug = viper.GetBool("DEBUG")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting MQTT Bridge")
	serverCert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		log.Fatal("Error loading server certificate", zap.Error(err))
	}

	log.Info("Connecting to registry", zap.String("host", devicesHost))
	conn, err = grpc.Dial(devicesHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error dialing device registry", zap.Error(err))
	}
	client = pb.NewDevicesServiceClient(conn)

	conn, err = grpc.Dial(shadowHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Error dialing shadow registry", zap.Error(err))
	} else {
		shadow = stpb.NewShadowServiceClient(conn)
	}

	SIGNING_KEY := []byte(viper.GetString("SIGNING_KEY"))
	auth.SetContext(log, SIGNING_KEY)
	token, err := auth.MakeToken(schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Error making token", zap.Error(err))
	}
	internal_ctx = metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)

	log.Info("Connecting to RabbitMQ", zap.String("url", RabbitMQConn))
	rbmq, err := amqp.Dial(RabbitMQConn)
	if err != nil {
		log.Fatal("Error dialing RabbitMQ", zap.Error(err))
	}
	defer rbmq.Close()

	ps, err = mqttps.Setup(log, rbmq, "mqtt.incoming", "mqtt.outgoing")
	if err != nil {
		log.Fatal("Error setting up pubsub", zap.Error(err))
	}

	tlsl, err := tls.Listen("tcp", ":8883", &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAnyClientCert, // Any Client Cert is OK in terms of what the go TLS package checks, further validation, e.g. if the cert belongs to a registered device, is performed in the VerifyPeerCertificate function
	})
	if err != nil {
		panic(err)
	}
	tcp, err := net.Listen("tcp", ":1883")
	if err != nil {
		panic(err)
	}

	go HandleTCPConnections(tcp)
	for {
		conn, _ := tlsl.Accept() // nolint: gosec
		if debug {
			printConnState(conn)
		}
		timeout := time.Second * 30
		errChannel := make(chan error, 2)
		go func() {
			errChannel <- conn.(*tls.Conn).Handshake()
		}()
		select {
		case err := <-errChannel:
			if err != nil {
				log.Info("Handshake failed", zap.Error(err))
				continue
			}
		case <-time.After(timeout):
			LogErrorAndClose(conn, errors.New("handshake failed due to timeout"))
			continue
		}

		p, err := packet.ReadPacket(conn, 0)
		if err != nil {
			LogErrorAndClose(conn, fmt.Errorf("error while reading connect packet: %v", err))
			continue
		}

		log.Debug("Control packet", zap.Any("packet", p))

		connectPacket, ok := p.(*packet.ConnectControlPacket)
		if !ok {
			LogErrorAndClose(conn, errors.New("first packet isn't ConnectControlPacket"))
			continue
		}
		log.Debug("ConnectPacket", zap.Any("packet", p))

		if len(conn.(*tls.Conn).ConnectionState().PeerCertificates) == 0 {
			LogErrorAndClose(conn, errors.New("no certificate given"))
			continue
		}

		rawcert := conn.(*tls.Conn).ConnectionState().PeerCertificates[0].Raw
		fingerprint := getFingerprint(rawcert)
		log.Debug("Fingerprint", zap.ByteString("fingerprint", fingerprint))

		device, err := GetByFingerprintAndVerify(fingerprint, func(device *devpb.Device) bool {
			if device.Enabled {
				log.Info("Device is enabled", zap.String("device", device.Uuid), zap.Strings("tags", device.Tags))
				return true
			} else {
				log.Error("Failed to verify client as the device is not enabled", zap.String("device", device.Uuid))
				return false
			}
		})
		if err != nil {
			LogErrorAndClose(conn, err)
			continue
		}

		log.Info("Client connected", zap.String("device", device.Uuid))

		go HandleConn(conn, connectPacket, device)
	}

}

func printConnState(con net.Conn) {
	conn := con.(*tls.Conn)
	state := conn.ConnectionState()

	log.Info("Connection state",
		zap.Uint16("version", state.Version),
		zap.Bool("handshake-complete", state.HandshakeComplete),
		zap.Bool("did-resume", state.DidResume),
		zap.Uint16("cipher-suite", state.CipherSuite),
		zap.String("proto-version", state.NegotiatedProtocol),
		zap.Any("certs", state.PeerCertificates),
	)
}

/*TopicChecker: to validate the subscribed topic name
  input : topic, deviceId string
  output : topicAltered
*/
func TopicChecker(topic, deviceId string) string {
	state := strings.Split(topic, "/")
	state[1] = deviceId
	topic = strings.Join(state, "/")
	return topic
}

//MQTT5 schema
const mqtt5Schema = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "type": "object",
    "properties": {
      "Timestamp": {
        "type": "string"
      },
      "Message": {
        "type": "array",
        "items": [
          {
            "type": "object",
            "properties": {
              "Topic": {
                "type": "string"
              },
              "Data": {
                "type": "object"
              }
            },
            "required": [
              "Topic",
              "Data"
            ]
          }
        ]
      }
    },
    "required": [
      "Timestamp",
      "Message"
    ]
  }`

func schemaValidation(data []byte, version int) bool {
	if version == 4 {
		return true
	}
	var payload mqtt.Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		log.Error("invalid payload format", zap.Error(err))
		return false
	}
	loader := gojsonschema.NewGoLoader(payload)
	//filename := "file:///mqtt-bridge/schema-mqtt5.json"
	//log.Printf("json file path: %v", filename)
	schemaLoader := gojsonschema.NewStringLoader(mqtt5Schema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Error("Loading new schema failed", zap.Error(err))
		return false
	}
	result, err := schema.Validate(loader)
	if err != nil {
		log.Error("Schema validation failed", zap.Error(err))
		return false
	}
	return result.Valid()
}
