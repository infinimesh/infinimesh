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

	"github.com/dgraph-io/dgo"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/grafana"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

//AccountController is a Data type for Account Controller file
type AccountController struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Grafana *grafana.Client
	Repo    Repo
}

//Validation method does the pre-checks for a REST request
func Validation(ctx context.Context, log *zap.Logger) (md metadata.MD, acc string, err error) {

	//Get the metadata from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//Added logging
		log.Error("Failed to get metadata from context", zap.Error(err))
		return nil, "", status.Error(codes.Aborted, "Failed to get metadata from context")
	}

	//Check for Authentication
	requestorID := md.Get("requestorID")
	if requestorID == nil {
		//Added logging
		log.Error("The account is not authenticated")
		return nil, "", status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	return md, requestorID[0], nil
}

//IsOwnedbyAdmin is a method to validation if the account is owned by admin or not
func (s *AccountController) IsOwnedbyAdmin(ctx context.Context, log *zap.Logger, requestorID string, accountID string) (isowned bool, err error) {

	//Get List of accounts owned by admin
	adminOwnedAcc, _ := s.ListAccounts(ctx, &nodepb.ListAccountsRequest{})

	//For Admin user the account to be updated should be owned by the Admin
	for i := 0; i < len(adminOwnedAcc.Accounts); i++ {
		if adminOwnedAcc.Accounts[i].Uid == accountID {
			break
		}

		if (i + 1) == len(adminOwnedAcc.Accounts) {
			//Added logging
			log.Error("The Account is not owned by the Admin")
			return false, status.Error(codes.PermissionDenied, "The Account is not owned by the Admin")
		}
	}

	//Added logging
	log.Info("Validation for Admin Owned Account", zap.Bool("Validation Result", true))
	return true, nil
}

