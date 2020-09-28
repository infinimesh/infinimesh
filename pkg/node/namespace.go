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

package node

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type NamespaceController struct {
	Repo Repo
}

//CreateNamespace is a method for creating Namespace
func (n *NamespaceController) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {
	id, err := n.Repo.CreateNamespace(ctx, request.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create namespace")
	}

	return &nodepb.Namespace{
		Id:   id,
		Name: request.GetName(),
	}, nil
}

//ListNamespaces is a method for Listing all the Namespaces
func (n *NamespaceController) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {
	namespaces, err := n.Repo.ListNamespaces(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed list namespaces")
	}

	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

//ListNamespacesForAccount is a method for Listing all the Namespaces for a specified account
func (n *NamespaceController) ListNamespacesForAccount(ctx context.Context, request *nodepb.ListNamespacesForAccountRequest) (response *nodepb.ListNamespacesResponse, err error) {
	namespaces, err := n.Repo.ListNamespacesForAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed list namespaces")
	}

	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

//GetNamespace is a method to get details of a Namespace using Namespace name
//Function to get the namespace based on namespace name
//Soon to be deprecated
func (n *NamespaceController) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {
	namespace, err := n.Repo.GetNamespace(ctx, request.GetNamespace())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed get namespace")
	}

	return namespace, nil
}

//GetNamespaceID is a method to get details of a Namespace using NamespaceID
//Function to get the namespace based on ID
func (n *NamespaceController) GetNamespaceID(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {
	namespace, err := n.Repo.GetNamespaceID(ctx, request.GetNamespace())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed get namespace")
	}

	return namespace, nil
}

//ListPermissions is a method to list all the accounts that have access to a Namespace
func (n *NamespaceController) ListPermissions(ctx context.Context, request *nodepb.ListPermissionsRequest) (response *nodepb.ListPermissionsResponse, err error) {
	permissions, err := n.Repo.ListPermissionsInNamespace(ctx, request.Namespace)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed get permissions")
	}

	return &nodepb.ListPermissionsResponse{Permissions: permissions}, nil
}

//DeletePermission is a method to delete access to a Namespace for a account
func (n *NamespaceController) DeletePermission(ctx context.Context, request *nodepb.DeletePermissionRequest) (response *nodepb.DeletePermissionResponse, err error) {
	err = n.Repo.DeletePermissionInNamespace(ctx, request.Namespace, request.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete permission")
	}
	return &nodepb.DeletePermissionResponse{}, nil
}

//DeleteNamespace is a method to delete a Namespace
func (n *NamespaceController) DeleteNamespace(ctx context.Context, request *nodepb.DeleteNamespaceRequest) (response *nodepb.DeleteNamespaceResponse, err error) {

	if !request.Revokedelete {
		//Action to perform when delete action is not revoked.
		if request.Harddelete {
			//Set the datecondition to 14days back date
			//This is to ensure that records that are older then 14 days or more will be only be deleted.
			datecondition := time.Now().AddDate(0, 0, -14)

			//Invokde Hardelete function with the date conidtion
			err = n.Repo.HardDeleteNamespace(ctx, datecondition.String())
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			//Soft delete will mark the record for deletion with the timestamp
			err = n.Repo.SoftDeleteNamespace(ctx, request.Namespaceid)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	} else {
		//Action to perform when delete action is revoked.
		ns, err := n.GetNamespaceID(ctx, &nodepb.GetNamespaceRequest{Namespace: request.Namespaceid})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if ns.Markfordeletion {
			//Add steps to update namespace
		} else {
			return nil, status.Error(codes.Internal, "The Namespace is not marked for deletion")
		}

	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete Namespace")
	}
	return &nodepb.DeleteNamespaceResponse{}, nil
}
