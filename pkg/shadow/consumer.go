package shadow

import (
	"fmt"
	"sync"

	"time"

	"encoding/json"

	sarama "github.com/Shopify/sarama"
)

type DeviceState struct {
	Version int64
	State   json.RawMessage
}

type StateMerger struct {
	SourceTopic    string
	ChangelogTopic string

	m           sync.Mutex
	localStates map[int32]map[string]*DeviceState // device id to state string

	ChangelogConsumerClient sarama.Client
	ChangelogProducerClient sarama.Client

	changelogProducer sarama.AsyncProducer
}

func (c *StateMerger) fetchLocalState(client sarama.Client, partitions []int32) (localStates map[int32]map[string]*DeviceState, offsets map[int32]int64, err error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		go consumer.Close()
	}()

	localStates = make(map[int32]map[string]*DeviceState)

	offsets = make(map[int32]int64)
	for _, partition := range partitions {
		localStates[partition] = make(map[string]*DeviceState)
		offsets[partition] = 0
		pc, err := consumer.ConsumePartition(c.ChangelogTopic, partition, int64(0))
		if err != nil {
			return nil, nil, err
		}
		defer pc.Close()

		newestOffset, err := client.GetOffset(c.ChangelogTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return nil, nil, err
		}

		if newestOffset == 0 {
			continue
		}

		for item := range pc.Messages() {
			var st DeviceState

			err := json.Unmarshal(item.Value, &st)
			if err != nil {
				return nil, nil, err
			}

			localStates[partition][string(item.Key)] = &st

			if item.Offset == pc.HighWaterMarkOffset()-1 {
				break
			}

		}
	}
	return
}

func (c *StateMerger) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Work with topic", c.SourceTopic)
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	c.localStates = make(map[int32]map[string]*DeviceState)

	//TODO enforce co-partitioning

	producer, err := sarama.NewAsyncProducerFromClient(c.ChangelogProducerClient)
	if err != nil {
		return err
	}

	c.changelogProducer = producer

	partitionsToFetch, ok := s.Claims()[c.SourceTopic]
	if !ok {
		fmt.Println("No partitions assigned. sleeping.")
	}

	start := time.Now()
	localStates, _, err := c.fetchLocalState(c.ChangelogConsumerClient, partitionsToFetch)
	if err != nil {
		return err
	}

	c.localStates = localStates
	// c.localStateMaxOffsets = offsets

	fmt.Printf("Restored local state: %v shadows in %v seconds.\n", len(localStates), time.Since(start).Seconds())
	fmt.Println(localStates)
	return nil
}
func (h *StateMerger) Cleanup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Cleaning consumer group session")

	h.changelogProducer.Close()
	h.changelogProducer = nil
	// h.localStateMaxOffsets = nil
	h.localStates = nil

	fmt.Println("return")
	return nil
}

func consumeBatch(c <-chan *sarama.ConsumerMessage, buf []*sarama.ConsumerMessage) (n int, ok bool) {
	item, ok := <-c
	if !ok {
		return 0, false
	}
	buf[0] = item

	i := 0

inner:
	for i = 1; i < len(buf); i++ {
		select {

		case item, ok := <-c:
			if !ok {
				return i, false
			}
			buf[i] = item
		default:
			break inner
		}
	}

	return i, true

}

// Topic infinimesh.bridge.incoming.raw
type MQTTBridgeData struct {
	SourceTopic  string
	SourceDevice string
	Data         []byte
}

func (h *StateMerger) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

		delta := string(message.Value)
		old := string(deviceState.State)

		newState, err := applyDelta(old, delta)

		if newState == old {
			fmt.Println("No change, skip")
			continue
		}

		deviceState.State = json.RawMessage(newState)
		deviceState.Version++

		stateDocument, err := json.Marshal(deviceState)
		if err != nil {
			panic(err)
		}

		h.changelogProducer.Input() <- &sarama.ProducerMessage{
			Topic:     h.ChangelogTopic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.StringEncoder(stateDocument),
			Partition: message.Partition,
		}
		sess.MarkMessage(message, "")
	}
	return nil
}
