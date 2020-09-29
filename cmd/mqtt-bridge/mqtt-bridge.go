//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
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
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/mqtt"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/mqtt-go/packet"
)

var verify = func(rawcerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	for _, rawcert := range rawcerts {
		digest := getFingerprint(rawcert)
		fmt.Printf("Validating certificate with fingerprint sha256-%X\n", digest)

		// Request information about all devices with this fingerprint
		reply, err := client.GetByFingerprint(context.Background(), &registrypb.GetByFingerprintRequest{Fingerprint: digest})
		if err != nil {
			fmt.Printf("Failed to find device for fingerprint: %v\n", err)
			continue
		}

		var enabled []*registrypb.Device
		for _, device := range reply.Devices {
			if device.Enabled.Value {
				enabled = append(enabled, device)
			}
		}

		if len(enabled) == 0 {
			return fmt.Errorf("no devices found for fingerprint %X\n", digest)
		}

		fmt.Printf("Verified connection with fingerprint [%v]. There are %v enabled devices with this fingerprint.\n", digest, len(enabled))
		return nil
	}
	return errors.New("Could not verify fingerprint")
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

	ps *pubsub.PubSub
)

func init() {
	viper.SetDefault("DEVICE_REGISTRY_URL", "localhost:8080")
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "mqtt.messages.incoming")
	viper.SetDefault("KAFKA_TOPIC_BACK", "mqtt.messages.outgoing")
	viper.SetDefault("TLS_CERT_FILE", "/cert/tls.crt")
	viper.SetDefault("TLS_KEY_FILE", "/cert/tls.key")
	viper.AutomaticEnv()

	deviceRegistryHost = viper.GetString("DEVICE_REGISTRY_URL")
	kafkaHost = viper.GetString("KAFKA_HOST")
	kafkaTopicTelemetry = viper.GetString("KAFKA_TOPIC")
	kafkaTopicBackChannel = viper.GetString("KAFKA_TOPIC_BACK")
	tlsCertFile = viper.GetString("TLS_CERT_FILE")
	tlsKeyFile = viper.GetString("TLS_KEY_FILE")

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

	conf := sarama.NewConfig()
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
		VerifyPeerCertificate: verify,
		ClientAuth:            tls.RequireAnyClientCert, // Any Client Cert is OK in terms of what the go TLS package checks, further validation, e.g. if the cert belongs to a registered device, is performed in the VerifyPeerCertificate function
	})
	if err != nil {
		panic(err)
	}

	go readBackchannelFromKafka()
	for {
		conn, _ := tlsl.Accept() // nolint: gosec
		timeout := time.Second / 2
		errChannel := make(chan error, 2)
		go func() {
			errChannel <- conn.(*tls.Conn).Handshake()
		}()
		select {
		case err := <-errChannel:
			if err != nil {
				fmt.Println("Handshake failed", err)
			}
		case <-time.After(timeout):
			fmt.Println("Handshake failed due to timeout")
			_ = conn.Close()
		}
		if len(conn.(*tls.Conn).ConnectionState().PeerCertificates) == 0 {
			continue
		}
		rawcert := conn.(*tls.Conn).ConnectionState().PeerCertificates[0].Raw
		reply, err := client.GetByFingerprint(context.Background(), &registrypb.GetByFingerprintRequest{
			Fingerprint: getFingerprint(rawcert),
		})
		if err != nil || len(reply.Devices) == 0 { //FIXME change logic so the client can send his id, and we track here which IDs are possible, but he can choose which identity he wants to use (in most cases it's only once, unless a device has multiple certs from multiple devices)
			_ = conn.Close()
			fmt.Printf("Failed to verify client, closing connection. err=%v\n", err)
			continue
		}

		var possibleIDs []string

		for _, device := range reply.Devices {
			if device.Enabled.Value {
				possibleIDs = append(possibleIDs, device.Id)
			}
		}

		fmt.Printf("Client connected, IDs: %v\n", possibleIDs)

		go handleConn(conn, possibleIDs)

	}

}

