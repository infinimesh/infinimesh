package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/infinimesh/infinimesh/pkg/proto/api"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	fmt.Printf("Connect with broker %v\n", broker)
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to broker")

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}
	for _, partition := range partitions {
		go func(partition int32) {
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				panic(err)
			}

			for message := range pc.Messages() {
				fmt.Println("Got Msg", message.Value)
				rawMessage := json.RawMessage{}
				err := json.Unmarshal(message.Value, &rawMessage)
				if err != nil {
					fmt.Printf("Invalid message at offset %v, err=%v\n", message.Offset, err)
					continue
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

	go func() {
		lis, err := net.Listen("tcp", ":8096")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		srv := grpc.NewServer()
		serverHandler := &Server{}
		api.RegisterShadowServer(srv, serverHandler)
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

type Server struct{}

func (s *Server) GetReported(context.Context, *api.GetReportedRequest) (*api.GetReportedResponse, error) {
	return &api.GetReportedResponse{}, nil
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
