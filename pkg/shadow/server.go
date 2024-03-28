/*
Copyright © 2018-2024 Infinite Devices GmbH Nikita Ivanovski info@slnt-opp.xyz

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
package shadow

import (
	"context"
	"errors"

	"encoding/json"

	"connectrpc.com/connect"
	redis "github.com/go-redis/redis/v8"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/pubsub"
	devpb "github.com/infinimesh/proto/node/devices"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/infinimesh/proto/shadow"
)

type ShadowServiceServer struct {
	pb.UnimplementedShadowServiceServer

	log  *zap.Logger
	rdb  redis.Cmdable
	ps   pubsub.PubSub
	repo graph.InfinimeshGenericActionsRepo[*devpb.Device]
}

func NewShadowServiceServer(log *zap.Logger, rdb redis.Cmdable, ps pubsub.PubSub, repo graph.InfinimeshGenericActionsRepo[*devpb.Device]) *ShadowServiceServer {
	return &ShadowServiceServer{
		log:  log.Named("shadow"),
		rdb:  rdb,
		ps:   ps,
		repo: repo,
	}
}

func (s *ShadowServiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log := s.log.Named("get")
	pool := req.GetPool()
	log.Debug("Request received", zap.Strings("pool", pool))

	keys := make([]string, len(pool)*3)
	for i, dev := range pool {
		keys[i*3] = Key(dev, pb.StateKey_REPORTED)
		keys[i*3+1] = Key(dev, pb.StateKey_DESIRED)
		keys[i*3+2] = Key(dev, pb.StateKey_CONNECTION)
	}
	if len(keys) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no devices specified"))
	}
	r := s.rdb.MGet(ctx, keys...)
	states, err := r.Result()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get Shadows"))
	}

	log.Debug("Got states", zap.Int("count", len(states)))
	shadows := make([]*pb.Shadow, len(pool))
	for i := range shadows {
		s := &pb.Shadow{
			Device:     pool[i],
			Reported:   &pb.State{},
			Desired:    &pb.State{},
			Connection: &pb.ConnectionState{},
		}
		if states[i*3] != nil {
			state := states[i*3].(string)
			json.Unmarshal([]byte(state), s.Reported)
		}
		if states[i*3+1] != nil {
			state := states[i*3+1].(string)
			json.Unmarshal([]byte(state), s.Desired)
		}
		if states[i*3+2] != nil {
			state := states[i*3+2].(string)
			json.Unmarshal([]byte(state), s.Connection)
		}
		shadows[i] = s
	}

	return &pb.GetResponse{Shadows: shadows}, nil
}

func (s *ShadowServiceServer) Patch(ctx context.Context, req *pb.Shadow) (*pb.Shadow, error) {
	log := s.log.Named("patch")
	log.Debug("Request received", zap.Any("req", req))

	if req.GetDevice() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no device specified"))
	}

	now := timestamppb.Now()
	topics := []string{}
	if req.Reported != nil {
		req.Reported.Timestamp = now
		topics = append(topics, "mqtt.incoming")
	}
	if req.Desired != nil {
		req.Desired.Timestamp = now
		topics = append(topics, "mqtt.outgoing")
	}

	s.ps.TryPub(req, topics...)

	return req, nil
}

func (s *ShadowServiceServer) Remove(ctx context.Context, req *pb.RemoveRequest) (*pb.Shadow, error) {
	log := s.log.Named("remove")
	log.Debug("Request received", zap.Any("req", req))

	if req.GetDevice() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("no device specified"))
	}
	if req.GetKey() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("key not specified"))
	}

	skey := Key(req.GetDevice(), req.StateKey)
	r := s.rdb.Get(ctx, skey)
	raw, err := r.Result()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get Shadow"))
	}

	var state pb.State
	err = json.Unmarshal([]byte(raw), &state)
	if err != nil {
		log.Warn("Cannot unmarshal state", zap.String("raw", raw), zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("cannot Unmarshal state"))
	}

	fields := state.Data.Fields
	delete(fields, req.GetKey())

	state.Timestamp = timestamppb.Now()
	log.Debug("Result", zap.Any("state", state))

	s.Store(log, req.Device, req.StateKey, &state)

	result := &pb.Shadow{
		Device: req.GetDevice(),
	}

	if req.StateKey == pb.StateKey_REPORTED {
		result.Reported = &state
	} else {
		result.Desired = &state
	}

	return result, nil
}

func (s *ShadowServiceServer) StreamShadow(req *pb.StreamShadowRequest, srv pb.ShadowService_StreamShadowServer) (err error) {
	log := s.log.Named("stream")
	log.Debug("Request received", zap.Any("req", req))

	if len(req.GetDevices()) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("no devices specified"))
	}
	devices := map[string]bool{}
	for _, id := range req.GetDevices() {
		devices[id] = true
	}

	if req.Sync {
		func() {
			log.Debug("Sending current state")
			r, err := s.Get(srv.Context(), &pb.GetRequest{
				Pool: req.GetDevices(),
			})
			if err != nil {
				log.Warn("Couldn't get current devices Shadow state", zap.Error(err))
				return
			}
			for _, s := range r.GetShadows() {
				srv.Send(s)
			}
		}()
	}

	messages := make(chan interface{}, 10)
	s.ps.AddSub(messages, "mqtt.incoming", "mqtt.outgoing")
	defer unsub(s.ps, messages)

	log.Debug("Listening for messages")

	for msg := range messages {
		shadow := msg.(*pb.Shadow)
		if _, ok := devices[shadow.GetDevice()]; !ok {
			continue
		}
		err := srv.Send(shadow)
		if err != nil {
			log.Warn("Unable to send message", zap.Error(err))
			return nil
		}
	}

	return nil
}

func unsub[T chan any](ps pubsub.PubSub, ch chan any) {
	go ps.Unsub(ch)

	for range ch {
	}
}
