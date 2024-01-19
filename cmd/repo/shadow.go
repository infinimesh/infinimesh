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

// TO BE DEPRECATED AND MOVED TO pkg/state

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"go.uber.org/zap"

	inf "github.com/infinimesh/infinimesh/pkg/shared"
	pb "github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/shadow"
)

// ShadowAPI data strcuture
type ShadowAPI struct {
	pb.UnimplementedShadowServiceServer

	log    *zap.Logger
	client shadow.ShadowServiceClient
}

func NewShadowAPI(log *zap.Logger, client shadow.ShadowServiceClient) *ShadowAPI {
	return &ShadowAPI{
		log: log.Named("ShadowAPI"), client: client,
	}
}

func (s *ShadowAPI) Get(ctx context.Context, _ *connect.Request[shadow.GetRequest]) (response *connect.Response[shadow.GetResponse], err error) {
	log := s.log.Named("Get")

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).(map[string]any)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope or not allowed to post"))
	}
	log.Debug("Scope", zap.Any("devices", devices_scope))

	var pool = make([]string, len(devices_scope))
	for device := range devices_scope {
		pool = append(pool, device)
	}

	res, err := s.client.Get(ctx, &shadow.GetRequest{Pool: pool})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

// PatchDesiredState is a method to update the current state of the device
func (s *ShadowAPI) Patch(ctx context.Context, request *connect.Request[shadow.Shadow]) (response *connect.Response[shadow.Shadow], err error) {
	log := s.log.Named("PatchDesiredState")
	shadow := request.Msg

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).(map[string]any)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope or not allowed to post"))
	}
	log.Debug("Scope", zap.Any("devices", devices_scope))

	found := false
	for device := range devices_scope {
		if device == shadow.Device {
			found = true
			break
		}
	}
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope"))
	}

	res, err := s.client.Patch(ctx, shadow)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

func (s *ShadowAPI) Remove(ctx context.Context, request *connect.Request[shadow.RemoveRequest]) (response *connect.Response[shadow.Shadow], err error) {
	log := s.log.Named("RemoveStateKey")
	req := request.Msg

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).(map[string]any)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope or not allowed to post"))
	}
	log.Debug("Scope", zap.Any("devices", devices_scope))

	found := false
	for device := range devices_scope {
		if device == req.Device {
			found = true
			break
		}
	}
	if !found {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope"))
	}

	res, err := s.client.Remove(ctx, req)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}

// StreamShadow is a method to get the stream for a device
func (s *ShadowAPI) StreamShadow(ctx context.Context, request *connect.Request[shadow.StreamShadowRequest], srv *connect.ServerStream[shadow.Shadow]) (err error) {
	log := s.log.Named("StreamReportedStateChanges")
	req := request.Msg

	devices_scope, ok := ctx.Value(inf.InfinimeshDevicesCtxKey).(map[string]any)
	if !ok {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("requested device is outside of token scope or not allowed to post"))
	}
	log.Debug("Scope", zap.Any("devices", devices_scope))

	var pool = make([]string, len(devices_scope))
	for device := range devices_scope {
		pool = append(pool, device)
	}

	req.Devices = pool

	log.Debug("Stream API Method: Streaming started", zap.Strings("devices", pool))

	c, err := s.client.StreamShadow(ctx, req)
	if err != nil {
		log.Warn("Stream API Method: Failed to start the Stream", zap.Error(err))
		return connect.NewError(connect.CodeAborted, errors.New("failed to start the Stream"))
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

func (s *ShadowAPI) StreamShadowSync(ctx context.Context, request *connect.Request[shadow.StreamShadowRequest], srv *connect.ServerStream[shadow.Shadow]) (err error) {
	request.Msg.Sync = true
	return s.StreamShadow(ctx, request, srv)
}
