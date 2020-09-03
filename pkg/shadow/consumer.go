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
	SourceTopic    string // Incoming ticks
	MergedTopic    string // Full state with version
	RealDeltaTopic string // Deltas for each full state transition, with version

	m           sync.Mutex
	localStates map[int32]map[string]*DeviceStateMessage // device id to state string

	ChangelogConsumerClient sarama.Client
	ChangelogProducerClient sarama.Client

	changelogProducer sarama.AsyncProducer
}

func (c *StateMerger) fetchLocalState(client sarama.Client, partitions []int32) (localStates map[int32]map[string]*DeviceStateMessage, offsets map[int32]int64, err error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		go consumer.Close()
	}()

	localStates = make(map[int32]map[string]*DeviceStateMessage)

	offsets = make(map[int32]int64)
	for _, partition := range partitions {
		fmt.Printf("Consumer partition reading : %v\n", partition)
		localStates[partition] = make(map[string]*DeviceStateMessage)
		offsets[partition] = 0
		pc, err := consumer.ConsumePartition(c.MergedTopic, partition, sarama.OffsetOldest)
		fmt.Printf("Partition Consumer :%v", pc)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("PC Messages :%v", pc.Messages())
		defer pc.Close()

		newestOffset, err := client.GetOffset(c.MergedTopic, partition, sarama.OffsetNewest)
		fmt.Printf("Consumer partition newestOffset : %v\n", newestOffset)
		if err != nil {
			return nil, nil, err
		}

		if newestOffset == 0 {
			continue
		}

		for item := range pc.Messages() {
			fmt.Printf("Consumer partition item : %v\n", item)
			var st DeviceStateMessage

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
	c.localStates = make(map[int32]map[string]*DeviceStateMessage)
	fmt.Printf("localstates initialized %v", c.localStates)
	//TODO enforce co-partitioning

	producer, err := sarama.NewAsyncProducerFromClient(c.ChangelogProducerClient)
	fmt.Printf("producer client created in shadow consumer %v", producer)
	if err != nil {
		return err
	}

	c.changelogProducer = producer

	partitionsToFetch, ok := s.Claims()[c.SourceTopic]
	fmt.Printf("Partitions to fetch %v", partitionsToFetch)
	if !ok {
		fmt.Println("No partitions assigned. sleeping.")
	}

	start := time.Now()
	localStates, offsets, err := c.fetchLocalState(c.ChangelogConsumerClient, partitionsToFetch)
	fmt.Printf("Local states fetched %v and offsets %v", localStates, offsets)
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
		sess.MarkMessage(message, "")
		fmt.Println("Recv message", string(message.Value))
		key := string(message.Key)

		var deviceState *DeviceStateMessage
		if ds, ok := localState[key]; !ok {
			deviceState = &DeviceStateMessage{}
			localState[key] = deviceState
		} else {
			deviceState = ds
		}

		// Input: any JSON
		delta := string(message.Value)
		old := string(deviceState.State)

		newState, err := applyDelta(old, delta)
		fmt.Printf("Consumer newState : %v", newState)
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
				SubPath:  "state/desired/delta",
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
		deviceState.Timestamp = time.Now() // TODO support user-provided timestamps

		stateDocument, err := json.Marshal(deviceState)
		if err != nil {
			fmt.Println("Failed to marshal JSON", err)
			continue
		}

		h.changelogProducer.Input() <- &sarama.ProducerMessage{
			Topic:     h.MergedTopic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.StringEncoder(stateDocument),
			Partition: message.Partition,
		}

		fmt.Println("Send msg to ", h.MergedTopic, " ", string(stateDocument))

		// Determine actual delta
		deltaMsg := &DeviceStateMessage{
			Version:   deviceState.Version,
			State:     []byte(mergePatch),
			Timestamp: time.Now(),
		}

		deltaMsgBytes, err := json.Marshal(&deltaMsg)
		if err != nil {
			fmt.Printf("Failed to marshal delta msg: %v\n", err)
		}

		h.changelogProducer.Input() <- &sarama.ProducerMessage{
			Topic:     h.RealDeltaTopic,
			Key:       sarama.StringEncoder(message.Key),
			Value:     sarama.ByteEncoder(deltaMsgBytes),
			Partition: message.Partition,
		}

		fmt.Println("Send msg to ", h.RealDeltaTopic, " ", string(deltaMsgBytes))
	}
	return nil
}
