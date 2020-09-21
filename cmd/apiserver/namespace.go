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

//API Method to list all Namespaces
func (n *namespaceAPI) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {

	//Added logging
	log.Info("List Namespaces API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("List Namespaces API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		//Added logging
		log.Error("List Namespaces API Method: Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	if resp.GetIsRoot() {
		return n.client.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
	} else {
		return n.client.ListNamespacesForAccount(ctx, &nodepb.ListNamespacesForAccountRequest{
			Account: account,
		})
	}
}

//API Method to create a Namespace
func (n *namespaceAPI) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {

	//Added logging
	log.Info("Create Namespace API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Create Namespace API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		//Added logging
		log.Error("Create Namespace API Method: Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	if resp.GetIsRoot() {
		// TODO this is not atomic and if the application crashes
		// between both calls, we'll have a problem. Maybe move it to
		// one operation into the repo, and do within a txn.
		ns, err := n.client.CreateNamespace(ctx, request)
		if err != nil {
			return nil, err
		}

		//Assign Permissions to the namespace to the account that was used to create namespace
		_, err = n.accountClient.AuthorizeNamespace(ctx, &nodepb.AuthorizeNamespaceRequest{
			Account:   account,
			Namespace: ns.GetId(),
			Action:    nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to authorize after creating ns")
		}

		//Added logging
		log.Info("Create Namespace API Method: Namespace Created", zap.String("Namespace ID", ns.Id))
		return &nodepb.Namespace{
			Id:   ns.Id,
			Name: ns.Name,
		}, nil

	}

	//Added logging
	log.Error("Create Namespace API Method: The Account does not have permission to create Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account does not have permission to create Namespace")
}

//API Method to get details of a Namespace
func (n *namespaceAPI) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {

	//Added logging
	log.Info("Get Namespace API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Get Namespace API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
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

	//Added logging
	log.Error("Get Namespace API Method: The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to access the Namespace")
}

//API Method to create permissions for ana ccount to access a Namespace
func (n *namespaceAPI) CreatePermission(ctx context.Context, request *apipb.CreateNamespacePermissionRequest) (response *apipb.CreateNamespacePermissionResponse, err error) {

	//Added logging
	log.Info("Create Permission API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Create Permission API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
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
			return &apipb.CreateNamespacePermissionResponse{}, status.Error(codes.Internal, err.Error())
		}
		return &apipb.CreateNamespacePermissionResponse{}, nil
	}

	//Added logging
	log.Error("Create Permission API Method: The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to access the Namespace")
}

//API Method to list all the permissions for a Namespace
func (n *namespaceAPI) ListPermissions(ctx context.Context, request *nodepb.ListPermissionsRequest) (response *nodepb.ListPermissionsResponse, err error) {

	//Added logging
	log.Info("List Permissions API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("List Permissions API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
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

	//Added logging
	log.Error("List Permissions API Method: The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The account is not allowed to access the Namespace")
}

//API Method to delete permissions for am account for a Namespace
func (n *namespaceAPI) DeletePermission(ctx context.Context, request *nodepb.DeletePermissionRequest) (response *nodepb.DeletePermissionResponse, err error) {

	//Added logging
	log.Info("Delete Permission API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Delete Permission API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
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

	//Added logging
	log.Error("Delete Permission API Method: The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to access the Namespace")

}

//API Method to delete namespace
func (n *namespaceAPI) DeleteNamespace(ctx context.Context, request *nodepb.DeleteNamespaceRequest) (response *nodepb.DeleteNamespaceResponse, err error) {

	//Added logging
	log.Info("Delete Namespace API Method: Function Invoked", zap.String("Account ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Delete Namespace API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   account,
		Namespace: request.Namespaceid,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		return n.client.DeleteNamespace(ctx, request)
	}

	//Added logging
	log.Error("Delete Namespace API Method: The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to access the Namespace")

}
