/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
package pubsub

import (
	"time"

	"github.com/cskr/pubsub"
	pb "github.com/infinimesh/proto/shadow"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	ps     *pubsub.PubSub
	logger *zap.Logger
)

func Setup(Log *zap.Logger, conn *amqp.Connection, pub, sub string) (*pubsub.PubSub, error) {
	logger = Log
	ps = pubsub.New(10)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	go HandlePublish(ch, pub)
	go HandleSubscribe(ch, sub)

	return ps, nil
}

// Reading messages from PubSub and publishing them to RabbitMQ Queue
func HandlePublish(ch *amqp.Channel, topic string) {
	log := logger.Named("publish")
init:
	q, err := ch.QueueDeclare(
		topic,
		true, false, false, true, nil,
	)
	if err != nil {
		log.Warn("Error declaring queue", zap.Error(err))
		time.Sleep(time.Second)
		goto init
	}
	log.Info("Queue declared", zap.String("name", q.Name))

	incoming := make(chan interface{}, 10)
	ps.AddSub(incoming, topic)

	for msg := range incoming {
		shadow := msg.(*pb.Shadow)
		log.Debug("Received message from PubSub", zap.Any("shadow", shadow))
		payload, err := proto.Marshal(shadow)
		if err != nil {
			log.Warn("Error while publishing message:", zap.Error(err))
			continue
		}
		ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain", Body: payload,
		})
	}
}

// Reading messages from RabbitMQ Queue and publishing them to PubSub
func HandleSubscribe(ch *amqp.Channel, topic string) {
	log := logger.Named("subscribe")
init:
	q, err := ch.QueueDeclare(
		topic,
		true, false, false, true, nil,
	)
	if err != nil {
		log.Warn("Error declaring queue", zap.Error(err))
		time.Sleep(time.Second)
		goto init
	}
	log.Info("Queue declared", zap.String("name", q.Name))

consume:
	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Warn("Error setting up consumer", zap.Error(err))
		time.Sleep(time.Second)
		goto consume
	}

	for msg := range messages {
		shadow := &pb.Shadow{}
		err = proto.Unmarshal(msg.Body, shadow)
		if err != nil {
			log.Warn("Error while consuming message:", zap.Error(err))
			continue
		}
		log.Debug("Received message from RabbitMQ", zap.Any("shadow", &shadow))
		ps.Pub(shadow, topic, topic+"/"+shadow.Device)
	}
}
