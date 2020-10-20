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

type namespaceAPI struct {
	client        nodepb.NamespacesClient
	accountClient nodepb.AccountServiceClient
}

//API Method to list all Namespaces
func (n *namespaceAPI) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {

	//Added logging
	log.Info("List Namespaces API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Create Namespace controller for server
	ns, err := n.client.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
	if err != nil {
		//Added logging
		log.Error("Create Namespace API Method: Failed to create Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("List Namespaces API Method: Namespace succesfully listed")
	return ns, nil
}

//API Method to create a Namespace
func (n *namespaceAPI) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {

	//Added logging
	log.Info("Create Namespace API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Create Namespace controller for server
	ns, err := n.client.CreateNamespace(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Create Namespace API Method: Failed to create Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Create Namespace API Method: Namespace succesfully created and Account auhtorized to access it")
	return ns, nil
}

//API Method to get details of a Namespace
func (n *namespaceAPI) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {

	//Added logging
	log.Info("Get Namespace API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Get Namespace controller for server
	ns, err := n.client.GetNamespace(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Get Namespace API Method: Failed to Get Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Get Namespace API Method: Namespace details succesfully obtained")
	return ns, nil
}

//API Method to create permissions for ana ccount to access a Namespace
func (n *namespaceAPI) CreatePermission(ctx context.Context, request *apipb.CreateNamespacePermissionRequest) (response *apipb.CreateNamespacePermissionResponse, err error) {

	//Added logging
	log.Info("Create Permission API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Create Permission API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     account,
		Namespaceid: request.Namespace,
		Action:      nodepb.Action_WRITE,
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
	log.Info("List Permissions API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("List Permissions API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     account,
		Namespaceid: request.Namespace,
		Action:      nodepb.Action_WRITE,
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
	log.Info("Delete Permission API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Delete Permission API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	resp, err := n.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     account,
		Namespaceid: request.Namespace,
		Action:      nodepb.Action_WRITE,
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
	log.Info("Delete Namespace API Method: Function Invoked",
		zap.String("Requestor ID", ctx.Value("account_id").(string)),
		zap.String("Namespace", request.Namespaceid),
		zap.Bool("HardDelete Flag", request.Harddelete),
	)

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Update Namespace controller for server
	ns, err := n.client.DeleteNamespace(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Delete Namespace API Method: Failed to delete Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Delete Namespace API Method: Namespace succesfully marked for deletion")
	return ns, nil

}

//API Method to Update namespace
func (n *namespaceAPI) UpdateNamespace(ctx context.Context, request *nodepb.UpdateNamespaceRequest) (response *nodepb.UpdateNamespaceResponse, err error) {

	//Added logging
	log.Info("Update Namespace API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Update Namespace controller for server
	ns, err := n.client.UpdateNamespace(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Update Namespace API Method: Failed to update Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Update Namespace API Method: Namespace succesfully updated")
	return ns, nil
}
