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
package state

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	redis "github.com/go-redis/redis/v8"
	pb "github.com/infinimesh/infinimesh/pkg/state/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StateServer struct {
	pb.UnimplementedStateServiceServer
	
	log *zap.Logger
	rdb *redis.Client

	NewRDB func() *redis.Client
}

var STATE_KEYS_PREFIX = "_st"
var rdbCtx context.Context = context.Background()

func NewStateServer(log *zap.Logger, redis_host string) *StateServer {
	rdb := redis.NewClient(&redis.Options{
		Addr: redis_host,
	})
	return &StateServer{
		log: log,
		rdb: rdb,
		NewRDB: func() *redis.Client {
			return redis.NewClient(&redis.Options{
				Addr: redis_host,
			})
		},
	}
}

func (s *StateServer) Post(ctx context.Context, req *pb.PostRequest) (*pb.EmptyMessage, error) {
	log := s.log.Named("Post")
	log.Debug("Request", zap.Any("req", req))

	key := fmt.Sprintf("%s:%s", STATE_KEYS_PREFIX, req.GetId())
	json, err := json.Marshal(req.GetData())
	if err != nil {
		s.log.Error("Error Marshal JSON",
			zap.String("key", key), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error  Marshal JSON")
	}

	r := s.rdb.Set(rdbCtx, key, json, 0)
	_, err = r.Result()
	if err != nil {
		s.log.Error("Error putting status to Redis",
			zap.String("key",key), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error putting status to Redis")
	}
	log.Debug("State stored in Redis")

	go func() {
		log.Debug("Storing State in Redis Channel", zap.String("key", key))
		err = s.rdb.Publish(rdbCtx, key, json).Err()
		if err != nil {
			s.log.Error("Error putting status to Redis channel", zap.String("key", key), zap.Error(err))
			return
		}
		log.Debug("State stored in Redis Channel")
	}()

	return &pb.EmptyMessage{}, nil
}

func (s *StateServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log := s.log.Named("Get")
	log.Debug("Request", zap.Any("req", req))

	states := make(map[string]*pb.State)

	keys := req.GetIds()
	for i, uuid := range keys {
		keys[i] = fmt.Sprintf("%s:%s", STATE_KEYS_PREFIX, uuid)
	}

	r := s.rdb.MGet(ctx, keys...)
	response, err := r.Result()
	if err != nil {
		log.Error("Error getting states from Redis", zap.Strings("keys", keys), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting states from Redis")
	}

	for i, state := range response {
		var dev_state pb.State
		key := strings.Replace(keys[i], STATE_KEYS_PREFIX + ":", "", 1)

		switch state := state.(type) {
		case string:
			err = json.Unmarshal([]byte(state), &dev_state)
		case nil:
			continue
		}
		if err != nil {
			log.Error("Error Unmarshal JSON",
				zap.String("key", keys[i]), zap.Error(err))
			return nil, status.Error(codes.Internal, "Error Unmarshal JSON")
		}

		states[key] = &dev_state
	}

	return &pb.GetResponse{States: states}, nil
}
