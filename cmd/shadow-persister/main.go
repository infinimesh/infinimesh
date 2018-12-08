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
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var (
	addr                = "postgresql://root@localhost:26257/postgres?sslmode=disable"
	broker              string
	sourceTopicReported = "private.changelog.reported-state"
	sourceTopicDesired  = "private.changelog.desired-state"

	consumerGroup = "persister"
)

type DeviceState struct {
	ID              string
	ReportedVersion int64
	ReportedState   postgres.Jsonb
	DesiredVersion  int64
	DesiredState    postgres.Jsonb
}

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")

}

type handler struct {
	db *gorm.DB
}

func main() {
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)
	db.SingularTable(true)
	db.AutoMigrate(&DeviceState{})

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

	handler := &handler{db: db}

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

		var stateFromKafka shadow.DeviceState
		if err := json.Unmarshal(message.Value, &stateFromKafka); err != nil {
			fmt.Println("Failed to deserialize message with offset ", message.Offset)
			continue
		}

		var deviceState DeviceState

		deviceState.ID = string(message.Key)

		switch message.Topic {
		case sourceTopicReported:
			deviceState.DesiredVersion = stateFromKafka.Version
			deviceState.DesiredState = postgres.Jsonb{stateFromKafka.State} //nolint
		case sourceTopicDesired:
			deviceState.ReportedVersion = stateFromKafka.Version
			deviceState.ReportedState = postgres.Jsonb{stateFromKafka.State} //nolint
		}

		if err := h.db.Save(deviceState).Error; err != nil {
			fmt.Println("Failed to persist message with offset", message.Offset)
		}

		s.MarkMessage(message, "")
	}
	return nil
}
