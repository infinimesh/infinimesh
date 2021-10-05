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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"github.com/Shopify/sarama"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"github.com/infinimesh/infinimesh/pkg/shadow"
)

const (
	sourceTopicReported = "shadow.reported-state.full"
	sourceTopicDesired  = "shadow.desired-state.full"
)

var (
	addr   string
	broker string

	consumerGroup string
)

func init() {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("DB_ADDR", ":6379")
	viper.SetDefault("KAFKA_CONSUMER_GROUP", "persister")
	viper.AutomaticEnv()
	broker = viper.GetString("KAFKA_HOST")
	consumerGroup = viper.GetString("KAFKA_CONSUMER_GROUP")
	addr = viper.GetString("DB_ADDR")
}

type handler struct {
	repo shadow.Repo
}

func main() {
	repo, err := shadow.NewRedisRepo(addr)
	if err != nil {
		log.Fatal(err)
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = false
	config.Version = sarama.V2_0_0_0

	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	group, err := sarama.NewConsumerGroupFromClient(consumerGroup, client)
	if err != nil {
		panic(err)
	}

	handler := &handler{repo: repo}

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

func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	return nil
}

func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		fmt.Println("got msg", string(message.Value))

		var stateFromKafka shadow.DeviceStateMessage
		if err := json.Unmarshal(message.Value, &stateFromKafka); err != nil {
			fmt.Println("Failed to deserialize message with offset ", message.Offset)
			continue
		}

		var dbErr error

		switch message.Topic {
		case sourceTopicReported:
			dbErr = h.repo.SetReported(shadow.DeviceState{
				ID:    string(message.Key),
				State: stateFromKafka,
			})
		case sourceTopicDesired:
			dbErr = h.repo.SetDesired(shadow.DeviceState{
				ID:    string(message.Key),
				State: stateFromKafka,
			})
		}

		// FIXME ignore errors for now
		if dbErr != nil {
			fmt.Println("Failed to persist message with offset", message.Offset, dbErr)
		}

		s.MarkMessage(message, "")
	}
	return nil
}
