package shadow

import (
	"fmt"
	"sync"

	"time"

	"encoding/json"

	sarama "github.com/Shopify/sarama"

	"github.com/infinimesh/infinimesh/pkg/mqtt"
)

type StateMerger struct {
	SourceTopic string
	MergedTopic string

	m           sync.Mutex
	localStates map[int32]map[string]*FullDeviceStateMessage // device id to state string

	ChangelogConsumerClient sarama.Client
	ChangelogProducerClient sarama.Client

	changelogProducer sarama.AsyncProducer
}

func (c *StateMerger) fetchLocalState(client sarama.Client, partitions []int32) (localStates map[int32]map[string]*FullDeviceStateMessage, offsets map[int32]int64, err error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		go consumer.Close()
	}()

	localStates = make(map[int32]map[string]*FullDeviceStateMessage)

	offsets = make(map[int32]int64)
	for _, partition := range partitions {
		localStates[partition] = make(map[string]*FullDeviceStateMessage)
		offsets[partition] = 0
		pc, err := consumer.ConsumePartition(c.MergedTopic, partition, sarama.OffsetOldest)
		if err != nil {
			return nil, nil, err
		}
		defer pc.Close()

		newestOffset, err := client.GetOffset(c.MergedTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return nil, nil, err
		}

		if newestOffset == 0 {
			continue
		}

		for item := range pc.Messages() {
			var st FullDeviceStateMessage

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
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	c.localStates = make(map[int32]map[string]*FullDeviceStateMessage)

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
	h.localStates = nil

	return nil
}

func (h *StateMerger) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.m.Lock()
	localState := h.localStates[claim.Partition()] // local state for exactly this partition
	h.m.Unlock()
	for message := range claim.Messages() {
		defer sess.MarkMessage(message, "")
		fmt.Println("Recv message", string(message.Value))
		key := string(message.Key)

		var deviceState *FullDeviceStateMessage
		if ds, ok := localState[key]; !ok {
			deviceState = &FullDeviceStateMessage{}
			localState[key] = deviceState
		} else {
			deviceState = ds
		}

		// Input: any JSON
		delta := string(message.Value)
		old := string(deviceState.State)

		newState, err := applyDelta(old, delta)
		if err != nil {
			fmt.Println("Failed to apply new delta. Ignoring message", err)
			continue
		}

		if newState == old {
			fmt.Println("No change, skip")
			continue
		}

		// We got a change, so also publish a message to the broker.
		// TODO better split this into multiple topics; ticks (deltas, that may or may not result in a change) -> full state + deltas

		mergePatch := calculateDelta(old, newState)

		// TODO workaround so we don't write needless messages in case of reported states
		if h.MergedTopic == "shadow.desired-state.full" {
			outgoing := mqtt.OutgoingMessage{
				DeviceID: string(message.Key),
				SubPath:  "shadow/updates",
				Data:     []byte(mergePatch),
			}

			outBytes, err := json.Marshal(&outgoing)
			if err != nil {
				fmt.Printf("Failed to marshal outgoing msg: %v\n", err)
			}

			h.changelogProducer.Input() <- &sarama.ProducerMessage{
				Topic: "mqtt.messages.outgoing",
				Key:   sarama.StringEncoder(message.Key),
				Value: sarama.ByteEncoder(outBytes),
			}

			fmt.Println("Send msg to ", "mqtt.messages.outgoing", " ", string(outBytes))
		}

		deviceState.State = json.RawMessage(newState)
		deviceState.Version++

		stateDocument, err := json.Marshal(deviceState)
		if err != nil {
			panic(err)
		}

		h.changelogProducer.Input() <- &sarama.ProducerMessage{
			Topic:     h.MergedTopic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.StringEncoder(stateDocument),
			Partition: message.Partition,
		}

		fmt.Println("Send msg to ", h.MergedTopic, " ", string(stateDocument))
	}
	return nil
}
