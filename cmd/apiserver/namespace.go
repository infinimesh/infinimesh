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
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type namespaceAPI struct {
	client        nodepb.NamespacesClient
	accountClient nodepb.AccountServiceClient
}

func (n *namespaceAPI) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check permissions")
	}

	if resp.GetIsRoot() {
		return n.client.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
	} else {
		return n.client.ListNamespacesForAccount(ctx, &nodepb.ListNamespacesForAccountRequest{
			Account: account,
		})
	}
}

func (n *namespaceAPI) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check permissions")
	}

	if resp.GetIsRoot() {
		// TODO this is not atomic and if the application crashes
		// between both calls, we'll have a problem. Maybe move it to
		// one operation into the repo, and do within a txn.
		_, err = n.client.CreateNamespace(ctx, request)
		if err != nil {
			return nil, err
		}

		_, err := n.accountClient.AuthorizeNamespace(ctx, &nodepb.AuthorizeNamespaceRequest{
			Account:   account,
			Namespace: request.GetName(),
			Action:    nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to authorize after creating ns")
		}
		return &nodepb.Namespace{}, nil

	}
	return nil, status.Error(codes.PermissionDenied, "Account is not root")
}

func (n *namespaceAPI) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   account,
		Namespace: request.GetNamespace(),
		Action:    nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		return n.client.GetNamespace(ctx, request)
	}
	return nil, status.Error(codes.PermissionDenied, "Account is not allowed to access this resource")
}

func (n *namespaceAPI) CreatePermission(ctx context.Context, request *apipb.CreateNamespacePermissionRequest) (response *apipb.CreateNamespacePermissionResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   account,
		Namespace: request.Namespace,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		_, err := n.accountClient.AuthorizeNamespace(ctx, &nodepb.AuthorizeNamespaceRequest{
			Account:   request.AccountId,
			Namespace: request.Namespace,
			Action:    request.Permission.Action,
		})
		if err != nil {
			log.Error("Failed to create permission", zap.String("account_id", request.AccountId), zap.String("namespace", request.Namespace), zap.Error(err))
			return &apipb.CreateNamespacePermissionResponse{}, status.Error(codes.Internal, "Failed to authorize for namespace")
		}
		return &apipb.CreateNamespacePermissionResponse{}, nil
	}

	return nil, status.Error(codes.PermissionDenied, "Account is not allowed to access this resource")
}

func (n *namespaceAPI) ListPermissions(ctx context.Context, request *nodepb.ListPermissionsRequest) (response *nodepb.ListPermissionsResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   account,
		Namespace: request.Namespace,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		return n.client.ListPermissions(ctx, request)
	}

	return nil, status.Error(codes.PermissionDenied, "Account is not allowed to access this resource")

}

func (n *namespaceAPI) DeletePermission(ctx context.Context, request *nodepb.DeletePermissionRequest) (response *nodepb.DeletePermissionResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   account,
		Namespace: request.Namespace,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		return n.client.DeletePermission(ctx, request)
	}

	return nil, status.Error(codes.PermissionDenied, "Account is not allowed to access this resource")

}
