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
package main

import (
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/mqtt/pubsub"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	pb "github.com/infinimesh/proto/shadow"
)

var (
	port         string
	redisHost    string
	RabbitMQConn string
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("REDIS_HOST", "redis:6379")
	viper.SetDefault("RABBITMQ_CONN", "amqp://infinimesh:infinimesh@rabbitmq:5672/")

	port = viper.GetString("PORT")
	redisHost = viper.GetString("REDIS_HOST")
	RabbitMQConn = viper.GetString("RABBITMQ_CONN")
}

func main() {
	log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Setting up RedisDB Connection")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0, // use default DB
	})
	log.Info("RedisDB connection established")

	log.Info("Connecting to RabbitMQ", zap.String("url", RabbitMQConn))
	rbmq, err := amqp.Dial(RabbitMQConn)
	if err != nil {
		log.Fatal("Error dialing RabbitMQ", zap.Error(err))
	}
	defer rbmq.Close()
	log.Info("Connected to RabbitMQ")

	log.Info("Setting up Pub/Sub")
	ps, err := pubsub.Setup(log, rbmq, "mqtt.outgoing", "mqtt.incoming")
	if err != nil {
		log.Fatal("Error setting up pubsub", zap.Error(err))
	}
	log.Info("Pub/Sub setup complete")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}

	srv := shadow.NewShadowServiceServer(log, rdb, ps)
	go srv.Persister()

	s := grpc.NewServer()
	pb.RegisterShadowServiceServer(s, srv)

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))
}
