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
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	inflog "github.com/slntopp/infinimesh/pkg/log"
	"github.com/slntopp/infinimesh/pkg/timeseries"
)

const (
	sourceTopic = "shadow.reported-state.delta.computed"
)

var (
	addr   string
	broker string

	log *zap.Logger

	consumerGroup = "timescale-connector-import"
)

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("DB_ADDR", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	viper.AutomaticEnv()
	broker = viper.GetString("KAFKA_HOST")
	addr = viper.GetString("DB_ADDR")

	log, _ = inflog.NewProdOrDev()
}

func main() {
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

	repo, err := timeseries.NewTimescaleRepo(
		log.Named("TimescaleRepo"),
		addr,
	)
	fmt.Printf("timescale called : %v", repo)
	if err != nil {
		panic(err)
	}

	handler := &timeseries.Consumer{
		Log:  log.Named("Consumer"),
		Repo: repo,
	}
	fmt.Printf("timescale Consumer : %v", handler)

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
