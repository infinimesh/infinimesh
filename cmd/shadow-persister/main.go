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
	addr        = "postgresql://root@localhost:26257/postgres?sslmode=disable"
	broker      string
	sourceTopic = "private.changelog.reported-state"
	table       = "reported"

	consumerGroup = "persister"
)

type State struct {
	ID      string
	Version int64
	State   postgres.Jsonb
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
	db.AutoMigrate(&State{})
	db.Table(table)

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
	return nil
}

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		fmt.Println("got msg", string(message.Value))

		var state shadow.DeviceState
		if err := json.Unmarshal(message.Value, &state); err != nil {
			fmt.Println("Failed to deserialize message with offset ", message.Offset)
			continue
		}

		if err := h.db.Save(&State{
			ID:      string(message.Key),
			Version: state.Version,
			State:   postgres.Jsonb{state.State},
		}).Error; err != nil {
			fmt.Println("Failed to persist message with offset", message.Offset)
		}

		s.MarkMessage(message, "")
	}
	return nil
}
