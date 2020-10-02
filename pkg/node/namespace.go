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

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

//NamespaceController is a Data type for Namespace Controller file
type NamespaceController struct {
	Repo Repo
	Log  *zap.Logger
}

var a *AccountController

//CreateNamespace is a method for creating Namespace
func (n *NamespaceController) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {

	log := n.Log.Named("Create Namespace Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace ", request.Name))

	id, err := n.Repo.CreateNamespace(ctx, request.GetName())
	if err != nil {
		//Added logging
		log.Error("Failed to create Namespace", zap.String("Name", request.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create Namespace")
	}

	//Added logging
	log.Info("Namespace Created", zap.String("Namespace ID", id))

	return &nodepb.Namespace{
		Id:   id,
		Name: request.GetName(),
	}, nil
}

//ListNamespaces is a method for Listing all the Namespaces
func (n *NamespaceController) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {

	log := n.Log.Named("List Namespaces Controller")
	//Added logging
	log.Info("Function Invoked")

	namespaces, err := n.Repo.ListNamespaces(ctx)
	if err != nil {
		//Added logging
		log.Error("Failed to list Namespaces", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("List Namespaces successful")
	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

//ListNamespacesForAccount is a method for Listing all the Namespaces for a specified account
func (n *NamespaceController) ListNamespacesForAccount(ctx context.Context, request *nodepb.ListNamespacesForAccountRequest) (response *nodepb.ListNamespacesResponse, err error) {

	log := n.Log.Named("List Namespaces for Account Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Account))

	namespaces, err := n.Repo.ListNamespacesForAccount(ctx, request.GetAccount())
	if err != nil {
		//Added logging
		log.Error("Failed to list Namespaces for the Account", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("List Namespaces for Account successful")
	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

//GetNamespace is a method to get details of a Namespace using Namespace name
func (n *NamespaceController) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {

	log := n.Log.Named("Get Namespace using name Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace", request.Namespace))

	namespace, err := n.Repo.GetNamespace(ctx, request.GetNamespace())
	if err != nil {
		//Added logging
		log.Error("Failed to get Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Get Namespace using name successful")
	return namespace, nil
}

//GetNamespaceID is a method to get details of a Namespace using NamespaceID
func (n *NamespaceController) GetNamespaceID(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {

	log := n.Log.Named("Get Namespace using ID Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace", request.Namespace))

	namespace, err := n.Repo.GetNamespaceID(ctx, request.GetNamespace())
	if err != nil {
		//Added logging
		log.Error("Failed to get Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Get Namespace using ID successful")
	return namespace, nil
}

//ListPermissions is a method to list all the accounts that have access to a Namespace
func (n *NamespaceController) ListPermissions(ctx context.Context, request *nodepb.ListPermissionsRequest) (response *nodepb.ListPermissionsResponse, err error) {

	log := n.Log.Named("List Permissions Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace", request.Namespace))

	permissions, err := n.Repo.ListPermissionsInNamespace(ctx, request.Namespace)
	if err != nil {
		//Added logging
		log.Error("Failed to list Permissions", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("List Permissions successful")
	return &nodepb.ListPermissionsResponse{Permissions: permissions}, nil
}

//DeletePermission is a method to delete access to a Namespace for a account
func (n *NamespaceController) DeletePermission(ctx context.Context, request *nodepb.DeletePermissionRequest) (response *nodepb.DeletePermissionResponse, err error) {

	log := n.Log.Named("Delete Permissions Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace", request.Namespace))

	err = n.Repo.DeletePermissionInNamespace(ctx, request.Namespace, request.AccountId)
	if err != nil {
		//Added logging
		log.Error("Failed to delete Permissions", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Delete Permission successful")
	return &nodepb.DeletePermissionResponse{}, nil
}

//DeleteNamespace is a method to delete a Namespace
func (n *NamespaceController) DeleteNamespace(ctx context.Context, request *nodepb.DeleteNamespaceRequest) (response *nodepb.DeleteNamespaceResponse, err error) {

	log := n.Log.Named("Delete Namespace Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Namespace", request.Namespaceid),
		zap.Bool("Hardelete Flag", request.Harddelete),
		zap.Bool("RevokeDelete Flag", request.Revokedelete),
	)

	if !request.Revokedelete {
		//Action to perform when delete is issued instead of revoke
		if request.Harddelete {
			//Set the datecondition to 14days back date
			//This is to ensure that records that are older then 14 days or more will be only be deleted.
			datecondition := time.Now().AddDate(0, 0, -14)

			//Added logging
			log.Info("Hard Delete Method Invoked")
			//Invokde Hardelete function with the date conidtion
			err = n.Repo.HardDeleteNamespace(ctx, datecondition.String())
			if err != nil {
				//Added logging
				log.Error("Failed to complete Hard delete Namespace process", zap.Error(err))
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			//Added logging
			log.Info("Soft Delete Method Invoked")
			//Soft delete will mark the record for deletion with the timestamp
			err = n.Repo.SoftDeleteNamespace(ctx, request.Namespaceid)
			if err != nil {
				//Added logging
				log.Error("Failed to Soft delete Namespace", zap.Error(err))
				return nil, status.Error(codes.Internal, err.Error())
			}
		}
	} else {
		//Action to perform when revoke is performed
		ns, err := n.Repo.GetNamespaceID(ctx, request.Namespaceid)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		//Added logging
		log.Info("Revoke Delete Method Invoked")
		//Initate Revoke
		if ns.Markfordeletion {
			err := n.Repo.RevokeNamespace(ctx, request.Namespaceid)
			if err != nil {
				//Added logging
				log.Error("Failed to Revoke delete Namespace", zap.Error(err))
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			//Added logging
			log.Error("Failed to Revoke as the Namespace is not marked for deletion")
			return nil, status.Error(codes.FailedPrecondition, "Failed to Revoke as the Namespace is not marked for deletion")
		}

	}

	if err != nil {
		//Added logging
		log.Error("Failed to delete Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to delete Namespace")
	}

	//Added logging
	log.Info("Delete Namespace successful")
	return &nodepb.DeleteNamespaceResponse{}, nil
}

//UpdateNamespace is a method to delete access to a Namespace for a account
func (n *NamespaceController) UpdateNamespace(ctx context.Context, request *nodepb.UpdateNamespaceRequest) (response *nodepb.UpdateNamespaceResponse, err error) {

	log := n.Log.Named("Update Namespace Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Namespace", request.Namespace.Id),
	)

	//Get the metadata from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//Added logging
		log.Error("Failed to get Account details", zap.Error(err))
		return nil, status.Error(codes.Aborted, "Failed to get Account details")
	}

	//Check for Authentication
	requestorID := md.Get("requestorID")
	if requestorID == nil {
		//Added logging
		log.Error("The Account is not authenticated", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	log.Info("Temp Logs", zap.Any("MD", md), zap.Any("Requestor ID", requestorID))

	//Check if the Account has WRITE access to Namespace
	resp, err := a.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:   requestorID[0],
		Namespace: request.Namespace.Id,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		//Added logging
		log.Error("Failed to get Authorization details for the Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.GetDecision().GetValue() {
		err = n.Repo.UpdateNamespace(ctx, request)
		if err != nil {
			//Added logging
			log.Error("Failed to update Namespace", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		//Added logging
		log.Error("The Account is not allowed to update the Namespace", zap.Error(err))
		return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to update the Namespace")
	}

	//Update logging
	log.Info("Delete Namespace successful")
	return &nodepb.UpdateNamespaceResponse{}, nil
}
