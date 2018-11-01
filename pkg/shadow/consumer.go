package shadow

import (
	"fmt"

	"time"

	sarama "github.com/Shopify/sarama"
)

type StateMerger struct {
	SourceTopic    string
	ChangelogTopic string

	localState           map[string]string // device id to state string
	localStateMaxOffsets map[int32]int64

	ChangelogConsumerClient sarama.Client
	ChangelogProducerClient sarama.Client

	changelogProducer sarama.SyncProducer
}

func (c *StateMerger) fetchLocalState(client sarama.Client, partitions []int32) (localState map[string]string, offsets map[int32]int64, err error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		go consumer.Close()
	}()

	localState = make(map[string]string)

	offsets = make(map[int32]int64)
	for _, partition := range partitions {
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
			localState[string(item.Key)] = string(item.Value)

			if item.Offset == pc.HighWaterMarkOffset()-1 {
				break
			}

		}
	}
	return
}

func (c *StateMerger) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	c.localState = make(map[string]string)

	producer, err := sarama.NewSyncProducerFromClient(c.ChangelogProducerClient)
	if err != nil {
		return err
	}

	c.changelogProducer = producer

	partitionsToFetch, ok := s.Claims()[c.SourceTopic]
	if !ok {
		fmt.Println("No partitions assigned. sleeping.")
	}

	start := time.Now()
	localState, offsets, err := c.fetchLocalState(c.ChangelogConsumerClient, partitionsToFetch)
	if err != nil {
		return err
	}

	c.localState = localState
	c.localStateMaxOffsets = offsets

	fmt.Printf("Restored local state: %v shadows in %v seconds.\n", len(localState), time.Now().Sub(start).Seconds())
	fmt.Println(localState)
	return nil
}
func (h *StateMerger) Cleanup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Cleaning consumer group session")

	h.changelogProducer.Close()
	h.changelogProducer = nil
	h.localStateMaxOffsets = nil
	h.localState = nil

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

func (h *StateMerger) processMessage(msg *sarama.ConsumerMessage) {
	localState, ok := h.localState[string(msg.Key)]
	if !ok {
		fmt.Println("Didn't find local state. assuming {}")
		localState = "{}"
	}

	newState, err := applyDelta(localState, string(msg.Value))
	if err != nil {
		fmt.Println("Failed to apply delta, ignoring msg")
		return
	}

	if newState == localState {
		return
	}

	h.localState[string(msg.Key)] = newState

	// TODO Model checks would happen here

	_, _, err = h.changelogProducer.SendMessage(&sarama.ProducerMessage{
		Topic:     h.ChangelogTopic,
		Key:       sarama.StringEncoder(msg.Key),
		Value:     sarama.StringEncoder(newState),
		Partition: msg.Partition, // Always use the same partition number in the target topic as in the source topic
	})
}

func (h *StateMerger) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		batch := make([]*sarama.ConsumerMessage, 100)
		n, ok := consumeBatch(claim.Messages(), batch)
		if !ok {
			break
		}
		for _, x := range batch[:n] {
			h.processMessage(x)

		}
		if n > 0 {
			// We mark the message after we wrote them to the local-state log. If the process crashes before the next line,
			// but after processMessage (or before the marked message offset is communicated to the broker as async). \
			// After restart the message will be processed again (at least once semantics).
			sess.MarkMessage(batch[n-1], "")
		}
	}

	// for msg := range claim.Messages() {
	// 	switch msg.Topic {
	// 	case h.SourceTopic:
	// 	case h.TopicDesired:
	// 	}
	// 	fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
	// 	sess.MarkMessage(msg, "")
	// }
	return nil
}
