//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
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
)

type deviceAPI struct {
	client        registrypb.DevicesClient
	accountClient nodepb.AccountServiceClient
}

//API Method to Create a Device
func (d *deviceAPI) Create(ctx context.Context, request *registrypb.CreateRequest) (response *registrypb.CreateResponse, err error) {

	//Added logging
	log.Info("Create Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Create Device controller for server
	res, err := d.client.Create(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Create Device API Method: Failed to create Device", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Create Device API Method: Device succesfully created")
	return res, nil
}

//API Method to Update a Device
func (d *deviceAPI) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {

	//Added logging
	log.Info("Update Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Update Device controller for server
	res, err := d.client.Update(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Update Device API Method: Failed to update Device", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Update Device API Method: Device succesfully updated")
	return res, nil
}

//API Method to Get a Device.
func (d *deviceAPI) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {

	//Added logging
	log.Debug("Get Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Get Device controller for server
	res, err := d.client.Get(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Get Device API Method: Failed to get Device", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Debug("Get Device API Method: Device Details succesfully obtained")
	return res, nil

}

//API Method to List a Device
func (d *deviceAPI) List(ctx context.Context, request *apipb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {

	//Added logging
	log.Debug("List Devices API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the List Device controller for server
	list, err := d.client.List(ctx, &registrypb.ListDevicesRequest{Namespaceid: request.Namespaceid, Account: ctx.Value("account_id").(string)})
	if err != nil {
		//Added logging
		log.Error("List Device API Method: Failed to list Devices", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Debug("List Device API Method: Device succesfully listed")
	return list, nil
}

//API Method to Delete a Device
func (d *deviceAPI) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {

	//Added logging
	log.Info("Delete Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Delete Device controller for server
	resp, err := d.client.Delete(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Delete Device API Method: Failed to list Devices", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Delete Device API Method: Device deleted successfully")
	return resp, nil
}

//API Method to Assign a Owner to a Device
func (d *deviceAPI) AssignOwnerDevices(ctx context.Context, request *registrypb.OwnershipRequestDevices) (response *registrypb.OwnershipResponseDevices, err error) {

	//Added logging
	log.Info("Assign Owner Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Delete Device controller for server
	resp, err := d.client.AssignOwnerDevices(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Assign Owner Device API Method: Failed to list Devices", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Assign Owner Device API Method: Owner Assigned successfully")
	return resp, nil
}

//API Method to Remove a Owner to a Device
func (d *deviceAPI) RemoveOwnerDevices(ctx context.Context, request *registrypb.OwnershipRequestDevices) (response *registrypb.OwnershipResponseDevices, err error) {

	//Added logging
	log.Info("Remove Owner Device API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Delete Device controller for server
	resp, err := d.client.RemoveOwnerDevices(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Remove Owner Device API Method: Failed to list Devices", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Remove Owner Device API Method: Owner Removed successfully")
	return resp, nil
}
