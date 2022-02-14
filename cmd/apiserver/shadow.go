//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package main

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

//shadowAPI data strcuture
type shadowAPI struct {
	apipb.UnimplementedStatesServer

	accountClient nodepb.AccountServiceClient
	devicesClient registrypb.DevicesClient
	client        shadowpb.ShadowsClient
}

//Get is a method to get the current state of the device
func (s *shadowAPI) Get(ctx context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {

	//Added logging
	log.Debug("Get State API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := s.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return s.client.Get(ctx, request)
}

func (s *shadowAPI) GetForNS(ctx context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetForNSResponse, err error) {
	log.Debug("Get State for Namespace API Method: Invoked", zap.String("namespace", request.GetId()))

	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Namespace ID not given")
	}

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	log.Debug("Requestor ID", zap.String("account", account))

	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", account)
	res, err := s.devicesClient.List(ctx, &registrypb.ListDevicesRequest{Namespaceid: request.GetId(), Account: account})
	if err != nil {
		log.Error("Error listing Devices for Namespace",
			zap.String("namespace", request.GetId()), zap.Error(err))
		return nil, err
	}

	type State struct {
		OK bool
		ID string
		State *shadowpb.Shadow
	}
	n := len(res.GetDevices())
	states := make(chan State, n)
	defer close(states)

	log.Debug("Gathering devices states", zap.Int("amount", n))
	for _, dev := range res.GetDevices() {
		go func(dev *registrypb.Device, r chan State) {
			res, err := s.client.Get(ctx, &shadowpb.GetRequest{Id: dev.GetId()})
			if err != nil {
				log.Error("Error getting Device state", zap.Error(err))
				r <- State{OK: false}
			}
			r <- State{true, dev.GetId(), res.GetShadow()}
		}(dev, states)
	}

	result := make(map[string]*shadowpb.Shadow)
	for i := 0; i < n; i++ {
		state := <- states
		if state.OK {
			result[state.ID] = state.State
		}
	}

	return &shadowpb.GetForNSResponse{Pool: result}, nil
}

//PatchDesiredState is a method to update the current state of the device
func (s *shadowAPI) PatchDesiredState(ctx context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {

	//Added logging
	log.Debug("Patch Desired State API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := s.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return s.client.PatchDesiredState(ctx, request)
}

//StreamReportedStateChanges is a method to get the stream for a device
func (s *shadowAPI) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv apipb.States_StreamReportedStateChangesServer) (err error) {

	//Added logging
	log.Debug("Stream API Method: Function Invoked", zap.String("Requestor ID", srv.Context().Value("account_id").(string)))

	account, ok := srv.Context().Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Stream API Method: The Account is not authenticated")
		return status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := s.accountClient.IsAuthorized(srv.Context(), &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return status.Error(codes.PermissionDenied, "Stream API Method: Failed to get Authorization for the Stream")
	}
	if !resp.GetDecision().GetValue() {
		return status.Error(codes.PermissionDenied, "Stream API Method: The account doesnot have permission to start the Stream")
	}

	//Added logging
	log.Debug("Stream API Method: Streaming started", zap.String("Device ID", request.Id))

	log.Debug("Stream API Method: Temporary Logs1", zap.Any("Context", srv.Context))
	c, err := s.client.StreamReportedStateChanges(srv.Context(), request)
	if err != nil {
		//Added logging
		log.Error("Stream API Method: Failed to start the Stream", zap.Error(err))
		return status.Error(codes.Unauthenticated, "Failed to start the Stream")
	}
	log.Debug("Stream API Method: Temporary Logs2", zap.Any("Context", srv.Context))

	for {
		msg, err := c.Recv()
		log.Debug("Stream API Method: Temporary Logs3", zap.Any("Message", msg))
		if err != nil {
			//Added logging
			log.Error("Stream API Method: Error while receving message", zap.Error(err))
			return err
		}

		err = srv.Send(msg)
		if err != nil {
			//Added logging
			log.Error("Stream API Method: Error while sending message", zap.Error(err))
			return err
		}
	}
}
