/*
Copyright © 2021-2022 Infinite Devices GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package plugins

import (
	"strings"

	"github.com/cskr/pubsub"
	devpb "github.com/infinimesh/proto/node/devices"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	logger  *zap.Logger
	channel *amqp.Channel
	qs      map[string]bool
)

type FetchDeviceFunc func(string) *devpb.Device

func Setup(Log *zap.Logger, conn *amqp.Connection, ps *pubsub.PubSub, fetcher FetchDeviceFunc) error {
	logger = Log.Named("PluginsExchange")

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare("plugins", "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	channel = ch
	qs = make(map[string]bool)

	go func(messages chan interface{}) {
		for msg := range messages {
			logger.Debug("Received message to Broadcast", zap.Any("msg", msg))
			shadow, ok := msg.(*pb.Shadow)
			if !ok {
				logger.Warn("Message corrupted, couldn't convert to Shadow")
			}
			dev := fetcher(shadow.GetDevice())
			if dev == nil {
				logger.Warn("Couldn't get Device from Shadow")
				continue
			}

			Publish(dev, shadow)
		}
	}(ps.Sub("mqtt.outgoing", "mqtt.incoming"))

	return nil
}

func Publish(dev *devpb.Device, state *pb.Shadow) {
	log := logger.Named(dev.GetUuid())
	log.Debug("Handling publish to Plugin Queue")
	for _, tag := range dev.Tags {
		plugin := strings.TrimPrefix(tag, "plugin:")
		log.Debug("Device data", zap.String("tag", tag), zap.String("plugin", plugin))
		if plugin != tag {
			PublishSingle(plugin, state)
		}
	}
}

func PublishSingle(plugin string, state *pb.Shadow) {
	log := logger.Named("Publisher")
	log.Debug("Publish request received", zap.String("plugin", plugin))

	_, ok := qs[plugin]
	if !ok {
		queue, err := channel.QueueDeclare(
			"plugin."+plugin, false, false, false, false, nil,
		)
		if err != nil {
			log.Warn("Error declaring queue", zap.Error(err))
			return
		}

		err = channel.QueueBind(queue.Name, plugin, "plugins", false, nil)
		if err != nil {
			log.Warn("Error binding queue", zap.Error(err))
			return
		}

		qs[plugin] = true
		log.Info("Queue declared", zap.String("name", queue.Name))
	}

	body, err := proto.Marshal(state)
	if err != nil {
		log.Warn("Error Marshaling state", zap.Error(err))
		return
	}
	err = channel.Publish("plugins", plugin, false, false, amqp.Publishing{
		ContentType: "text/plain", Body: body,
	})
	if err != nil {
		log.Warn("Error Publishing message", zap.Error(err))
		delete(qs, plugin)
		return
	}
	log.Debug("State published")
}
