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
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"

	"github.com/infinimesh/infinimesh/pkg/shadow"
)

var (
	consumerGroupReported           = "shadow-reported"
	consumerGroupDesired            = "shadow-desired"
	broker                          string
	topicReportedState              = "shadow.reported-state.delta"
	topicDesiredState               = "shadow.desired-state.delta"
	mergedTopicReported             = "shadow.reported-state.full"
	mergedTopicDesired              = "shadow.desired-state.full"
	topicComputedDeltaReportedState = "shadow.reported-state.delta.computed"
	topicComputedDeltaDesiredState  = "shadow.desired-state.delta.computed"
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
}

func runMerger(inputTopic, outputTopic, realDeltaTopic, consumerGroup string, stop chan bool, ctx context.Context) (close io.Closer, done chan bool) {
	done = make(chan bool)
	consumerGroupClient := sarama.NewConfig()
	consumerGroupClient.Version = sarama.V1_0_0_0
	consumerGroupClient.Consumer.Return.Errors = true
	consumerGroupClient.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient([]string{broker}, consumerGroupClient)
	fmt.Printf("client created %v\n", client)
	if err != nil {
		panic(err)
	}

	group, err := sarama.NewConsumerGroupFromClient(consumerGroup, client)
	fmt.Printf("consumer group %v, client %v, group %v \n", consumerGroup, client, group)
	if err != nil {
		panic(err)
	}

	go func() {
		for err := range group.Errors() {
			fmt.Printf("Consumer group error: %v\n", err)
		}

	}()

	pconfig := sarama.NewConfig()
	pconfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	pconfig.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	pconfig.Producer.Partitioner = sarama.NewManualPartitioner
	pconfig.Producer.Return.Errors = false
	pconfig.Producer.Return.Successes = false
	pconfig.Version = sarama.V1_1_0_0

	producerClient, err := sarama.NewClient([]string{broker}, pconfig)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	pconfig.Producer.Return.Errors = false
	pconfig.Producer.Return.Successes = false
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message

	localStateConsumerClient, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	handler := &shadow.StateMerger{
		SourceTopic:             inputTopic,
		MergedTopic:             outputTopic,
		RealDeltaTopic:          realDeltaTopic,
		ChangelogProducerClient: producerClient,
		ChangelogConsumerClient: localStateConsumerClient,
	}

	go func() {
	outer:
		for {
			fmt.Printf("group : %v, %v,  %v \n", group, inputTopic, handler)
			err = group.Consume(ctx, []string{inputTopic}, handler)
			fmt.Printf("Handler called : %v and err: %v ", inputTopic, err)
			if err != nil {
				panic(err)
			}

			select {
			case <-stop:
				done <- true
				break outer
			default:
			}

		}

	}()
	return group, done
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	stopReported := make(chan bool, 2)
	stopDesired := make(chan bool, 2)

	ctx, cancel := context.WithCancel(context.Background())

	closeReported, doneReported := runMerger(topicReportedState, mergedTopicReported, topicComputedDeltaReportedState, consumerGroupReported, stopReported, ctx)
	fmt.Printf("closeReported, doneReported: %v, %v\n", closeReported, doneReported)
	closeDesired, doneDesired := runMerger(topicDesiredState, mergedTopicDesired, topicComputedDeltaDesiredState, consumerGroupDesired, stopDesired, ctx)
	//fmt.Printf("closeReported, doneReported: %v,%v", closeDesired, doneDesired)
	// TODO consume from desired.delta and write to mqtt.messages.outgoing
	// TODO adjust code to new topology

	go func() {
		<-signals
		stopDesired <- true
		stopReported <- true
		cancel()
		closeDesired.Close()
		closeReported.Close()
	}()

	<-doneReported
	<-doneDesired
}
