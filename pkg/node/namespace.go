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
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

//NamespaceController is a Data type for Namespace Controller file
type NamespaceController struct {
	Repo Repo
	Log  *zap.Logger
}

//Account controller to get access to method for Account validation
var a AccountController

//CreateNamespace is a method for creating Namespace
func (n *NamespaceController) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {

	log := n.Log.Named("Create Namespace Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Namespace ", request.Name))

	//Initialize the Account Controller with Namespace controller data
	a.Repo = n.Repo
	a.Log = n.Log

	//Get metadata from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//Validated that required data is populated with values
	if request.Name == "" {
		//Added logging
		log.Error("Data Validation for Namespace Creation", zap.String("Error", "The Name cannot not be empty"))
		return nil, status.Error(codes.FailedPrecondition, "The Name cannot not be empty")
	}

	//Check if the account is root
	isroot, err := a.IsRoot(ctx, &nodepb.IsRootRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	//Check if the account is admin
	isadmin, err := a.IsAdmin(ctx, &nodepb.IsAdminRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	var id string
	//Create the namespace if the account is root or admin
	if isroot.GetIsRoot() || isadmin.GetIsAdmin() {
		log.Info("Create Namespace initiated")
		id, err = n.Repo.CreateNamespace(ctx, request.GetName())
		if err != nil {
			//Added logging
			log.Error("Failed to create Namespace", zap.String("Name", request.Name), zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to create Namespace")
		}
	} else {
		//Added logging
		log.Error("The Account does not have permission to create Namespace")
		return &nodepb.Namespace{}, status.Error(codes.PermissionDenied, "The Account does not have permission to create Namespace")
	}

	//Added logging
	log.Info("Namespace Created", zap.String("Namespace ID", id), zap.String("Namespace Name", request.GetName()))

	//Assign Permissions to the account that was used to create namespace
	_, err = a.AuthorizeNamespace(ctx, &nodepb.AuthorizeNamespaceRequest{
		Account:   requestorID,
		Namespace: id,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to assign permissions to the Account for the Namespace")
	}

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

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//Initialize the Account Controller with Namespace controller data
	a.Repo = n.Repo
	a.Log = n.Log

	//Check if the account is root
	isroot, err := a.IsRoot(ctx, &nodepb.IsRootRequest{Account: requestorID})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var namespaces []*nodepb.Namespace
	if isroot.GetIsRoot() {
		//Get the namespaces for root account
		namespaces, err = n.Repo.ListNamespaces(ctx)
		if err != nil {
			//Added logging
			log.Error("Failed to list Namespaces", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to list Namespaces")
		}
	} else {
		//Check is the account is present
		_, err := n.Repo.UserExists(ctx, requestorID)
		if err != nil {
			//Added logging
			log.Error("Failed to list Namespaces for the Account", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to list Namespaces for the Account")
		}

		//Get the namespaces for a specific account
		namespaces, err = n.Repo.ListNamespacesForAccount(ctx, requestorID)
		if err != nil {
			//Added logging
			log.Error("Failed to list Namespaces for the Account", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to list Namespaces for the Account")
		}
	}

	//Added logging
	log.Info("List Namespaces successful")
	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

//GetNamespace is a method to get details of a Namespace using Namespace name
func (n *NamespaceController) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {

	log := n.Log.Named("Get Namespace using name Controller")

	//Added logging
	log.Info("Function Invoked", zap.String("Namespace", request.Namespace))

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//This is the only way to get namespace id from name and then send the namespace id for for Authorizaion check
	namespace, err := n.Repo.GetNamespace(ctx, request.GetNamespace())
	if err != nil {
		//Added logging
		log.Error("Failed to get Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Initialize the Account Controller with Namespace controller data
	a.Repo = n.Repo
	a.Log = n.Log

	//Check if the account has access to Namespace
	resp, err := a.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     requestorID,
		Namespaceid: namespace.Id,
		Action:      nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to Authorize the user"+err.Error())
	}

	if resp.GetDecision().GetValue() {

		//Added logging
		log.Info("Get Namespace using name successful")
		return namespace, nil
	}

	//Added logging
	log.Error("The Account is not allowed to access the Namespace")
	return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to access the Namespace")
}

//GetNamespaceID is a method to get details of a Namespace using Namespace ID
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
	)

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//Initialize the Account Controller with Namespace controller data
	a.Repo = n.Repo
	a.Log = n.Log

	//Check if the account is root
	isroot, err := a.IsRoot(ctx, &nodepb.IsRootRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the Account")
	}

	//Check if the account is admin
	isadmin, err := a.IsAdmin(ctx, &nodepb.IsAdminRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the Account")
	}

	namespace, err := n.Repo.GetNamespaceID(ctx, request.Namespaceid)
	if err != nil {
		//Added logging
		log.Error("Failed to get Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Validate that namespace is not root
	if namespace.Name == "root" {
		//Added logging
		log.Error("Cannot delete root Namespace")
		return nil, status.Error(codes.FailedPrecondition, "Cannot delete root Namespace")
	}

	//Check if the Account has WRITE access to Namespace
	resp, err := a.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     requestorID,
		Namespaceid: request.Namespaceid,
		Action:      nodepb.Action_WRITE,
	})
	if err != nil {
		//Added logging
		log.Error("Failed to get Authorization details for the Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Initiate delete if the account has access
	if resp.GetDecision().GetValue() && (isroot.GetIsRoot() || isadmin.GetIsAdmin()) {
		//Action to perform when delete is issued instead of revoke
		if request.Harddelete {
			//Set the datecondition to 14days back date
			//This is to ensure that records that are older then 14 days or more will be only be deleted.
			datecondition := time.Now().AddDate(0, 0, -14).Format(time.RFC3339)

			//Added logging
			log.Info("Hard Delete Process Invoked")
			//Invokde Hardelete function with the date conidtion
			err = n.Repo.HardDeleteNamespace(ctx, datecondition)
			if err != nil {
				if status.Code(err) != 5 { //5 is the error code for NotFound in GRPC
					//Added logging
					log.Error("Failed to complete Hard delete Namespace process", zap.Error(err))
					return nil, status.Error(codes.Internal, err.Error())
				} else {
					log.Error("Failed to complete Hard delete Namespace process", zap.Error(err))
					return nil, status.Error(codes.Internal, err.Error())
				}
			}
			//Added logging
			log.Info("Hard Delete Process Successful")
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
		//Added logging
		log.Error("The Account is not allowed to delete the Namespace")
		return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to delete the Namespace")
	}

	//Added logging
	log.Info("Namespace successfully deleted")
	return &nodepb.DeleteNamespaceResponse{}, nil
}

//UpdateNamespace is a method to delete access to a Namespace for a account
func (n *NamespaceController) UpdateNamespace(ctx context.Context, request *nodepb.UpdateNamespaceRequest) (response *nodepb.UpdateNamespaceResponse, err error) {

	log := n.Log.Named("Update Namespace Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Namespace", request.Namespace.Id),
		zap.Any("FieldMask Paths", request.GetNamespaceMask()),
	)

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//Initialize the Account Controller with Namespace controller data
	a.Repo = n.Repo
	a.Log = n.Log

	//Check if the account is root
	isroot, err := a.IsRoot(ctx, &nodepb.IsRootRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	//Check if the account is admin
	isadmin, err := a.IsAdmin(ctx, &nodepb.IsAdminRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	//Check if the Account has access to Namespace
	resp, err := a.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Account:     requestorID,
		Namespaceid: request.Namespace.Id,
		Action:      nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Initiate update if the account has access
	if resp.GetDecision().GetValue() && (isroot.GetIsRoot() || isadmin.GetIsAdmin()) {
		err = n.Repo.UpdateNamespace(ctx, request)
		if err != nil {
			//Added logging
			log.Error("Failed to update Namespace", zap.Error(err))
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		//Added logging
		log.Error("The Account is not allowed to update the Namespace")
		return nil, status.Error(codes.PermissionDenied, "The Account is not allowed to update the Namespace")
	}

	//Added logging
	log.Info("Namespace successfully updated")
	return &nodepb.UpdateNamespaceResponse{}, nil
}
