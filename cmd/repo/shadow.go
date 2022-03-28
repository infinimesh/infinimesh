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

// TO BE DEPRECATED AND MOVED TO pkg/state

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
)

//ShadowAPI data strcuture
type ShadowAPI struct {
	pb.UnimplementedShadowServiceServer

	log *zap.Logger
	client        shadowpb.ShadowsClient
}

func NewShadowAPI(log *zap.Logger, client shadowpb.ShadowsClient) *ShadowAPI {
	return &ShadowAPI{
		log: log.Named("ShadowAPI"), client: client,
	}
}	

//Get is a method to get the current state of the device
func (s *ShadowAPI) Get(ctx context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	log := s.log.Named("Get")
	log.Debug("Get request received", zap.Any("request", request), zap.Any("context", ctx))

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	found := false
	for _, device := range devices_scope {
		if device == request.Id {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}

	return s.client.Get(ctx, request)
}

func (s *ShadowAPI) GetMultiple(ctx context.Context, _ *shadowpb.Empty) (response *shadowpb.GetMultipleResponse, err error) {
	log := s.log.Named("GetMultiple")

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	type State struct {
		OK bool
		ID string
		State *shadowpb.Shadow
	}
	n := len(devices_scope)
	states := make(chan State, n)
	defer close(states)

	log.Debug("Gathering devices states", zap.Int("amount", n))
	for _, dev := range devices_scope {
		go func(dev string, r chan State) {
			res, err := s.client.Get(ctx, &shadowpb.GetRequest{Id: dev})
			if err != nil {
				log.Error("Error getting Device state", zap.Error(err))
				r <- State{OK: false}
			}
			r <- State{true, dev, res.GetShadow()}
		}(dev, states)
	}

	result := make(map[string]*shadowpb.Shadow)
	for i := 0; i < n; i++ {
		state := <- states
		if state.OK {
			result[state.ID] = state.State
		}
	}

	return &shadowpb.GetMultipleResponse{Pool: result}, nil
}

//PatchDesiredState is a method to update the current state of the device
func (s *ShadowAPI) PatchDesiredState(ctx context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {
	log := s.log.Named("PatchDesiredState")

	post_allowed, ok := ctx.Value(inf.InfinimeshPostAllowedCtxKey).(bool)
	if !ok || !post_allowed {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope or not allowed to post")
	}

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope or not allowed to post")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	found := false
	for _, device := range devices_scope {
		if device == request.Id {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}

	return s.client.PatchDesiredState(ctx, request)
}

//StreamReportedStateChanges is a method to get the stream for a device
func (s *ShadowAPI) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv pb.ShadowService_StreamReportedStateChangesServer) (err error) {
	log := s.log.Named("StreamReportedStateChanges")

	devices_scope, ok := srv.Context().Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return status.Error(codes.Unauthenticated, "Requested device is outside of token scope or not allowed to post")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	found := false
	for _, device := range devices_scope {
		if device == request.Id {
			found = true
			break
		}
	}
	if !found {
		return status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}

	log.Debug("Stream API Method: Streaming started", zap.String("Device ID", request.Id))

	c, err := s.client.StreamReportedStateChanges(srv.Context(), request)
	if err != nil {
		log.Error("Stream API Method: Failed to start the Stream", zap.Error(err))
		return status.Error(codes.Unauthenticated, "Failed to start the Stream")
	}

	for {
		msg, err := c.Recv()
		if err != nil {
			log.Error("Stream API Method: Error while receving message", zap.Error(err))
			return err
		}

		err = srv.Send(msg)
		if err != nil {
			log.Error("Stream API Method: Error while sending message", zap.Error(err))
			return err
		}
	}
}