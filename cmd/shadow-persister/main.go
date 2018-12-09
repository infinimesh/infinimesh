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
	"github.com/infinimesh/infinimesh/pkg/shadow"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

const (
	sourceTopicReported = "private.changelog.reported-state"
	sourceTopicDesired  = "private.changelog.desired-state"
)

var (
	addr   = "postgresql://root@localhost:26257/postgres?sslmode=disable"
	broker string

	consumerGroup = "persister"
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("DB_ADDR", "postgresql://root@localhost:26257/postgres?sslmode=disable")
	viper.AutomaticEnv()
	broker = viper.GetString("KAFKA_HOST")
	addr = viper.GetString("DB_ADDR")
}

type handler struct {
	repo shadow.Repo
}

func main() {
	repo, err := shadow.NewPostgresRepo(addr)
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

		var stateFromKafka shadow.FullDeviceStateMessage
		if err := json.Unmarshal(message.Value, &stateFromKafka); err != nil {
			fmt.Println("Failed to deserialize message with offset ", message.Offset)
			continue
		}

		var dbErr error

		switch message.Topic {
		case sourceTopicReported:
			dbErr = h.repo.SetReported(shadow.DeviceState{
				ID:      string(message.Key),
				Version: stateFromKafka.Version,
				State:   string(stateFromKafka.State),
			})
		case sourceTopicDesired:
		}

		if dbErr != nil {
			fmt.Println("Failed to persist message with offset", message.Offset, dbErr)
		}

		s.MarkMessage(message, "")
	}
	return nil
}
