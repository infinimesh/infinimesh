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
	"log"
	"net"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/infinimesh/infinimesh/pkg/mqtt"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/slntopp/mqtt-go/packet"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/grpc"
)

func verifyBasicAuth(p *packet.ConnectControlPacket) (fingerprint []byte, err error) {
	if p.ConnectPayload.Password == "" {
		return nil, errors.New("Payload Password is Empty")
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
	client      registrypb.DevicesClient
	debug       bool

	deviceRegistryHost    string
	kafkaHost             string
	kafkaTopicTelemetry   string
	kafkaTopicBackChannel string
	tlsCertFile           string
	tlsKeyFile            string
	dbAddr                string

	ps *pubsub.PubSub
)

func init() {
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
	dbAddr = viper.GetString("DB_ADDR2")
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
				fmt.Println("Failed to unmarshal message from kafka", err)
			}
			topic := fqTopic(m.DeviceID, m.SubPath)
			fmt.Println("pub to topic", fqTopic(m.DeviceID, m.SubPath))
			ps.Pub(&m, topic)
		}
	}
}

func fqTopic(deviceID, subPath string) string {
	return "devices/" + deviceID + "/" + subPath
}

func main() {

	serverCert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		log.Println(err)
		return
	}

	conn, err = grpc.Dial(deviceRegistryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = registrypb.NewDevicesClient(conn)

	fmt.Printf("KAFKA HOST :%v\n", kafkaHost)
	conf := sarama.NewConfig()
	conf.ClientID = "mqtt-bridge"
	kafkaClient, err = sarama.NewClient([]string{kafkaHost}, conf)
	if err != nil {
		panic(err)
	}

	producer, err = sarama.NewAsyncProducerFromClient(kafkaClient)
	if err != nil {
		panic(err)
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
				fmt.Println("Handshake failed", err)
				continue
			}
		case <-time.After(timeout):
			LogErrorAndClose(conn, errors.New("Handshake failed due to timeout"))
			continue
		}

		p, err := packet.ReadPacket(conn, 0)
		if err != nil {
			LogErrorAndClose(conn, fmt.Errorf("Error while reading connect packet: %v\n", err))
			continue
		}
		if debug {
			fmt.Println("ControlPacket", p)
		}

		connectPacket, ok := p.(*packet.ConnectControlPacket)
		if !ok {
			LogErrorAndClose(conn, errors.New("Got wrong packet as first packet..need connect!"))
			continue
		}
		if debug {
			fmt.Println("ConnectPacket", p)
		}

		if len(conn.(*tls.Conn).ConnectionState().PeerCertificates) == 0 {
			LogErrorAndClose(conn, errors.New("No certificate given"))
			continue
		}

		rawcert := conn.(*tls.Conn).ConnectionState().PeerCertificates[0].Raw
		fingerprint := getFingerprint(rawcert)

		if debug {
			fmt.Println("Fingerprint", fingerprint)
		}

		possibleIDs, err := GetByFingerprintAndVerify(fingerprint, func(device *registrypb.Device) (bool) {
			if device.Enabled.Value {
				fmt.Println(device.Tags)
				return true
			} else {
				fmt.Printf("Failed to verify client as the device is not enabled. Device ID:%v\n", device.Id)
				return false
			}
		})
		if err != nil {
			LogErrorAndClose(conn, err)
			continue
		}

		fmt.Printf("Client connected, IDs: %v\n", possibleIDs)

		go HandleConn(conn, connectPacket, possibleIDs)
	}

}

/*
func changeDeviceStatus(deviceID string, deviceStatus bool) {
	deviceStatusMap[deviceID] = deviceStatus
}*/

func handleBackChannel(c net.Conn, deviceID string, backChannel chan interface{}, protocolLevel byte) {
	// Everything from this channel is "vetted", i.e. it's legit that this client is subscribed to the topic.
	for message := range backChannel {
		fmt.Printf("Inside reading backchannel")
		m := message.(*mqtt.OutgoingMessage)
		// TODO PacketID
		fmt.Printf("m.subpath : %v", m.SubPath)
		topic := fqTopic(m.DeviceID, m.SubPath)
		fmt.Println("Publish to topic ", topic, "of client", deviceID)
		p := packet.NewPublish(topic, uint16(0), m.Data, protocolLevel)
		_, err := p.WriteTo(c)
		if err != nil {
			panic(err)
		}
	}

}

func printConnState(con net.Conn) {
	conn := con.(*tls.Conn)
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	state := conn.ConnectionState()
	log.Printf("Version: %x", state.Version)
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("DidResume: %t", state.DidResume)
	log.Printf("CipherSuite: %x", state.CipherSuite)
	log.Printf("NegotiatedProtocol: %s", state.NegotiatedProtocol)
	log.Printf("NegotiatedProtocolIsMutual: %t", state.NegotiatedProtocolIsMutual)

	log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string, topicAliasPublishMap map[string]int, protocolLevel int) (map[string]int, error) {
	fmt.Println("Handle publish", deviceID, p.VariableHeader.Topic, string(p.Payload))
	topic := TopicChecker(p.VariableHeader.Topic, deviceID)
	if p.VariableHeader.PublishProperties.TopicAlias > 0 {
		if val, ok := topicAliasPublishMap[topic]; ok {
			if val == p.VariableHeader.PublishProperties.TopicAlias {
				if err := publishTelemetry(topic, p.Payload, deviceID, protocolLevel); err != nil {
					return topicAliasPublishMap, err
				}
			} else {
				fmt.Printf("Please use the correct topic alias")
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
		fmt.Println("Payload schema invalid")
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
		log.Printf("invalid payload format")
		return false
	}
	loader := gojsonschema.NewGoLoader(payload)
	//filename := "file:///mqtt-bridge/schema-mqtt5.json"
	//log.Printf("json file path: %v", filename)
	schemaLoader := gojsonschema.NewStringLoader(mqtt5Schema)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Printf("Loading new schema failed %v", err)
		return false
	}
	result, err := schema.Validate(loader)
	if err != nil {
		log.Printf("Schema validation failed %v", err)
		return false
	}
	return result.Valid()
}
