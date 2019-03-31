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

		// Request information about a potential device with this fingerprint
		reply, err := client.GetByFingerprint(context.Background(), &registrypb.GetByFingerprintRequest{Fingerprint: digest})
		if err != nil {
			fmt.Printf("Failed to find device for fingerprint: %v\n", err)
			continue
		}

		if len(reply.Devices) == 0 {
			return errors.New("Invalid device")
		}

		result := reply.Devices[0] // FIXME silly Workaround

		// TODO verify with configured CACert.
		fmt.Printf("Verified client with fingerprint device for fingerprint: %v\n", result.Name)
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
	return "/devices/" + deviceID + "/" + subPath
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
		err := conn.(*tls.Conn).Handshake()
		if err != nil {
			fmt.Println("Handshake of client failed", err)
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
		response := reply.Devices[0] // FIXME
		fmt.Printf("Client connected, id: %v, name: %v\n", response.Id, response.Name)

		backChannel := ps.Sub()

		go handleConn(conn, response.Id, backChannel)
		go handleBackChannel(conn, response.Id, backChannel)
	}

}

func handleBackChannel(c net.Conn, deviceID string, backChannel chan interface{}) {
	// Everything from this channel is "vetted", i.e. it's legit that this client is subscribed to the topic.
	for message := range backChannel {
		m := message.(*mqtt.OutgoingMessage)
		// TODO PacketID
		topic := fqTopic(m.DeviceID, m.SubPath)
		fmt.Println("Publish to topic ", topic, "of client", deviceID)
		p := packet.NewPublish(topic, uint16(0), m.Data)
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
func handleConn(c net.Conn, deviceID string, backChannel chan interface{}) {
	defer fmt.Println("Client disconnected ", deviceID)
	defer func() {
		ps.Unsub(backChannel)
	}()
	p, err := packet.ReadPacket(c)

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

	id := connectPacket.ConnectPayload.ClientID
	fmt.Printf("Client with ID %v connected!\n", id)
	// TODO ignore/compare this ID with the given ID from the verify function

	resp := packet.ConnAckControlPacket{
		FixedHeader: packet.FixedHeader{
			ControlPacketType: packet.CONNACK,
		},
		VariableHeader: packet.ConnAckVariableHeader{},
	}

	_, err = resp.WriteTo(c)
	if err != nil {
		fmt.Println("Failed to write ConnAck. Closing connection.")
		return
	}

	for {
		p, err := packet.ReadPacket(c)
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
			err = handlePublish(p, c, deviceID)
			if err != nil {
				fmt.Printf("Failed to handle Publish packet: %v.", err)
			}
		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				fmt.Println("Failed to write SubAck:", err)
			}

			// TODO better loop over subscribing topics..
			topic := p.Payload.Subscriptions[0].Topic

			ps.AddSub(backChannel, topic)
			fmt.Println("Added Subscription", topic, deviceID)
		}
	}
}

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string) error {
	fmt.Println("Handle publish", deviceID, p.VariableHeader.Topic, string(p.Payload))
	if err := publishTelemetry(p.VariableHeader.Topic, p.Payload, deviceID); err != nil {
		return err
	}
	if p.FixedHeaderFlags.QoS >= packet.QoSLevelAtLeastOnce {
		pubAck := packet.NewPubAckControlPacket(uint16(p.VariableHeader.PacketID)) // TODO better always use directly uint16 for PacketIDs,everywhere
		_, err := pubAck.WriteTo(c)
		if err != nil {
			return err
		}

	}
	return nil
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
