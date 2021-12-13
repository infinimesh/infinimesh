//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/slntopp/infinimesh/pkg/avro"
	"github.com/slntopp/infinimesh/pkg/avro/avropb"
	log1 "github.com/slntopp/infinimesh/pkg/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	sourceTopicReported = "shadow.reported-state.full"
	sourceTopicDesired  = "shadow.desired-state.full"
	consumerGroup       = "avro-persister"
)

var (
	broker     string
	avroHost   string
	avroClient avropb.AvroreposClient
)

func init() {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("AVRO_HOST", "localhost:50054")
	viper.SetDefault("KAFKA_CONSUMER_GROUP", "avro-persister")
	viper.AutomaticEnv()
	broker = viper.GetString("KAFKA_HOST")
	avroHost = viper.GetString("AVRO_HOST")
	//consumerGroup = viper.GetString("KAFKA_CONSUMER_GROUP")
}

func main() {
	log1, err := log1.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log1.Sync()
	}()

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = false
	config.Version = sarama.V2_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10

	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		panic(err)
	}
	// gRPC client initialization
	conn, err := grpc.Dial(avroHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("unable to connect to localhost:50054")
	}
	defer conn.Close()

	avroClient = avropb.NewAvroreposClient(conn)

	group, err := sarama.NewConsumerGroupFromClient(consumerGroup, client)
	if err != nil {
		panic(err)
	}

	handler := &avro.Consumer{
		SourceTopicReported: sourceTopicReported,
		SourceTopicDesired:  sourceTopicDesired,
		ConsumerGroup:       consumerGroup,
		AvroClient:          avroClient,
		Log:                 log1.Named("Avro Connector Controller"),
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {
	outer:
		for {

			err = group.Consume(context.Background(), []string{sourceTopicDesired, sourceTopicReported}, handler)
			if err != nil {
				panic(err)
			}

			select {
			case <-done:
				break outer
			default:
			}

		}

	}()

	<-c
	done <- true
	err = group.Close()
	if err != nil {
		panic(err)
	}
}
