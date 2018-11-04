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

	"github.com/Shopify/sarama"
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

	deviceRegistryHost string
	kafkaHost          string
	kafkaTopic         string
)

func init() {
	viper.SetDefault("DEVICE_REGISTRY_URL", "localhost:8080")
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "public.bridge.mqtt")
	viper.AutomaticEnv()

	deviceRegistryHost = viper.GetString("DEVICE_REGISTRY_URL")
	kafkaHost = viper.GetString("KAFKA_HOST")
	kafkaTopic = viper.GetString("KAFKA_TOPIC")

}

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

	tlsl, err := tls.Listen("tcp", ":8089", &tls.Config{
		Certificates:          []tls.Certificate{serverCert},
		VerifyPeerCertificate: verify,
		ClientAuth:            tls.RequireAnyClientCert, // Any Client Cert is OK in terms of what the go TLS package checks, further validation, e.g. if the cert belongs to a registered device, is performed in the VerifyPeerCertificate function
	})

	if err != nil {
		panic(err)
	}

	for {
		conn, _ := tlsl.Accept() // nolint: gosec
		err := conn.(*tls.Conn).Handshake()
		if err != nil {
			fmt.Println("Handshake of client failed", err)
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
		go handleConn(conn, reply.GetName())
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
// TODO maybe provide a connection-context struct with metadata about this client, e.g. the associated device, ..
func handleConn(c net.Conn, deviceID string) {
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
		fmt.Println("Got wrong packet as first packjet..need connect!")
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
		}
	}
}

func handlePublish(p *packet.PublishControlPacket, c net.Conn, deviceID string) error {
	if err := publishTelemetry(p.VariableHeader.Topic, p.Payload, deviceID); err != nil {
		return err
	}
	if p.FixedHeaderFlags.QoS >= packet.QoSLevelAtLeastOnce {
		pubAck := packet.NewPubAckControlPacket(uint16(p.VariableHeader.PacketID))
		_, err := pubAck.WriteTo(c)
		if err != nil {
			return err
		}

	}
	return nil
}

func publishTelemetry(topic string, data []byte, deviceID string) error {
	var targetTopic string

	fmt.Println("Send to topic", topic)
	if topic == "_shadow" {
		targetTopic = "public.delta.reported-state"
	} else {
		targetTopic = kafkaTopic
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: targetTopic,
		Key:   sarama.StringEncoder(deviceID), // TODO
		Value: sarama.ByteEncoder(data),
	}
	return nil
}
