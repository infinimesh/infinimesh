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
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	jwt "github.com/golang-jwt/jwt"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type accountAPI struct {
	signingSecret []byte
	client        nodepb.AccountServiceClient
}

//SelfAccount API Method to get details of own Account.
func (a *accountAPI) SelfAccount(ctx context.Context, request *empty.Empty) (response *nodepb.Account, err error) {

	//Added logging
	log.Debug("Self Account API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Self Account API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	//Added logging
	log.Debug("Self Account API Method: Own Account Details Obtained")
	return a.client.GetAccount(ctx, &nodepb.GetAccountRequest{
		Id: account,
	})
}

//API Method to Get Details of an Account
func (a *accountAPI) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {
	account, ok := ctx.Value("account_id").(string)

	//Added logging
	log.Debug("Get Account API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	if !ok {
		//Added logging
		log.Error("Get Account API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	//Validate if the account is root or not
	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {
		return a.client.GetAccount(ctx, request)
	}

	//Validate if the account is admin or not
	if res, err := a.client.IsAdmin(ctx, &nodepb.IsAdminRequest{
		Account: account,
	}); err == nil && res.GetIsAdmin() {
		res, err := a.client.GetAccount(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("Get Account API Method: The Account does not have permission to get details")
	return &nodepb.Account{}, status.Error(codes.PermissionDenied, "The Account does not have permission to get details")
}

//Token is method to get token for an Account
func (a *accountAPI) Token(ctx context.Context, request *apipb.TokenRequest) (response *apipb.TokenResponse, err error) {

	//Added logging
	log.Debug("Generate Token Method: Function Invoked", zap.String("Requestor ID", request.Username))

	resp, err := a.client.Authenticate(ctx, &nodepb.AuthenticateRequest{Username: request.GetUsername(), Password: request.GetPassword()})
	if err != nil {
		//Added logging
		log.Error("Generate Token Method: Authentication for User failed", zap.Error(err))
		return nil, err
	}

	if resp.GetSuccess() {
		if resp.Account == nil {
			//Added logging
			log.Error("Generate Token Method: Failed to check credentials")
			return nil, status.Error(codes.Internal, "Failed to check credentials")
		}

		claim := jwt.MapClaims{}
		claim[accountIDClaim] = resp.Account.Uid

		if request.GetExpireTime() != "" {
			exp, err := strconv.Atoi(request.GetExpireTime())
			if err != nil {
				//Added logging
				log.Error("Generate Token Method: Parising for Expiry Time failed", zap.Error(err))
				return nil, status.Error(codes.InvalidArgument, "Can't parse expire time")
			}
			claim[expiresAt] = time.Now().UTC().Add(time.Duration(exp) * time.Second).Unix()
		}

		if ruleset := request.GetRuleset(); len(ruleset) > 0 {
			claim[tokenRestrictedClaim] = true
			prefix := "infinimesh.api." // Should be equal to api.proto package name + .
			for _, rule := range ruleset {
				claim[prefix+rule.GetType()] = rule.GetIds()
			}
		} else {
			claim[tokenRestrictedClaim] = false
		}

		//Added logging
		log.Debug("Generate Token Method: Get Token for the Authenticated User")

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(a.signingSecret)
		if err != nil {
			//Added logging
			log.Error("Generate Token Method: Failed to sign token", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to sign token")
		}

		//Added logging
		log.Debug("Generate Token Method: Token generation successful")
		return &apipb.TokenResponse{Token: tokenString}, nil
	}

	//Added logging
	log.Error("Generate Token Method: The User credentials are not valid", zap.Error(err))
	return nil, status.Error(codes.Unauthenticated, "The User credentials are not valid")
}

//API Method to Update an Account
func (a *accountAPI) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.UpdateAccountResponse, err error) {

	//Added logging
	log.Info("Update Account API Method: Function Invoked",
		zap.String("Account", request.Account.Uid),
		zap.Any("FieldMask", request.FieldMask))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the Update Account controller for server
	res, err := a.client.UpdateAccount(ctx, request)
	if err != nil {
		//Added logging
		log.Error("Update Account API Method: Failed to update Account", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Update Account API Method: Account succesfully updated")
	return res, nil
}

//API Method to Create an Account
func (a *accountAPI) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {

	//Temporary Assigning the request.password to request.account.password
	if request.Account.Password == "" {
		request.Account.Password = request.Password
	}

	//Added logging
	log.Info("Create Account API Method: Function Invoked",
		zap.String("Requestor ID", ctx.Value("account_id").(string)),
		zap.String("Account", request.Account.Name),
		zap.Bool("Enabled", request.Account.Enabled),
		zap.Bool("IsRoot", request.Account.IsRoot),
		zap.Bool("IsAdmin", request.Account.IsAdmin),
	)

	account, ok := ctx.Value("account_id").(string)

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", account)

	if !ok {
		//Added logging
		log.Error("Create Account API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	//Validated that required data is populated with values
	if request.Account.Name == "" {
		//Added logging
		log.Error("Create Account API Method: Data Validation for Account Creation", zap.String("Error", "The Name cannot not be empty"))
		return nil, status.Error(codes.FailedPrecondition, "The Name cannot not be empty")
	}

	if request.Account.Password == "" {
		//Added logging
		log.Error("Create Account API Method: Data Validation for Account Creation", zap.String("Error", "The Password cannot not be empty"))
		return nil, status.Error(codes.FailedPrecondition, "The Password cannot not be empty")
	}

	//Validate if the account is root or not
	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {
		res, err := a.client.CreateUserAccount(ctx, request)
		return res, err
	}

	//Validate if the account is admin or not
	if res, err := a.client.IsAdmin(ctx, &nodepb.IsAdminRequest{
		Account: account,
	}); err == nil && res.GetIsAdmin() {
		res, err := a.client.CreateUserAccount(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("Create Account API Method: The Account does not have permission to create another account")
	return &nodepb.CreateUserAccountResponse{}, status.Error(codes.PermissionDenied, "The Account does not have permission to create another account")
}

//API Method to List all account.
func (a *accountAPI) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {

	//Added logging
	log.Debug("List Accounts API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("List Accounts API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", account)

	//Invoke the List Account controller for server
	res, err := a.client.ListAccounts(ctx, request)
	return res, err

}

//API Method to Delete an Account
func (a *accountAPI) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (response *nodepb.DeleteAccountResponse, err error) {

	//Added logging
	log.Info("Delete Account API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Delete Account API Method: The Account is not authenticated")
		return nil, status.Error(codes.Unauthenticated, "The Account is not authenticated")
	}

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", account)

	//Invoke the controller for server
	res, err := a.client.DeleteAccount(ctx, request)
	return res, err

}

//API Method to Assign Owner to an Account
func (a *accountAPI) AssignOwner(ctx context.Context, request *nodepb.OwnershipRequest) (response *nodepb.OwnershipResponse, err error) {

	//Added logging
	log.Info("Assign Owner API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the controller for server
	res, err := a.client.AssignOwner(ctx, request)
	return res, err

}

//API Method to Remove Owner from an Account
func (a *accountAPI) RemoveOwner(ctx context.Context, request *nodepb.OwnershipRequest) (response *nodepb.OwnershipResponse, err error) {

	//Added logging
	log.Info("Remove Owner API Method: Function Invoked", zap.String("Requestor ID", ctx.Value("account_id").(string)))

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", ctx.Value("account_id").(string))

	//Invoke the controller for server
	res, err := a.client.RemoveOwner(ctx, request)
	return res, err

}
