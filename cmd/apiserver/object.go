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
)

type objectAPI struct {
	objectClient  nodepb.ObjectServiceClient
	accountClient nodepb.AccountServiceClient
}

func (o *objectAPI) CreateObject(ctx context.Context, request *apipb.CreateObjectRequest) (response *nodepb.Object, err error) {

	//Added logging
	log.Info("Create Object API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	var parent string
	if request.Parent != nil {
		parent = request.Parent.Value
	}

	//Invoke the List Objects controller for server
	obj, err := o.objectClient.CreateObject(ctx, &nodepb.CreateObjectRequest{
		Parent:      parent,
		Name:        request.Object.Name,
		Namespaceid: request.Namespace,
		Kind:        request.Object.Kind,
	})

	if err != nil {
		//Added logging
		log.Error("Create Object API Method: Failed to creat Object", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Create Object API Method: Create Object succesful")
	return obj, nil

}

func (o *objectAPI) ListObjects(ctx context.Context, request *apipb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {

	//Added logging
	log.Debug("List Objects API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the List Objects controller for server
	obj, err := o.objectClient.ListObjects(ctx, &nodepb.ListObjectsRequest{Account: ctx.Value("account_id").(string), Namespace: request.GetNamespace(), Recurse: request.Recurse})
	if err != nil {
		//Added logging
		log.Error("List Objects API Method: Failed to list Objects", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Debug("List Objects API Method: List Objects succesfull")
	return obj, nil
}

func (o *objectAPI) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {

	//Added logging
	log.Info("Delete Object API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Delete Object controller for server
	obj, err := o.objectClient.DeleteObject(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Delete Object API Method: Failed to delete Object", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Delete Object API Method: Delete Object succesfull")
	return obj, nil
}
