package avro

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/infinimesh/infinimesh/pkg/avro/avropb"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	"go.uber.org/zap"
)

type Consumer struct {
	Log                 *zap.Logger
	avroClient          avropb.AvroreposClient
	SourceTopicReported string
	SourceTopicDesired  string
	ConsumerGroup       string
}

func (h *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	fmt.Println("Rebalance, assigned partitions:", s.Claims())
	return nil
}

func (h *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		fmt.Println("got msg", string(message.Value))

		var stateFromKafka shadow.DeviceStateMessage
		if err := json.Unmarshal(message.Value, &stateFromKafka); err != nil {
			fmt.Println("Failed to deserialize message with offset ", message.Offset)
			continue
		}

		var dbErr error

		switch message.Topic {
		case h.SourceTopicReported:
			_, dbErr = h.avroClient.SetDeviceState(context.Background(), &avropb.SaveDeviceStateRequest{
				DeviceId:    string(message.Key),
				NamespaceId: string(message.Key),
				Version:     stateFromKafka.Version,
				Ds: &avropb.DeviceState{
					ReportedState: stateFromKafka.State,
					//DesiredState:,
				}})
		case h.SourceTopicDesired:
			_, dbErr = h.avroClient.SetDeviceState(context.Background(), &avropb.SaveDeviceStateRequest{
				DeviceId:    string(message.Key),
				NamespaceId: string(message.Key),
				Version:     stateFromKafka.Version,
				Ds: &avropb.DeviceState{
					//ReportedState:,
					DesiredState: stateFromKafka.State,
				}})
		}
		// FIXME ignore errors for now
		if dbErr != nil {
			fmt.Println("Failed to persist message with offset", message.Offset, dbErr)
		}

		s.MarkMessage(message, "")
	}
	return nil
}
