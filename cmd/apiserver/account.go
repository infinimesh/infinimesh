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

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type accountAPI struct {
	signingSecret []byte
	client        nodepb.AccountServiceClient
}

//API Method to get details of own Account
func (a *accountAPI) SelfAccount(ctx context.Context, request *empty.Empty) (response *nodepb.Account, err error) {

	//Added logging
	log.Info("Self Account API Method", zap.Bool("Function Invoked", true), zap.String("Account ID:", ctx.Value("account_id").(string)))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Self Account API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	//Added logging
	log.Info("Self Account API Method", zap.Bool("Own Account Details Obtained", true))
	return a.client.GetAccount(ctx, &nodepb.GetAccountRequest{
		Id: account,
	})
}

//API Method to Get Details of an Account
func (a *accountAPI) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {
	account, ok := ctx.Value("account_id").(string)

	//Added logging
	log.Info("Get Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account", request.Id))

	if !ok {
		//Added logging
		log.Error("Get Account API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {
		return a.client.GetAccount(ctx, request)
	}

	//Added logging
	log.Error("Get Account API Method", zap.Bool("The account does not have permission to get details", true))
	return &nodepb.Account{}, status.Error(codes.PermissionDenied, "The account does not have permission to get details")
}

//Method to get token for an Account
func (a *accountAPI) Token(ctx context.Context, request *apipb.TokenRequest) (response *apipb.TokenResponse, err error) {

	//Added logging
	log.Info("Generate Token Method", zap.Bool("Function Invoked", true), zap.String("Account ID:", request.Username))

	resp, err := a.client.Authenticate(ctx, &nodepb.AuthenticateRequest{Username: request.GetUsername(), Password: request.GetPassword()})
	if err != nil {
		//Added logging
		log.Error("Generate Token Method", zap.Bool("Authentication for User failed", true), zap.Error(err))
		return nil, err
	}

	if resp.GetSuccess() {
		if resp.Account == nil {
			return nil, status.Error(codes.Internal, "Failed to check credentials")
		}

		claim := jwt.MapClaims{}
		claim[accountIDClaim] = resp.Account.Uid

		if request.GetExpireTime() != "" {
			exp, err := strconv.Atoi(request.GetExpireTime())
			if err != nil {
				//Added logging
				log.Error("Generate Token Method", zap.Bool("Parising for Expiry Time failed", true), zap.Error(err))
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
		log.Info("Generate Token Method", zap.Bool("Get Token for the Authenticated User", true))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(a.signingSecret)
		if err != nil {
			//Added logging
			log.Error("Generate Token Method", zap.Bool("Token Signing failed", true), zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to sign token")
		}

		//Added logging
		log.Info("Generate Token Method", zap.Bool("Token generation successful", true))
		return &apipb.TokenResponse{Token: tokenString}, nil
	}

	//Added logging
	log.Error("Generate Token Method", zap.Bool("The User credentials are not valid", true), zap.Error(err))
	return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
}

//API Method to Update an Account
func (a *accountAPI) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.Account, err error) {

	//Added logging
	log.Info("Update Account API Method", zap.Bool("Function Invoked", true),
		zap.Any("Account", request.Account.Uid),
		zap.Any("Name", request.Account.Name))

	account, ok := ctx.Value("account_id").(string)

	if !ok {
		//Added logging
		log.Error("Update Account API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		res, err := a.client.UpdateAccount(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("Update Account API Method", zap.Bool("The account does not have permission to update details", true))
	return &nodepb.Account{}, status.Error(codes.PermissionDenied, "The account does not have permission to update details")
}

//API Method to Create an Account
func (a *accountAPI) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {

	//Temporary Assigning the request.password to request.account.password
	if request.Account.Password == "" {
		request.Account.Password = request.Password
	}

	//Added logging
	log.Info("Create Account API Method", zap.Bool("Function Invoked", true),
		zap.String("Account", request.Account.Name),
		zap.Bool("Enabled", request.Account.Enabled),
		zap.Bool("IsRoot", request.Account.IsRoot))

	account, ok := ctx.Value("account_id").(string)

	if !ok {
		//Added logging
		log.Error("Create Account API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	//Validated that required data is populated with values
	if request.Account.Name == "" {
		//Added logging
		log.Error("Create Account API Method", zap.Bool("Data Validation for Account Creation", true), zap.String("Error", "The Name cannot not be empty"))
		return nil, status.Error(codes.FailedPrecondition, "The Name cannot not be empty")
	}

	if request.Account.Password == "" {
		//Added logging
		log.Error("Create Account API Method", zap.Bool("Data Validation for Account Creation", true), zap.String("Error", "The Password cannot not be empty"))
		return nil, status.Error(codes.FailedPrecondition, "The Password cannot not be empty")
	}

	//Validate if the account is root or not
	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		//Added logging
		log.Info("Create Account API Method", zap.Bool("Validation for Root Account", res.GetIsRoot()))

		res, err := a.client.CreateUserAccount(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("Create Account API Method", zap.Bool("The account does not have permission to create another account", true))
	return &nodepb.CreateUserAccountResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to create another account")
}

//API Method to List all account
func (a *accountAPI) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {

	//Added logging
	log.Info("List Accounts API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("List Accounts API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		//Added logging
		log.Info("List Accounts API Method", zap.Bool("Validation for Root Account", res.GetIsRoot()))

		res, err := a.client.ListAccounts(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("List Accounts API Method", zap.Bool("The account does not have permission to list details", true))
	return &nodepb.ListAccountsResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to list details")

}

//API Method to Delete an Account
func (a *accountAPI) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (response *nodepb.DeleteAccountResponse, err error) {

	//Added logging
	log.Info("Delete Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		//Added logging
		log.Error("Delete Account API Method", zap.Bool("The account is not authenticated", true), zap.Bool("Authentication", ok))
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated")
	}

	//Added the requestor account id to context metadata so that it can be passed on to the server
	ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", account)

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		//Added logging
		log.Info("Delete Account API Method", zap.Bool("Validation for Root Account", res.GetIsRoot()))

		res, err := a.client.DeleteAccount(ctx, request)
		return res, err
	}

	//Added logging
	log.Error("Delete Account API Method", zap.Bool("The account does not have permission to delete another account", true))
	return &nodepb.DeleteAccountResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to delete another account")
}
