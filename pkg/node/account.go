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

//IsRoot is a method that returns if the account has root priviledges or not
func (s *AccountController) IsRoot(ctx context.Context, request *nodepb.IsRootRequest) (response *nodepb.IsRootResponse, err error) {
	account, err := s.Repo.GetAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find account")
	}

	return &nodepb.IsRootResponse{IsRoot: account.IsRoot}, nil
}

//CreateUserAccount is a method for creating user account
func (s *AccountController) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {

	log := s.Log.Named("CreateUserAccount Controller")

	//Added logging
	log.Info("Create Account Controller Method", zap.Any("Function Invoked", nil), zap.String("Account ID:", request.Account.Uid))

	uid, err := s.Repo.CreateUserAccount(ctx, request.Account.Name, request.Account.Password, request.Account.IsRoot, request.Account.Enabled)
	if err != nil {
		//Added logging
		log.Error("Failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Create Account Controller Method", zap.Bool("User Created", true), zap.String("username", request.Account.Name), zap.String("uid", uid))

	if request.CreateGfUser {
		err = s.Grafana.CreateUser(request.Account.Name)
		if err != nil {
			//Added logging
			log.Error("Failed to create grafana user", zap.String("name", request.Account.Name), zap.Error(err))
		} else {
			//Added logging
			log.Info("Create Account Controller Method", zap.Any("Graphana User Created", nil), zap.String("username", request.Account.Name), zap.String("password", request.Account.Password), zap.String("uid", uid))
		}
	}

	return &nodepb.CreateUserAccountResponse{Uid: uid}, nil
}

//AuthorizeNamespace is a method that reutrns if the user has access to namespace
func (s *AccountController) AuthorizeNamespace(ctx context.Context, request *nodepb.AuthorizeNamespaceRequest) (response *nodepb.AuthorizeNamespaceResponse, err error) {
	err = s.Repo.AuthorizeNamespace(ctx, request.GetAccount(), request.GetNamespace(), request.GetAction())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to authorize")
	}

	return &nodepb.AuthorizeNamespaceResponse{}, nil
}

//Authorize is a method that reutrns if the user has access to a particulare node in Dgraph
func (s *AccountController) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {
	err = s.Repo.Authorize(ctx, request.GetAccount(), request.GetNode(), request.GetAction(), request.GetInherit())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to authorize")
	}

	return &nodepb.AuthorizeResponse{}, nil
}

//IsAuthorizedNamespace is a method that reutrns if the user has access to namespace
func (s *AccountController) IsAuthorizedNamespace(ctx context.Context, request *nodepb.IsAuthorizedNamespaceRequest) (response *nodepb.IsAuthorizedNamespaceResponse, err error) {
	root, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: request.GetAccount(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Authorization check failed")
	}

	if root.GetIsRoot() {
		return &nodepb.IsAuthorizedNamespaceResponse{
			Decision: &wrappers.BoolValue{Value: true},
		}, nil
	}

	decision, err := s.Repo.IsAuthorizedNamespace(ctx, request.GetNamespace(), request.GetAccount(), request.GetAction())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.IsAuthorizedNamespaceResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

//SetPassword is a method that allows to change password for the account
func (s *AccountController) SetPassword(ctx context.Context, request *nodepb.SetPasswordRequest) (response *nodepb.SetPasswordResponse, err error) {
	err = s.Repo.SetPassword(ctx, request.Username, request.Password)
	if err != nil {
		return &nodepb.SetPasswordResponse{}, err
	}

	return &nodepb.SetPasswordResponse{}, nil
}

//IsAuthorized is a method that reutrns if the user has access to a node
func (s *AccountController) IsAuthorized(ctx context.Context, request *nodepb.IsAuthorizedRequest) (response *nodepb.IsAuthorizedResponse, err error) {
	log := s.Log.Named("Authorize").With(
		zap.String("request.account", request.GetAccount()),
		zap.String("request.action", request.GetAction().String()),
		zap.String("request.node", request.GetNode()),
	)

	root, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: request.GetAccount(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Authorization check failed")
	}

	if root.GetIsRoot() {
		return &nodepb.IsAuthorizedResponse{
			Decision: &wrappers.BoolValue{Value: true},
		}, nil
	}

	decision, err := s.Repo.IsAuthorized(ctx, request.GetNode(), request.GetAccount(), request.GetAction().String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info("Return decision", zap.Bool("decision", decision))
	return &nodepb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

//GetAccount is a method that reutrns details of the an account
func (s *AccountController) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {
	account, err := s.Repo.GetAccount(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return account, nil
}

//Authenticate is a method that validates user credentials
func (s *AccountController) Authenticate(ctx context.Context, request *nodepb.AuthenticateRequest) (response *nodepb.AuthenticateResponse, err error) {
	ok, uid, defaultNs, err := s.Repo.Authenticate(ctx, request.GetUsername(), request.GetPassword())
	if !ok || (err != nil) {

		return &nodepb.AuthenticateResponse{}, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	return &nodepb.AuthenticateResponse{Success: ok, Account: &nodepb.Account{Uid: uid}, DefaultNamespace: defaultNs}, nil
}

//ListAccounts is a method that list details of the all account
func (s *AccountController) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {
	accounts, err := s.Repo.ListAccounts(ctx)
	if err != nil {
		return &nodepb.ListAccountsResponse{}, status.Error(codes.Internal, "Failed to list accounts")
	}

	return &nodepb.ListAccountsResponse{
		Accounts: accounts,
	}, nil
}

//UpdateAccount is a method that update details of the an account
func (s *AccountController) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.Account, err error) {
	log := s.Log.Named("UpdateUserAccount")
	err = s.Repo.UpdateAccount(ctx, request)

	if err != nil {
		return &nodepb.Account{}, err
	}

	log.Info("Account Updated", zap.String("Account updated:", request.Account.Name))
	return request.Account, nil
}

//DeleteAccount is a method that deletes an account
func (s *AccountController) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (response *nodepb.DeleteAccountResponse, err error) {
	err = s.Repo.DeleteAccount(ctx, request)
	if err != nil {
		return nil, err
	}
	return &nodepb.DeleteAccountResponse{}, nil
}
