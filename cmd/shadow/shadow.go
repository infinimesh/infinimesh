/*
Copyright © 2021-2023 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"

	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/mqtt/pubsub"
	"github.com/infinimesh/infinimesh/pkg/shadow"
	fanoutpublisher "github.com/infinimesh/infinimesh/pkg/shadow/fanout_publisher"
	"github.com/infinimesh/infinimesh/pkg/shadow/plugins"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	nodepb "github.com/infinimesh/proto/node"
	devpb "github.com/infinimesh/proto/node/devices"
	pb "github.com/infinimesh/proto/shadow"
)

var (
	log *zap.Logger

	port            string
	redisHost       string
	devicesHost     string
	RabbitMQConn    string
	SIGNING_KEY     string
	buffer_capacity int
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("REDIS_HOST", "redis:6379")
	viper.SetDefault("DEVICES_HOST", "repo:8000")
	viper.SetDefault("RABBITMQ_CONN", "amqp://infinimesh:infinimesh@rabbitmq:5672/")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")
	viper.SetDefault("BUFFER_CAPACITY", 10)

	port = viper.GetString("PORT")
	redisHost = viper.GetString("REDIS_HOST")
	devicesHost = viper.GetString("DEVICES_HOST")
	RabbitMQConn = viper.GetString("RABBITMQ_CONN")
	SIGNING_KEY = viper.GetString("SIGNING_KEY")
	buffer_capacity = viper.GetInt("BUFFER_CAPACITY")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Setting up RedisDB Connection", zap.String("host", redisHost))
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
	ps, err := pubsub.Setup(log, rbmq, "mqtt.outgoing", "mqtt.incoming", buffer_capacity)
	if err != nil {
		log.Fatal("Error setting up pubsub", zap.Error(err))
	}
	log.Info("Pub/Sub setup complete")

	SIGNING_KEY := []byte(viper.GetString("SIGNING_KEY"))
	auth := auth.NewAuthInterceptor(log, rdb, nil, SIGNING_KEY)
	token, err := auth.MakeToken(schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Error making token", zap.Error(err))
	}
	internal_ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)

	log.Info("Connecting to registry", zap.String("host", devicesHost))
	conn, err := grpc.Dial(devicesHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error dialing device registry", zap.Error(err))
	}
	client := nodepb.NewDevicesServiceClient(conn)

	fetcher_log := log.Named("DevicesFetcher")
	err = plugins.Setup(log, rbmq, ps, func(uuid string) *devpb.Device {
		fetcher_log.Debug("Attempt getting device", zap.String("uuid", uuid))
		dev, err := client.Get(internal_ctx, &devpb.Device{
			Uuid: uuid,
		})
		if err != nil {
			fetcher_log.Warn("Coudln't get device", zap.Error(err))
			return nil
		}
		return dev
	})
	if err != nil {
		log.Fatal("Error setting up plugins", zap.Error(err))
	}

	err = fanoutpublisher.Setup(log, rbmq, ps)
	if err != nil {
		log.Fatal("Error setting up fanout publisher", zap.Error(err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}

	srv := shadow.NewShadowServiceServer(log, rdb, ps)

	s := grpc.NewServer()
	pb.RegisterShadowServiceServer(s, srv)

	healthpb.RegisterHealthServer(s, health.NewServer())

	go func() {
		log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port))
		log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))
	}()

	srv.Persister()

	s.Stop()
}
