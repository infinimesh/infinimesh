package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"crypto/sha256"

	"strings"

	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/infinimesh/infinimesh/pkg/registry"
	"github.com/infinimesh/mqtt-go/packet"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var verify = func(rawcerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	for _, rawcert := range rawcerts {
		digest := getFingerprint(rawcert)
		fmt.Printf("Validating certificate with fingerprint sha256-%X\n", digest)

		// Request information about a potential device with this fingerprint
		reply, err := client.GetByFingerprint(context.Background(), &registry.GetByFingerprintRequest{Fingerprint: digest})
		if err != nil {
			fmt.Printf("Failed to find device for fingerprint")
			continue
		}

		// TODO verify with configured CACert.
		fmt.Printf("Verified client with fingerprint device for fingerprint: %v\n", reply.Name)
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
	client      registry.DevicesClient
	debug       bool

	deviceRegistryHost    string
	kafkaHost             string
	kafkaTopicTelemetry   string
	kafkaTopicBackChannel string
)

type Message struct {
	Topic string
	Data  []byte
}

func init() {
	viper.SetDefault("DEVICE_REGISTRY_URL", "localhost:8080")
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "public.bridge.mqtt")
	viper.SetDefault("KAFKA_TOPIC_BACK", "public.bridge.mqtt.back-channel")
	viper.AutomaticEnv()

	deviceRegistryHost = viper.GetString("DEVICE_REGISTRY_URL")
	kafkaHost = viper.GetString("KAFKA_HOST")
	kafkaTopicTelemetry = viper.GetString("KAFKA_TOPIC")
	kafkaTopicBackChannel = viper.GetString("KAFKA_TOPIC_BACK")

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
			var m Message
			err = json.Unmarshal(message.Value, &m)
			if err != nil {
				fmt.Println("Failed to unmarshal message from kafka", err)
			}

			ps.Pub(&m, m.Topic)
		}
	}
}

var ps *pubsub.PubSub

func main() {
	serverCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	conn, err = grpc.Dial(deviceRegistryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = registry.NewDevicesClient(conn)

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
		reply, err := client.GetByFingerprint(context.Background(), &registry.GetByFingerprintRequest{
			Fingerprint: getFingerprint(rawcert),
		})
		if err != nil {
			_ = conn.Close()
			fmt.Printf("Failed to verify client, closing connection. err=%v\n", err)
			continue
		}
		fmt.Printf("Client connected, device name according to registry: %v\n", reply.GetName())

		backChannel := ps.Sub()

		go handleConn(conn, reply.GetName(), backChannel)
		go handleBackChannel(conn, reply.GetName(), backChannel)
	}

}

var (

// TODO maybe better use a publish/subscribe package for this instead of
// doing this ourselves
// mtx               sync.Mutex
// subscribedClients map[string]map[string]chan *msg
)

func init() {
	// subscribedClients = make(map[string]map[string]chan *msg)
}

// TODO need context struct for connection; context struct has []subscription-channels, select from those?
func handleBackChannel(c net.Conn, deviceID string, backChannel chan interface{}) {
	for message := range backChannel {

		// Everything from this channel is "vetted", i.e. it's legit that this client is subscribed to the topic.
		m := message.(*Message)
		p := packet.NewPublish(m.Topic /* TODO */, uint16(0), m.Data)
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
		case *packet.PublishControlPacket:
			err = handlePublish(p, c, deviceID)
			if err != nil {
				fmt.Printf("Failed to handle Publish packet: %v.", err)
			}
		case *packet.SubscribeControlPacket:
			response := packet.NewSubAck(uint16(p.VariableHeader.PacketID), []byte{1})
			_, err := response.WriteTo(c)
			if err != nil {
				panic(err)
			}

			// TODO better loop over subscribing topics..
			topic := p.Payload.Subscriptions[0].Topic

			ps.AddSub(backChannel, topic)
		}
	}
}

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string) error {
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

type MQTTBridgeData struct {
	SourceTopic  string
	SourceDevice string
	Data         []byte
}

func publishTelemetry(topic string, data []byte, deviceID string) error {
	var targetTopic string

	// Current allowed layout:
	// /devices/%v/... for any telemetry data
	// /shadows/%v/... for any shadows (json-diffs)

	fmt.Println("Try to send to topic", topic)
	// Currently, only the exact shadow & telemetry topic of the device is
	// allowed. Later, this will be dynamic, a device can have multiple
	// shadows or even write to other device's shadows + telemetry topics

	// TODO specific data types! shadow = text or json only; telemetry is raw binary.

	// Anything below the shadow topic or device topic of the device is allowed
	if strings.HasPrefix(topic, fmt.Sprintf("/shadows/%v", deviceID)) {
		targetTopic = "public.delta.reported-state" // shadows
	} else if strings.HasPrefix(topic, fmt.Sprintf("/devices/%v", deviceID)) {
		targetTopic = kafkaTopicTelemetry // public.bridge.mqtt - "raw" telemetry
	} else {
		// Dead letter queue - we didn't have permission,...
		targetTopic = "public.bridge.dlq"
		// TODO maybe write a reason
	}

	message := MQTTBridgeData{
		SourceTopic:  topic,
		SourceDevice: deviceID,
		Data:         data,
	}

	serialized, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: targetTopic,
		Key:   sarama.StringEncoder(deviceID), // TODO
		Value: sarama.ByteEncoder(serialized),
	}
	return nil
}
