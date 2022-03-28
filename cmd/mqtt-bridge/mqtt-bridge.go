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
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/infinimesh/infinimesh/pkg/mqtt"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"github.com/slntopp/mqtt-go/packet"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	conn        *grpc.ClientConn
	kafkaClient sarama.Client
	producer    sarama.AsyncProducer
	client      pb.DevicesServiceClient
	debug       bool

	deviceRegistryHost    string
	kafkaHost             string
	kafkaTopicTelemetry   string
	kafkaTopicBackChannel string
	tlsCertFile           string
	tlsKeyFile            string

	ps *pubsub.PubSub

	log *zap.Logger
)

func init() {
	var err error
	log, err = inflog.NewProdOrDev()
	if err != nil {
		panic(err)
	}

	viper.SetDefault("DEVICE_REGISTRY_URL", "localhost:8080")
	viper.SetDefault("DB_ADDR2", ":6379")
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "mqtt.messages.incoming")
	viper.SetDefault("KAFKA_TOPIC_BACK", "mqtt.messages.outgoing")
	viper.SetDefault("TLS_CERT_FILE", "/cert/tls.crt")
	viper.SetDefault("TLS_KEY_FILE", "/cert/tls.key")
	viper.SetDefault("DEBUG", false)
	viper.AutomaticEnv()

	deviceRegistryHost = viper.GetString("DEVICE_REGISTRY_URL")
	kafkaHost = viper.GetString("KAFKA_HOST")
	kafkaTopicTelemetry = viper.GetString("KAFKA_TOPIC")
	kafkaTopicBackChannel = viper.GetString("KAFKA_TOPIC_BACK")
	tlsCertFile = viper.GetString("TLS_CERT_FILE")
	tlsKeyFile = viper.GetString("TLS_KEY_FILE")
	debug = viper.GetBool("DEBUG")
}

func readBackchannelFromKafka() {
	consumer, err := sarama.NewConsumerFromClient(kafkaClient)
	if err != nil {
		panic(err)
	}

	partitions, err := consumer.Partitions(kafkaTopicBackChannel)
	if err != nil {
		panic(err)
	}
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(kafkaTopicBackChannel, partition, sarama.OffsetNewest) // TODO, currently no guarantees, just process new messages
		if err != nil {
			panic(err)
		}

		for message := range pc.Messages() {
			var m mqtt.OutgoingMessage
			err = json.Unmarshal(message.Value, &m)
			if err != nil {
				log.Error("Failed to unmarshal message from kafka", zap.Error(err))
			}
			topic := fqTopic(m.DeviceID, m.SubPath)
			log.Info("Publish to the topic", zap.String("topic", topic))
			ps.Pub(&m, topic)
		}
	}
}

func fqTopic(deviceID, subPath string) string {
	return "devices/" + deviceID + "/" + subPath
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

	log.Info("Connecting to registry", zap.String("host", deviceRegistryHost))
	conn, err = grpc.Dial(deviceRegistryHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error dialing device registry", zap.Error(err))
	}
	client = pb.NewDevicesServiceClient(conn)

	log.Info("Connecting to Kafka", zap.String("host", kafkaHost))
	conf := sarama.NewConfig()
	kafkaClient, err = sarama.NewClient([]string{kafkaHost}, conf)
	if err != nil {
		log.Fatal("Error creating kafka client", zap.Error(err))
	}

	producer, err = sarama.NewAsyncProducerFromClient(kafkaClient)
	if err != nil {
		log.Fatal("Error creating kafka producer", zap.Error(err))
	}

	ps = pubsub.New(10)

	tlsl, err := tls.Listen("tcp", ":8089", &tls.Config{
		Certificates:          []tls.Certificate{serverCert},
		ClientAuth:            tls.RequireAnyClientCert, // Any Client Cert is OK in terms of what the go TLS package checks, further validation, e.g. if the cert belongs to a registered device, is performed in the VerifyPeerCertificate function
	})
	if err != nil {
		panic(err)
	}
	tcp, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go readBackchannelFromKafka()
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

		device, err := GetByFingerprintAndVerify(fingerprint, func(device *devpb.Device) (bool) {
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

/*
func changeDeviceStatus(deviceID string, deviceStatus bool) {
	deviceStatusMap[deviceID] = deviceStatus
}*/

func handleBackChannel(c net.Conn, deviceID string, backChannel chan interface{}, protocolLevel byte) {
	// Everything from this channel is "vetted", i.e. it's legit that this client is subscribed to the topic.
	for message := range backChannel {
		m := message.(*mqtt.OutgoingMessage)
		// TODO PacketID
		topic := fqTopic(m.DeviceID, m.SubPath)
		log.Info("Publish to the topic", zap.String("topic", topic), zap.String("client", deviceID))
		p := packet.NewPublish(topic, uint16(0), m.Data, protocolLevel)
		_, err := p.WriteTo(c)
		if err != nil {
			panic(err)
		}
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

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string, topicAliasPublishMap map[string]int, protocolLevel int) (map[string]int, error) {
	log.Info("Handle publish", zap.String("device", deviceID), zap.String("topic", p.VariableHeader.Topic), zap.ByteString("payload", p.Payload))
	topic := TopicChecker(p.VariableHeader.Topic, deviceID)
	if p.VariableHeader.PublishProperties.TopicAlias > 0 {
		if val, ok := topicAliasPublishMap[topic]; ok {
			if val == p.VariableHeader.PublishProperties.TopicAlias {
				if err := publishTelemetry(topic, p.Payload, deviceID, protocolLevel); err != nil {
					return topicAliasPublishMap, err
				}
			} else {
				log.Error("Please use correct topic alias")
			}
		} else {
			topicAliasPublishMap[topic] = p.VariableHeader.PublishProperties.TopicAlias
			if err := publishTelemetry(topic, p.Payload, deviceID, protocolLevel); err != nil {
				return topicAliasPublishMap, err
			}
		}
	} else {
		if err := publishTelemetry(topic, p.Payload, deviceID, protocolLevel); err != nil {
			return topicAliasPublishMap, err
		}
	}
	if p.FixedHeaderFlags.QoS >= packet.QoSLevelAtLeastOnce {
		pubAck := packet.NewPubAckControlPacket(uint16(p.VariableHeader.PacketID)) // TODO better always use directly uint16 for PacketIDs,everywhere
		_, err := pubAck.WriteTo(c)
		if err != nil {
			return topicAliasPublishMap, err
		}
	}
	return topicAliasPublishMap, nil
}

/*TopicChecker: to validate the subscribed topic name
  input : topic, deviceId string
  output : topicAltered
*/
func TopicChecker(topic, deviceId string) (string) {
	state := strings.Split(topic, "/")
	state[1] = deviceId
	topic = strings.Join(state, "/")
	return topic
}

func publishTelemetry(topic string, data []byte, deviceID string, version int) error {
	valid := schemaValidation(data, version)
	if valid {
		message := mqtt.IncomingMessage{
			ProtoLevel:   version,
			SourceTopic:  topic,
			SourceDevice: deviceID,
			Data:         data,
		}

		serialized, err := json.Marshal(&message)
		if err != nil {
			return err
		}
		producer.Input() <- &sarama.ProducerMessage{
			Topic: kafkaTopicTelemetry,
			Key:   sarama.StringEncoder(deviceID), // TODO
			Value: sarama.ByteEncoder(serialized),
		}
	} else {
		log.Error("Payload schema is invalid")
	}
	return nil
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
