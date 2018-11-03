package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

var (
	broker string
	topic  string

	localStateMtx sync.Mutex
	localState    = make(map[string]*DeviceState)

	subMtx      sync.Mutex
	subscribers = make(map[string]map[chan *DeviceState]bool)
)

type DeviceState json.RawMessage

func init() {
	viper.SetDefault("KAFKA_HOST", "localhost:9092")
	viper.SetDefault("KAFKA_TOPIC", "private.changelog.reported-state")
	viper.AutomaticEnv()

	broker = viper.GetString("KAFKA_HOST")
	topic = viper.GetString("KAFKA_TOPIC")
}

func main() {

	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		panic(err)
	}

	partitions, err := consumer.Partitions(topic)
	for _, partition := range partitions {
		go func(partition int32) {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				panic(err)
			}

			for message := range pc.Messages() {
				rawMessage := json.RawMessage{}
				err := json.Unmarshal(message.Value, &rawMessage)
				if err != nil {
					fmt.Printf("Invalid message at offset %v, err=%v\n", message.Offset, err)
				}

				d := DeviceState(rawMessage)

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

	r := httprouter.New()
	r.HandlerFunc("GET", "/:id", handler)
	err = http.ListenAndServe(":8085", r)
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

	ch := make(chan *DeviceState)
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

	notify := w.(http.CloseNotifier).CloseNotify()

outer:
	for {
		select {
		case doc := <-ch:
			str, _ := json.Marshal(json.RawMessage(*doc))
			fmt.Fprintf(w, "data: %s\n\n", str)
			flusher.Flush()
		case _ = <-notify:
			break outer

		}

	}

	return
}
