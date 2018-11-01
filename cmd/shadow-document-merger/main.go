package main

import (
	"context"
	"fmt"
	"sync"

	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/infinimesh/infinimesh/pkg/shadow"
)

var (
	brokers         = []string{"localhost:9092"}
	groupIDreported = "mixer-reported"
	groupIDdesired  = "mixer-desired"
	topicReported   = "private.changelog.reported-state"
	topicDesired    = "private.changelog.desired-state"
	changelogTopic  = "public.shadow.states"
)

func main() {
	// TODO create topics / at least check if they are correct, e.g. co-partitioned and compacted
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Strategy = &shadow.BalanceStrategyCoPartitioned{}
	config.Producer.Retry.Max = 10

	// Ensure co-partitioning

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		panic(err)
	}

	group, err := sarama.NewConsumerGroupFromClient(groupIDreported, client)
	if err != nil {
		panic(err)
	}

	h := &handler{
		config: config,
		client: client,
	}

	for {
		_ = group.Consume(context.Background(), []string{topicReported, topicDesired}, h)
	}
}

type handler struct {
	config *sarama.Config
	client sarama.Client

	producer sarama.AsyncProducer

	// localState map[string]*DeviceState

	// Partitions -> DeviceID -> State this allows us to fetch the
	// correct maps once after rebalance, instead of locking it every time
	// we access it
	m           sync.Mutex
	localStates map[int32]map[string]*DeviceState
}

type DeviceState struct {
	Version  int64
	Reported json.RawMessage
	Desired  json.RawMessage
}

func (h *handler) fetchLocalState(partitions []int32) (localStates map[int32]map[string]*DeviceState, err error) {
	consumer, err := sarama.NewConsumerFromClient(h.client)
	if err != nil {
		return nil, err
	}
	defer func() {
		go consumer.Close()
	}()

	localStates = make(map[int32]map[string]*DeviceState)

	// TODO do this in parallel. beware, don't access localStates concurrently without mutex.
	for _, partition := range partitions {
		localStates[partition] = make(map[string]*DeviceState)
		pc, err := consumer.ConsumePartition(changelogTopic, partition, int64(0))
		if err != nil {
			return nil, err
		}
		defer pc.Close()

		newestOffset, err := h.client.GetOffset(changelogTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return nil, err
		}

		if newestOffset == 0 {
			continue
		}

		for message := range pc.Messages() {
			key := string(message.Key)

			var d *DeviceState
			if localState, ok := localStates[partition][key]; ok {
				d = localState
			} else {
				d = &DeviceState{}
				localStates[partition][key] = d
			}

			err := json.Unmarshal(message.Value, &d)
			if err != nil {
				panic(err)
			}

			if message.Offset == pc.HighWaterMarkOffset()-1 {
				break
			}

		}
	}
	return
}

func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Rebalance, assigned partitions:", s.Claims())

	// Fetch state from changelog
	partitionsToFetch, ok := s.Claims()[topicReported]
	if !ok {
		fmt.Println("No partitions assigned. sleeping.")
	}

	localStates, err := h.fetchLocalState(partitionsToFetch)
	if err != nil {
		return err
	}

	h.localStates = localStates
	c, _ := sarama.NewClient(brokers, h.config)
	producer, err := sarama.NewAsyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	h.producer = producer

	return nil
}

func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error {
	h.producer.Close()
	h.localStates = nil
	fmt.Println("Cleaning consumer group session")
	return nil
}

func (h *handler) processMessage(message *sarama.ConsumerMessage) {
}

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.m.Lock()
	localState := h.localStates[claim.Partition()] // local state for exactly this partition
	h.m.Unlock()
	for message := range claim.Messages() {
		key := string(message.Key)

		var deviceState *DeviceState
		if ds, ok := localState[key]; !ok {
			deviceState = &DeviceState{}
			localState[key] = deviceState
		} else {
			deviceState = ds
		}

		if message.Topic == topicDesired {
			deviceState.Desired = json.RawMessage(string(message.Value))
		} else if message.Topic == topicReported {
			deviceState.Reported = json.RawMessage(string(message.Value))
		}
		deviceState.Version++

		stateDocument, err := json.Marshal(deviceState)
		if err != nil {
			panic(err)
		}

		h.producer.Input() <- &sarama.ProducerMessage{
			Topic: changelogTopic,
			Key:   sarama.StringEncoder(message.Key),
			Value: sarama.StringEncoder(stateDocument),
		}

		s.MarkMessage(message, "")
	}
	return nil
}