//IsRoot is a method that returns if the account has root priviledges or not
func (s *AccountController) IsRoot(ctx context.Context, request *nodepb.IsRootRequest) (response *nodepb.IsRootResponse, err error) {

	log := s.Log.Named("IsRoot Validation Controller")

	account, err := s.Repo.GetAccount(ctx, request.GetAccount())
	if err != nil {
		//Added logging
		log.Error("Could not find account", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Could not find account")
	}

	//Added logging
	log.Info("Validation for Root Account", zap.String("Account", request.Account), zap.Bool("Validation Result", account.IsRoot))
	return &nodepb.IsRootResponse{IsRoot: account.IsRoot}, nil
}

//IsAdmin is a method that returns if the account has root priviledges or not
func (s *AccountController) IsAdmin(ctx context.Context, request *nodepb.IsAdminRequest) (response *nodepb.IsAdminResponse, err error) {

	log := s.Log.Named("IsAdmin Validation Controller")

	account, err := s.Repo.GetAccount(ctx, request.GetAccount())
	if err != nil {
		//Added logging
		log.Error("Could not find account", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Could not find account")
	}

	//Added logging
	log.Info("Validation for Admin Account", zap.String("Account", request.Account), zap.Bool("Validation Result", account.IsAdmin))
	return &nodepb.IsAdminResponse{IsAdmin: account.IsAdmin}, nil
}

//CreateUserAccount is a method for creating user account
func (s *AccountController) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {

	log := s.Log.Named("CreateUserAccount Controller")
	//Added logging
	log.Info("Function Invoked", zap.Any("Account", request.Account.Name))

	//Get metadata from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	uid, err := s.Repo.CreateUserAccount(ctx, request.Account.Name, request.Account.Password, request.Account.IsRoot, request.Account.IsAdmin, request.Account.Enabled)
	if err != nil {
		//Added logging
		log.Error("Failed to create User", zap.String("Name", request.Account.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if request.CreateGfUser {
		err = s.Grafana.CreateUser(request.Account.Name)
		if err != nil {
			//Added logging
			log.Error("Failed to create Grafana user", zap.String("Name", request.Account.Name), zap.Error(err))
		} else {
			//Added logging
			log.Info("Graphana User Created", zap.String("Grafana UserName", request.Account.Name), zap.String("uid", uid))
		}
	}

	if len(request.Account.Owner) > 0 {
		err := s.Repo.AssignOwner(ctx, requestorID, uid)
		if err != nil {
			//Added logging
			log.Error("Failed to assign owner to the account", zap.String("Account", uid), zap.Error(err))
			return nil, err
		}
	}

	//Added logging
	log.Info("Infinimesh User Created", zap.String("UserName", request.Account.Name), zap.String("Account ID", uid), zap.String("Owned by", requestorID))
	return &nodepb.CreateUserAccountResponse{Uid: uid}, nil
}

//AuthorizeNamespace is a method that provides the user access to namespace
func (s *AccountController) AuthorizeNamespace(ctx context.Context, request *nodepb.AuthorizeNamespaceRequest) (response *nodepb.AuthorizeNamespaceResponse, err error) {

	log := s.Log.Named("Authorize Namespace Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Account", request.Account),
		zap.String("Namespace", request.Namespace),
		zap.String("Action", request.Action.String()))

	err = s.Repo.AuthorizeNamespace(ctx, request.GetAccount(), request.GetNamespace(), request.GetAction())
	if err != nil {
		//Added logging
		log.Error("Failed to provide Authorization to Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to provide Authorization to Namespace")
	}

	//Added logging
	log.Info("Account Authorized to Access Namespace")
	return &nodepb.AuthorizeNamespaceResponse{}, nil
}

//Authorize is a method that provides the user access to a particulare node in Dgraph
func (s *AccountController) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {

	log := s.Log.Named("Authorize Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Account", request.Account),
		zap.String("Node", request.Node),
		zap.String("Action", request.Action))

	err = s.Repo.Authorize(ctx, request.GetAccount(), request.GetNode(), request.GetAction(), request.GetInherit())
	if err != nil {
		//Added logging
		log.Error("Failed to provide Authorization to Node", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to provide Authorization to Node")
	}

	//Added logging
	log.Info("Account Authorized to Access Node")
	return &nodepb.AuthorizeResponse{}, nil
}

//IsAuthorizedNamespace is a method that returns true if the user has access to namespace
func (s *AccountController) IsAuthorizedNamespace(ctx context.Context, request *nodepb.IsAuthorizedNamespaceRequest) (response *nodepb.IsAuthorizedNamespaceResponse, err error) {

	log := s.Log.Named("Is Authorize Namespace Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Account", request.Account),
		zap.String("Namespace", request.Namespaceid),
		zap.String("Action", request.Action.String()))

	//Check if the account is root
	isroot, err := s.IsRoot(ctx, &nodepb.IsRootRequest{Account: request.Account})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the account")
	}

	//Provide access if the account is root
	if isroot.GetIsRoot() {
		log.Info("Authorization check successful for the Account and the Namespace as root")
		return &nodepb.IsAuthorizedNamespaceResponse{
			Decision: &wrappers.BoolValue{Value: true},
		}, nil
	}

	decision, err := s.Repo.IsAuthorizedNamespace(ctx, request.GetNamespaceid(), request.GetAccount(), request.GetAction())
	if err != nil {
		//Added logging
		log.Error("Authorization check failed for the Account and the Namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Authorization check successful for the Account and the Namespace", zap.Bool("Decision for Access is", decision))
	return &nodepb.IsAuthorizedNamespaceResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

//SetPassword is a method that allows to change password for the account
func (s *AccountController) SetPassword(ctx context.Context, request *nodepb.SetPasswordRequest) (response *nodepb.SetPasswordResponse, err error) {

	log := s.Log.Named("Set Password Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Username))

	err = s.Repo.SetPassword(ctx, request.Username, request.Password)
	if err != nil {
		//Added logging
		log.Error("Password change failed", zap.Error(err))
		return &nodepb.SetPasswordResponse{}, err
	}

	//Added logging
	log.Info("Password changed sucesssful")
	return &nodepb.SetPasswordResponse{}, nil
}

//IsAuthorized is a method that reutrns if the user has access to a node
func (s *AccountController) IsAuthorized(ctx context.Context, request *nodepb.IsAuthorizedRequest) (response *nodepb.IsAuthorizedResponse, err error) {

	log := s.Log.Named("Is Authorized Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Account", request.Account),
		zap.String("Node", request.Node),
		zap.String("Action", request.Action.String()))

	root, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: request.GetAccount(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Authorization check failed")
	}

	if root.GetIsRoot() {
		//Added logging
		log.Info("Authorization check successful for the Account and the Node as root account")
		return &nodepb.IsAuthorizedResponse{
			Decision: &wrappers.BoolValue{Value: true},
		}, nil
	}

	decision, err := s.Repo.IsAuthorized(ctx, request.GetNode(), request.GetAccount(), request.GetAction().String())
	if err != nil {
		//Added logging
		log.Error("Authorization check failed for the Account and the Node", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Authorization check successful for the Account and the Node")
	return &nodepb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

//GetAccount is a method that reutrns details of the an account
func (s *AccountController) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {

	log := s.Log.Named("Get Account Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Id))

	account, err := s.Repo.GetAccount(ctx, request.Id)
	if err != nil {
		//Added logging
		log.Error("Unable able to get Account", zap.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	//Added logging
	log.Info("Account details obtained")
	return account, nil
}

//Authenticate is a method that validates user credentials
func (s *AccountController) Authenticate(ctx context.Context, request *nodepb.AuthenticateRequest) (response *nodepb.AuthenticateResponse, err error) {

	log := s.Log.Named("Authenticate Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Username))

	ok, uid, defaultNs, err := s.Repo.Authenticate(ctx, request.GetUsername(), request.GetPassword())
	if !ok || (err != nil) {
		//Added logging
		log.Error("Authentication Failed", zap.Error(err))
		return &nodepb.AuthenticateResponse{}, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	//Added logging
	log.Info("Authentication successsful")
	return &nodepb.AuthenticateResponse{Success: ok, Account: &nodepb.Account{Uid: uid}, DefaultNamespace: defaultNs}, nil
}

//ListAccounts is a method that list details of the all account
func (s *AccountController) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {

	log := s.Log.Named("List Accounts Controller")
	//Added logging
	log.Info("Function Invoked")

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	var accounts []*nodepb.Account

	//Check the account priviledges
	if res, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: requestorID,
	}); err == nil && res.GetIsRoot() {
		//Get the list if the account has root permissions
		log.Info("List Accounts with root privilidges")
		accounts, err = s.Repo.ListAccounts(ctx)
		if err != nil {
			//Added logging
			log.Error("Failed to list Accounts as root", zap.Error(err))
			return &nodepb.ListAccountsResponse{}, err
		}
	} else if res, err := s.IsAdmin(ctx, &nodepb.IsAdminRequest{
		Account: requestorID,
	}); err == nil && res.GetIsAdmin() {
		//Get the list if the account has admin permissions
		log.Info("List Accounts with admin privilidges")
		accounts, err = s.Repo.ListAccountsforAdmin(ctx, requestorID)
		if err != nil {
			//Added logging
			log.Error("Failed to list Accounts as admin", zap.Error(err))
			return &nodepb.ListAccountsResponse{}, status.Error(codes.OK, err.Error())
		}
	} else {
		//Added logging
		log.Error("The Account does not have permission to list details")
		return &nodepb.ListAccountsResponse{}, status.Error(codes.PermissionDenied, "The Account does not have permission to list details")
	}

	//Added logging
	log.Info("List Account successful")
	return &nodepb.ListAccountsResponse{
		Accounts: accounts,
	}, nil
}

//UpdateAccount is a method that update details of the an account
func (s *AccountController) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.UpdateAccountResponse, err error) {

	log := s.Log.Named("Update Account Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Account.Uid))

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	//Check if the account is root
	isroot, err := s.IsRoot(ctx, &nodepb.IsRootRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the Account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the Account")
	}

	//Check if the account is admin
	isadmin, err := s.IsAdmin(ctx, &nodepb.IsAdminRequest{Account: requestorID})
	if err != nil {
		//Added logging
		log.Error("Unable to get permissions for the Account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Unable to get permissions for the Account")
	}

	//isself is to make sure that you can only update name and password for your own account
	var isself bool
	//Check if the account is self account i.e. user trying to update his account
	if requestorID == request.GetAccount().Uid {
		isself = true
	}

	//Check if the account is owned by the admin
	_, err = s.IsOwnedbyAdmin(ctx, log, requestorID, request.Account.Uid)
	if err != nil {
		return nil, err
	}

	//Added logging
	log.Info("Validation for Self Account", zap.Bool("Validation Result", isself))

	//Perform update account if the requestor as access
	if isroot.IsRoot || isadmin.IsAdmin || isself {
		err = s.Repo.UpdateAccount(ctx, request, isself)
		if err != nil {
			//Added logging
			log.Error("Failed to update Account", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to update Account : "+err.Error())
		}
	} else {
		//Added logging
		log.Error("The Account does not have permission to update details")
		return nil, status.Error(codes.PermissionDenied, "The Account does not have permission to update details")
	}

	//Added Logging
	log.Info("Update Account successful")
	return &nodepb.UpdateAccountResponse{}, nil
}

//DeleteAccount is a method that deletes an account
func (s *AccountController) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (response *nodepb.DeleteAccountResponse, err error) {

	log := s.Log.Named("Delete Account Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Uid))

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	if res, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: requestorID,
	}); err == nil && !res.GetIsRoot() {
		if res, err := s.IsAdmin(ctx, &nodepb.IsAdminRequest{
			Account: requestorID,
		}); err == nil && !res.GetIsAdmin() {
			//Validate if the account is admin or not
			//Added logging
			log.Error("The Account does not have permission to delete another account")
			return &nodepb.DeleteAccountResponse{}, status.Error(codes.PermissionDenied, "The Account does not have permission to delete another account")
		}
	}

	//Get account details for validation
	account, err := s.Repo.GetAccount(ctx, request.Uid)
	if err != nil {
		//Added logging
		log.Error("Failed to get account details", zap.Error(err))
		return nil, status.Error(codes.Aborted, "Failed to get account details"+err.Error())
	}

	//Validation to make sure root account cannot be deleted
	if account.IsRoot {
		//Added logging
		log.Error("Cannot delete root Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot delete root Account")
	}

	//Validation to make sure admin account cannot be deleted
	if account.IsAdmin {
		//Added logging
		log.Error("Cannot delete admin Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot delete admin Account")
	}

	//Validation to make sure enabled account cannot be deleted
	if account.Enabled {
		//Added logging
		log.Error("Cannot delete enabled Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot delete enabled Account")
	}

	//Check if the account is owned by the admin
	_, err = s.IsOwnedbyAdmin(ctx, log, requestorID, request.Uid)
	if err != nil {
		return nil, err
	}

	//Call the delete query when all the validation pass
	err = s.Repo.DeleteAccount(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Failed to delete Account", zap.Error(err))
		return nil, err
	}

	//Added Logging
	log.Info("Delete Account successful")
	return &nodepb.DeleteAccountResponse{}, nil
}

//AssignOwner is a method that assigns an Admin to an account
func (s *AccountController) AssignOwner(ctx context.Context, request *nodepb.OwnershipRequest) (response *nodepb.OwnershipResponse, err error) {

	log := s.Log.Named("Assign Owner Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Owner", request.Ownerid), zap.String("Account", request.Accountid))

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	if res, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: requestorID,
	}); err == nil && !res.GetIsRoot() {
		if res, err := s.IsAdmin(ctx, &nodepb.IsAdminRequest{
			Account: requestorID,
		}); err == nil && !res.GetIsAdmin() {
			//Validate if the account is admin or not
			//Added logging
			log.Error("The Account does not have permission to assign owner to the Account")
			return &nodepb.OwnershipResponse{}, status.Error(codes.PermissionDenied, "The Account does not have permission to assign owner to the Account")
		}
	}

	//Get account details for validation
	account, err := s.Repo.GetAccount(ctx, request.Accountid)
	if err != nil {
		//Added logging
		log.Error("Failed to get account details", zap.Error(err))
		return nil, status.Error(codes.Aborted, "Failed to get account details"+err.Error())
	}

	//Get account details for validation
	owner, err := s.Repo.GetAccount(ctx, request.Ownerid)
	if err != nil {
		//Added logging
		log.Error("Failed to get account details", zap.Error(err))
		return nil, status.Error(codes.Aborted, "Failed to get account details"+err.Error())
	}

	//Validation to make sure root account cannot be deleted
	if account.IsRoot {
		//Added logging
		log.Error("Cannot assign to root Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot assign to root Account")
	}

	//Validation to make sure the new owner is admin
	if !owner.IsAdmin {
		//Added logging
		log.Error("Cannot assign owner as the Account is not admin Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot assign owner as the Account is not admin Account")
	}

	//Check if the account is owned by the admin
	_, err = s.IsOwnedbyAdmin(ctx, log, requestorID, request.Accountid)
	if err != nil {
		return nil, err
	}

	//Call the Assign Owner query when all the validation pass
	err = s.Repo.AssignOwner(ctx, request.Ownerid, request.Accountid)
	if err != nil {
		//Added logging
		log.Error("Failed to assign owner to the Account", zap.Error(err))
		return nil, err
	}

	//Added Logging
	log.Info("Assign Owner successful")
	return &nodepb.OwnershipResponse{}, nil
}

//RemoveOwner is a method that assigns an Admin to an account
func (s *AccountController) RemoveOwner(ctx context.Context, request *nodepb.OwnershipRequest) (response *nodepb.OwnershipResponse, err error) {

	log := s.Log.Named("Remove Owner Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Owner", request.Ownerid), zap.String("Account", request.Accountid))

	//Get metadata and from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	if res, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: requestorID,
	}); err == nil && !res.IsRoot {
		if res, err := s.IsAdmin(ctx, &nodepb.IsAdminRequest{
			Account: requestorID,
		}); err == nil && !res.IsAdmin {
			//Validate if the account is admin or not
			//Added logging
			log.Error("The Account does not have permission to assign owner to the Account")
			return &nodepb.OwnershipResponse{}, status.Error(codes.PermissionDenied, "The Account does not have permission to assign owner to the Account")
		}
	}

	//Get account details for validation
	_, err = s.Repo.GetAccount(ctx, request.Ownerid)
	if err != nil {
		//Added logging
		log.Error("Failed to get account details", zap.Error(err))
		return nil, err
	}

	//Get account details for validation
	account, err := s.Repo.GetAccount(ctx, request.Accountid)
	if err != nil {
		//Added logging
		log.Error("Failed to get account details", zap.Error(err))
		return nil, err
	}

	//Validation to make sure root account cannot be deleted
	if account.IsRoot {
		//Added logging
		log.Error("Cannot remove owner from root Account")
		return nil, status.Error(codes.FailedPrecondition, "Cannot remove owner from root Account")
	}

	//Check if the account is owned by the admin
	_, err = s.IsOwnedbyAdmin(ctx, log, requestorID, request.Accountid)
	if err != nil {
		return nil, err
	}

	//Call the Remove Owner query when all the validation pass
	err = s.Repo.RemoveOwner(ctx, request.Ownerid, request.Accountid)
	if err != nil {
		//Added logging
		log.Error("Failed to remove owner from the Account", zap.Error(err))
		return nil, err
	}

	//Added Logging
	log.Info("Remove Owner successful")
	return &nodepb.OwnershipResponse{}, nil
}
