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
	consumerGroupReported = "shadow-reported"
	consumerGroupDesired  = "shadow-desired"
	broker                string
	topicReportedState    = "shadow.reported-state.delta"
	topicDesiredState     = "shadow.desired-state.delta"
	mergedTopicReported   = "shadow.reported-state.full"
	mergedTopicDesired    = "shadow.desired-state.full"
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
}

func runMerger(inputTopic, outputTopic, consumerGroup string, stop chan bool, ctx context.Context) (close io.Closer, done chan bool) {
	done = make(chan bool)
	consumerGroupClient := sarama.NewConfig()
	consumerGroupClient.Version = sarama.V1_0_0_0
	consumerGroupClient.Consumer.Return.Errors = true
	consumerGroupClient.Consumer.Offsets.Initial = sarama.OffsetOldest

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
		SourceTopic:             inputTopic,
		MergedTopic:             outputTopic,
		ChangelogProducerClient: producerClient,
		ChangelogConsumerClient: localStateConsumerClient,
	}

	go func() {
	outer:
		for {

			err = group.Consume(ctx, []string{inputTopic}, handler)
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

	closeReported, doneReported := runMerger(topicReportedState, mergedTopicReported, consumerGroupReported, stopReported, ctx)
	closeDesired, doneDesired := runMerger(topicDesiredState, mergedTopicDesired, consumerGroupDesired, stopDesired, ctx)

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
