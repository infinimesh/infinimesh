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

	inf "github.com/infinimesh/infinimesh/pkg/shared"
	pb "github.com/infinimesh/proto/node"
	shadowpb "github.com/infinimesh/proto/shadow"
)

//ShadowAPI data strcuture
type ShadowAPI struct {
	pb.UnimplementedShadowServiceServer

	log    *zap.Logger
	client shadowpb.ShadowServiceClient
}

func NewShadowAPI(log *zap.Logger, client shadowpb.ShadowServiceClient) *ShadowAPI {
	return &ShadowAPI{
		log: log.Named("ShadowAPI"), client: client,
	}
}

func (s *ShadowAPI) Get(ctx context.Context, _ *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	log := s.log.Named("Get")

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	return s.client.Get(ctx, &shadowpb.GetRequest{Pool: devices_scope})
}

//PatchDesiredState is a method to update the current state of the device
func (s *ShadowAPI) Patch(ctx context.Context, request *shadowpb.Shadow) (response *shadowpb.Shadow, err error) {
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
		if device == request.Device {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.Unauthenticated, "Requested device is outside of token scope")
	}

	return s.client.Patch(ctx, request)
}

//StreamShadow is a method to get the stream for a device
func (s *ShadowAPI) StreamShadow(request *shadowpb.StreamShadowRequest, srv pb.ShadowService_StreamShadowServer) (err error) {
	log := s.log.Named("StreamReportedStateChanges")

	devices_scope, ok := srv.Context().Value(inf.InfinimeshDevicesCtxKey).([]string)
	if !ok {
		return status.Error(codes.Unauthenticated, "Requested device is outside of token scope or not allowed to post")
	}
	log.Debug("Scope", zap.Strings("devices", devices_scope))

	request.Devices = devices_scope

	log.Debug("Stream API Method: Streaming started", zap.Strings("devices", devices_scope))

	c, err := s.client.StreamShadow(srv.Context(), request)
	if err != nil {
		log.Warn("Stream API Method: Failed to start the Stream", zap.Error(err))
		return status.Error(codes.Unauthenticated, "Failed to start the Stream")
	}

	for {
		msg, err := c.Recv()
		if err != nil {
			log.Info("Error receiving message, closing stream", zap.Error(err))
			return err
		}

		err = srv.Send(msg)
		if err != nil {
			log.Info("Error sending message, closing stream", zap.Error(err))
			return err
		}
	}
}
