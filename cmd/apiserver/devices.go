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
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

type deviceAPI struct {
	client        registrypb.DevicesClient
	accountClient nodepb.AccountServiceClient
}

func (d *deviceAPI) Create(ctx context.Context, request *registrypb.CreateRequest) (response *registrypb.CreateResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := d.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Namespace: request.Device.Namespace,
		Account:   account,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Could not get permission to create device.")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied to create device.")
	}
	return d.client.Create(ctx, request)
}

func (d *deviceAPI) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := d.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Device.Id,
		Account: account,
		Action:  nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return d.client.Update(ctx, request)
}

func (d *deviceAPI) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := d.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	fmt.Println("decision", resp.Decision.Value)
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return d.client.Get(ctx, request)

}
func (d *deviceAPI) List(ctx context.Context, request *apipb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isRootResp, err := d.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	})
	if err != nil {
		return nil, err
	}

	if isRootResp.IsRoot {
		return d.client.List(ctx, &registrypb.ListDevicesRequest{Namespaceid: request.Namespace})
	}

	resp, err := d.client.ListForAccount(ctx, &registrypb.ListDevicesRequest{Namespaceid: request.Namespace, Account: account})
	return resp, err
}
func (d *deviceAPI) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	return d.client.Delete(ctx, request)
}
