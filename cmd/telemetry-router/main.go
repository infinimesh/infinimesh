package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/infinimesh/infinimesh/pkg/mqtt"
	"github.com/infinimesh/infinimesh/pkg/router"
	"github.com/spf13/viper"
)

type handler struct {
	producer sarama.AsyncProducer
	router   *router.Router
}

var (
	consumerGroup = "dispatcher"
	broker        string

	sourceTopic  = "mqtt.messages.incoming"
	defaultRoute = "public.bridge.dlq"
	routes       = map[string]string{
		"shadows/": "public.delta.reported-state",
	}
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
}

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Errors = false
	config.Producer.Return.Successes = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = false
	config.Version = sarama.V2_0_0_0

	client, err := sarama.NewClient([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	producer, err := sarama.NewAsyncProducerFromClient(client)
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

	handler := &handler{
		producer: producer,
		router:   router.New(defaultRoute, routes),
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {
	outer:
		for {

			err = group.Consume(context.Background(), []string{sourceTopic}, handler)
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
	h.producer.Close()
	return nil
}

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg mqtt.IncomingMessage
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to deserialize msg with offset %v", message.Offset)
		}

		target := h.router.Route(msg.SourceTopic)

		h.producer.Input() <- &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(msg.SourceDevice),
			Topic: target,
			Value: sarama.ByteEncoder(msg.Data),
		}

		s.MarkMessage(message, "")
	}
	return nil
}
