package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"flag"

	"github.com/yosssi/gmq/mqtt/client"
)

var (
	topic  string
	broker string
)

func init() {
	flag.StringVar(&topic, "topic", "/", "MQTT Topic name")
	flag.StringVar(&broker, "broker", "localhost:8089", "MQTT Broker port. Defaults to localhost:8089")

}

func main() {
	flag.Parse()

	fmt.Println(os.Args)
	// Create an MQTT Client.
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println("ERR", err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	b, err := ioutil.ReadFile("hack/server.crt")
	if err != nil {
		panic(err)
	}

	kp, err := tls.LoadX509KeyPair("hack/server.crt", "hack/server.key")
	if err != nil {
		log.Println(err)
		return
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(b); !ok {
		panic("failed to parse root certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{kp},
		RootCAs:      roots,
	}

	// Connect to the MQTT Server using TLS.
	err = cli.Connect(&client.ConnectOptions{
		// Network is the network on which the Client connects to.
		ClientID: []byte("test"),
		Network:  "tcp",
		// Address is the address which the Client connects to.
		Address: broker,
		// TLSConfig is the configuration for the TLS connection.
		// If this property is not nil, the Client tries to use TLS
		// for the connection.
		TLSConfig: tlsConfig,
	})
	if err != nil {
		panic(err)
	}

	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{&client.SubReq{
			Handler: func(topicName, message []byte) {
				fmt.Println("recv", string(topicName), string(message))
			},
			TopicFilter: []byte("/updates"),
			QoS:         byte(0),
		}},
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 100)

	// if err := cli.Disconnect(); err != nil {
	// 	panic(err)
	// }

}
