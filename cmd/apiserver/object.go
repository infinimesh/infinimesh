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
)

type objectAPI struct {
	objectClient  nodepb.ObjectServiceClient
	accountClient nodepb.AccountServiceClient
}

func (o *objectAPI) CreateObject(ctx context.Context, request *apipb.CreateObjectRequest) (response *nodepb.Object, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	if request.Object == nil || request.Object.Name == "" {
		return nil, status.Error(codes.FailedPrecondition, "Invalid object given.")
	}

	// If a parent is given, we need permission on the parent. otherwise, we need permission on the namespace as it's created without a parent
	var authorized bool
	var parent string
	if request.Parent != nil {
		parent = request.Parent.Value
		resp, err := o.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
			Node:    request.Parent.GetValue(),
			Account: account,
			Action:  nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, err
		}
		authorized = resp.Decision.GetValue()
	} else {
		resp, err := o.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
			Namespaceid: request.Namespace,
			Account:     account,
			Action:      nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, err
		}

		authorized = resp.Decision.GetValue()
	}

	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "The account does not have permission to create object.")
	}

	return o.objectClient.CreateObject(ctx, &nodepb.CreateObjectRequest{
		Parent:      parent,
		Name:        request.Object.Name,
		Namespaceid: request.Namespace,
		Kind:        request.Object.Kind,
	})
}

func (o *objectAPI) ListObjects(ctx context.Context, request *apipb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	fmt.Println("rec?", request.Recurse)

	// This request automatically runs in the scope of the user, no need to call IsAuthorized
	return o.objectClient.ListObjects(ctx, &nodepb.ListObjectsRequest{Account: account, Namespace: request.GetNamespace(), Recurse: request.Recurse})
}

func (o *objectAPI) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	resp, err := o.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.GetUid(),
		Account: account,
		Action:  nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Decision.GetValue() {
		return nil, status.Error(codes.PermissionDenied, "The account does not have permission to access the resource.")
	}

	return o.objectClient.DeleteObject(ctx, request)
}
