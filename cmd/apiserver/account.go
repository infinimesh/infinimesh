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
	account, ok := ctx.Value("account_id").(string)

	//Added logging
	log.Info("Self Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	return a.client.GetAccount(ctx, &nodepb.GetAccountRequest{
		Id: account,
	})
}

//API Method to Get Details of an Account
func (a *accountAPI) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {
	account, ok := ctx.Value("account_id").(string)

	//Added logging
	log.Info("Get Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {
		return a.client.GetAccount(ctx, request)
	}

	return &nodepb.Account{}, status.Error(codes.PermissionDenied, "The account does not have permission to get details.")
}

//Method to get token for an Account
func (a *accountAPI) Token(ctx context.Context, request *apipb.TokenRequest) (response *apipb.TokenResponse, err error) {

	//Added logging
	log.Info("Token Generate Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	resp, err := a.client.Authenticate(ctx, &nodepb.AuthenticateRequest{Username: request.GetUsername(), Password: request.GetPassword()})
	if err != nil {
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
		log.Info("Token Generate Method", zap.String("Get Token for the Authenticated User", request.Username))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(a.signingSecret)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to sign token")
		}

		return &apipb.TokenResponse{Token: tokenString}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
}

//API Method to Update an Account
func (a *accountAPI) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.Account, err error) {
	account, ok := ctx.Value("account_id").(string)

	//Added logging
	log.Info("Update Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		res, err := a.client.UpdateAccount(ctx, request)
		return res, err
	}

	return &nodepb.Account{}, status.Error(codes.PermissionDenied, "The account does not have permission to update details.")
}

//API Method to Create an Account
func (a *accountAPI) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {

	//Added logging
	log.Info("Create Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	account, ok := ctx.Value("account_id").(string)

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	//Validated that required data is populated with values
	if request.Account.Name == "" {
		//Added logging
		log.Error("Create Account API Method", zap.Bool("Data Validation for Account Creation", false), zap.String("Error", "The Name cannot not be empty."))
		return nil, status.Error(codes.FailedPrecondition, "The Name cannot not be empty.")
	}

	if request.Account.Password == "" {
		//Added logging
		log.Error("Create Account API Method", zap.Bool("Data Validation for Account Creation", false), zap.String("Error", "The Password cannot not be empty."))
		return nil, status.Error(codes.FailedPrecondition, "The Password cannot not be empty.")
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

	return &nodepb.CreateUserAccountResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to create another account.")
}

//API Method to List all account
func (a *accountAPI) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {

	//Added logging
	log.Info("List Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		//Added logging
		log.Info("List Account API Method", zap.Bool("Validation for Root Account", res.GetIsRoot()))

		res, err := a.client.ListAccounts(ctx, request)
		return res, err
	}

	return &nodepb.ListAccountsResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to list details.")

}

//API Method to Delete an Account
func (a *accountAPI) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (response *nodepb.DeleteAccountResponse, err error) {

	//Added logging
	log.Info("Delete Account API Method", zap.Bool("Function Invoked", true), zap.Any("Account ID:", ctx.Value("account_id")))

	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "The account is not authenticated.")
	}

	if res, err := a.client.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: account,
	}); err == nil && res.GetIsRoot() {

		//Added logging
		log.Info("Delete Account API Method", zap.Bool("Validation for Root Account", res.GetIsRoot()))

		res, err := a.client.DeleteAccount(ctx, request)
		return res, err
	}
	return &nodepb.DeleteAccountResponse{}, status.Error(codes.PermissionDenied, "The account does not have permission to update details.")
}
