package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"flag"

	"github.com/yosssi/gmq/mqtt/client"
)

var (
	topic string
)

func init() {
	flag.StringVar(&topic, "topic", "/", "MQTT Topic name")
}

func main() {
	flag.Parse()
	// Create an MQTT Client.
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Read the certificate file.
	b, err := ioutil.ReadFile("../mqtt-bridge/server.crt")
	if err != nil {
		panic(err)
	}

	kp, err := tls.LoadX509KeyPair("../mqtt-bridge/server.crt", "../mqtt-bridge/server.key")
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
		Address: "localhost:8089",
		// TLSConfig is the configuration for the TLS connection.
		// If this property is not nil, the Client tries to use TLS
		// for the connection.
		TLSConfig: tlsConfig,
	})
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	_, _ = io.Copy(&buf, os.Stdin)

	err = cli.Publish(&client.PublishOptions{
		QoS:       byte(1),
		TopicName: []byte("_shadow"),
		Message:   buf.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Millisecond * 10)

	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