func handleBackChannel(c net.Conn, deviceID string, backChannel chan interface{}, protocolLevel byte) {
	// Everything from this channel is "vetted", i.e. it's legit that this client is subscribed to the topic.
	for message := range backChannel {
		m := message.(*mqtt.OutgoingMessage)
		// TODO PacketID
		fmt.Printf("m.subpath %v :", m.SubPath)
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

// Connection is expected to be valid & legitimate at this point
func handleConn(c net.Conn, deviceIDs []string) {
	p, err := packet.ReadPacket(c, 0)

	if debug {
		printConnState(c)
	}

	if err != nil {
		fmt.Printf("Error while reading connect packet: %v\n", err)
		return
	}

	connectPacket, ok := p.(*packet.ConnectControlPacket)
	if !ok {
		fmt.Println("Got wrong packet as first packet..need connect!")
		return
	}

	defer fmt.Println("Client disconnected ", connectPacket.ConnectPayload.ClientID)

	fmt.Printf("Client with ID %v connected!\n", connectPacket.ConnectPayload.ClientID)
	// TODO ignore/compare this ID with the given ID from the verify function

	var clientIDOK bool
	var deviceID string
	if len(deviceIDs) == 1 {
		// only one ID is possible with this cert; no need to have clientID set
		clientIDOK = true
		deviceID = deviceIDs[0]
	} else {
		fmt.Printf("Client used duplicate fingerprint, Please use unique certificate for your device\n")
		_ = c.Close()
		return
		//TODO : when multiple devices have single fingerprint authentication
		/*
			for _, possibleID := range deviceIDs {
				if connectPacket.ConnectPayload.ClientID == possibleID {
					fmt.Printf("Using ClientID: %v\n", possibleID)
					clientIDOK = true
					deviceID = possibleID
				}
			}
		*/
	}

	if !clientIDOK {
		fmt.Printf("Client used invalid clientID, disconnecting\n")
		_ = c.Close()
		return

	}
	//TODO : MQTT CONNACK Properties need to add here
	resp := packet.ConnAckControlPacket{
		FixedHeader: packet.FixedHeader{
			ControlPacketType: packet.CONNACK,
		},
		VariableHeader: packet.ConnAckVariableHeader{},
	}

	// Only open Back-channel after conn packet was received

	// Create empty subscription
	backChannel := ps.Sub()
	go handleBackChannel(c, deviceID, backChannel, connectPacket.VariableHeader.ProtocolLevel)
	defer func() {
		fmt.Printf("Unsubbed channel %v\n", deviceID)
		ps.Unsub(backChannel)
	}()

	_, err = resp.WriteTo(c)
	if err != nil {
		fmt.Println("Failed to write ConnAck. Closing connection.")
		return
	}
	var topicAliasPublishMap map[string]int
	topicAliasPublishMap = make(map[string]int)

	for {
		p, err := packet.ReadPacket(c, connectPacket.VariableHeader.ProtocolLevel)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client closed connection.\n")
			} else {
				fmt.Printf("Error while reading packet in client loop: %v. Disconnecting client.\n", err)
			}
			_ = c.Close() // nolint: gosec
			break
		}

		switch p := p.(type) {
		case *packet.PingReqControlPacket:
			pong := packet.NewPingRespControlPacket()
			_, err := pong.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write PingResp", err)
			}
		case *packet.PublishControlPacket:
			topicAliasPublishMap, err = handlePublish(p, c, deviceID, topicAliasPublishMap)
			if err != nil {
				fmt.Printf("Failed to handle Publish packet: %v.", err)
			}
		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write SubAck:", err)
			}
			for _, sub := range p.Payload.Subscriptions {
				ps.AddSub(backChannel, sub.Topic)
				go handleBackChannel(c, deviceID, backChannel, connectPacket.VariableHeader.ProtocolLevel)
				fmt.Println("Added Subscription", sub.Topic, deviceID)
			}
		case *packet.UnsubscribeControlPacket:
			response := packet.NewUnSubAck(uint16(p.VariableHeader.PacketID), connectPacket.VariableHeader.ProtocolLevel, []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write UnSubAck:", err)
			}
			for _, unsub := range p.Payload.UnSubscriptions {
				ps.Unsub(backChannel, unsub.Topic)
				fmt.Println("Removed Subscription", unsub.Topic, deviceID)
			}
		}
	}
}

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string, topicAliasPublishMap map[string]int) (map[string]int, error) {
	fmt.Println("Handle publish", deviceID, p.VariableHeader.Topic, string(p.Payload))
	if p.VariableHeader.PublishProperties.TopicAlias > 0 {
		if val, ok := topicAliasPublishMap[p.VariableHeader.Topic]; ok {
			if val == p.VariableHeader.PublishProperties.TopicAlias {
				if err := publishTelemetry(p.VariableHeader.Topic, p.Payload, deviceID); err != nil {
					return topicAliasPublishMap, err
				}
			} else {
				fmt.Printf("Please use the correct topic alias")
			}
		} else {
			topicAliasPublishMap[p.VariableHeader.Topic] = p.VariableHeader.PublishProperties.TopicAlias
			if err := publishTelemetry(p.VariableHeader.Topic, p.Payload, deviceID); err != nil {
				return topicAliasPublishMap, err
			}
		}
	} else {
		if err := publishTelemetry(p.VariableHeader.Topic, p.Payload, deviceID); err != nil {
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

func publishTelemetry(topic string, data []byte, deviceID string) error {
	message := mqtt.IncomingMessage{
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
	return nil
}
