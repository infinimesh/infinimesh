package main

import (
	"context"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	sarama "github.com/Shopify/sarama"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	"github.com/spf13/viper"
)

var (
	consumerGroup = "shadow"
	broker        string
	topics        = []string{topic}
	topic         = "public.delta.reported-state"
	mergedTopic   = "private.changelog.reported-state"
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
}

func main() {
	consumerGroupClient := sarama.NewConfig()
	consumerGroupClient.Version = sarama.V1_0_0_0
	consumerGroupClient.Consumer.Return.Errors = true
	consumerGroupClient.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroupClient.Consumer.Group.Member.UserData = []byte("test:8080")

	client, err := sarama.NewClient([]string{broker}, consumerGroupClient)
	if err != nil {
		panic(err)
	}

	group, err := sarama.NewConsumerGroupFromClient(consumerGroup, client)
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
		SourceTopic:             topic,
		ChangelogTopic:          mergedTopic,
		ChangelogProducerClient: producerClient,
		ChangelogConsumerClient: localStateConsumerClient,
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {
	outer:
		for {

			err = group.Consume(context.Background(), topics, handler)
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
