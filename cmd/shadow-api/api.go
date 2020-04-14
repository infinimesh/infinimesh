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
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cskr/pubsub"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/infinimesh/infinimesh/pkg/shadow"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

var (
	broker string

	topicDesiredDelta = "shadow.desired-state.delta"

	topicReportedFull          = "shadow.reported-state.full"
	topicReportedDeltaComputed = "shadow.reported-state.delta.computed"
	topicDesiredFull           = "shadow.desired-state.full"
	topicDesiredDeltaComputed  = "shadow.desired-state.delta.computed"

	localStateMtx sync.Mutex
	localState    = make(map[string]*DeviceState)

	subMtx      sync.Mutex
	subscribers = make(map[string]map[chan *DeviceState]bool)

	dbAddr string
)

type DeviceState json.RawMessage

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "shadow.reported-state.full")
	viper.SetDefault("DB_ADDR", ":6379")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
	dbAddr = viper.GetString("DB_ADDR")
}

func subscribe(consumer sarama.Consumer, ps *pubsub.PubSub, topic, subPath string) {
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}
	fmt.Println("Consuming from " + topic)
	for _, partition := range partitions {
		go func(partition int32) {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				panic(err)
			}

			for message := range pc.Messages() {
				deltaMsg := shadow.DeviceStateMessage{}

				err := json.Unmarshal(message.Value, &deltaMsg)
				if err != nil {
					fmt.Printf("Invalid message on topic"+topic+" at offset %v, err=%v\n", message.Offset, err)
					continue
				}

				ps.Pub(&deltaMsg, string(message.Key)+subPath)

				d := DeviceState(deltaMsg.State)

				localStateMtx.Lock()
				localState[string(message.Key)] = &d
				localStateMtx.Unlock()

				// notify
				subMtx.Lock()
				if sub, ok := subscribers[string(message.Key)]; ok {
					for subscriber := range sub {
						subscriber <- &d
					}
				}
				subMtx.Unlock()
			}

		}(partition)
	}
}

func main() {
	repo, err := shadow.NewRedisRepo(dbAddr)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.Return.Successes = true

	fmt.Printf("Connect with broker %v\n", broker)
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to broker")

	ps := pubsub.New(10)

	subscribe(consumer, ps, topicReportedFull, "/reported/full")
	subscribe(consumer, ps, topicReportedDeltaComputed, "/reported/delta")
	subscribe(consumer, ps, topicDesiredFull, "/desired/full")
	subscribe(consumer, ps, topicDesiredDeltaComputed, "/desired/delta")

	go func() {
		lis, err := net.Listen("tcp", ":8096")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		srv := grpc.NewServer()

		producer, err := sarama.NewSyncProducer([]string{broker}, config)
		if err != nil {
			panic(err)
		}
		serverHandler := &shadow.Server{
			Repo:         repo,
			Producer:     producer,
			ProduceTopic: topicDesiredDelta,
			PubSub:       ps,
		}

		shadowpb.RegisterShadowsServer(srv, serverHandler)
		reflection.Register(srv)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := httprouter.New()
	r.HandlerFunc("GET", "/:id", handler)
	err = http.ListenAndServe(":8084", r)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")

	ch := make(chan *DeviceState, 10)

	// FIXME deadlock possible if kafka consume loop is blocked while writing to channel and they also hold the mutex, nobody can make progress.
	subMtx.Lock()
	if _, ok := subscribers[id]; !ok {
		subscribers[id] = make(map[chan *DeviceState]bool)
	}
	subscribers[id][ch] = true
	subMtx.Unlock()

	defer func() {
		subMtx.Lock()
		delete(subscribers[id], ch)
		subMtx.Unlock()
	}()

	notify := r.Context().Done()

outer:
	for {
		select {
		case doc := <-ch:
			str, _ := json.Marshal(json.RawMessage(*doc))
			fmt.Fprintf(w, "data: %s\n\n", str)
			flusher.Flush()
		case <-notify:
			break outer

		}

	}
}
