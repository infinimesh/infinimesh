package fanoutpublisher

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/pubsub"
	"github.com/infinimesh/infinimesh/pkg/shadow/fanout"
	pb "github.com/infinimesh/proto/shadow"
	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

func Setup(logger *zap.Logger, conn *amqp.Connection, ps pubsub.PubSub) error {
	log := logger.Named("FanoutPublisher")
	f, err := fanout.Setup(log, conn, "shadow.fanout", 1000)
	if err != nil {
		return err
	}

	go func(messages chan interface{}) {
		for msg := range messages {
			log.Debug("Received message to Broadcast", zap.Any("msg", msg))
			shadow, ok := msg.(*pb.Shadow)
			if !ok {
				log.Warn("Message corrupted, couldn't convert to Shadow")
				continue
			}

			err := f.Publish(context.Background(), shadow)
			if err != nil {
				log.Warn("Error publishing to fanout", zap.Error(err))
			}
		}
	}(ps.Sub("mqtt.outgoing", "mqtt.incoming"))

	return nil
}
