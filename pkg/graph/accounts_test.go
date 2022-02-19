/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graph

import (
	"context"
	"errors"
	"testing"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/golang-jwt/jwt/v4"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/internal"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"github.com/infinimesh/infinimesh/pkg/node/proto/accounts"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	log *zap.Logger
	arangodbHost string
	arangodbCred string

	ctrl AccountsController

	rootCtx context.Context
)

func init() {
	viper.AutomaticEnv()
	log = zap.NewExample()

	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	db := schema.InitDB(log, arangodbHost, arangodbCred, "infinimesh")

	ctrl = NewAccountsController(log, db)

	md := metadata.New(map[string]string{"requestorid": "infinimesh"})
	rootCtx = metadata.NewIncomingContext(context.Background(), md)
}

func TestValidate(t *testing.T) {
	t.Log("Testing Validate function")
	_, id, err := Validate(rootCtx, log)
	if err != nil {
		t.Fatalf("Error validating correct context: %v", err)
	}

	if id != "infinimesh" {
		t.Fatalf("Requestor ID supposed to be 'infinimesh', got: %s", id)
	}
}

func TestValidateEmptyContext(t *testing.T) {
	t.Log("Testing Validate function with Empty Context")
	_, _, err := Validate(context.TODO(), log)
	if err == nil {
		t.Fatal("Error is nil, while it's supposed to be 'Failed to get metadata from context'")
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Error reading status from error, original error: %v", err)
	}

	if s.Code() != codes.Aborted && s.Message() != ("Failed to get metadata from context") {
		t.Fatalf("Error supposed to be Aborted: Failed to get metadata from context, but received %s: %s", s.Code().String(), s.Message())
	}
}

func TestValidateEmptyRequestor(t *testing.T) {
	t.Log("Testing Validate function with Empty Context")
	md := metadata.New(map[string]string{"notrequestorid": "infinimesh"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, _, err := Validate(ctx, log)
	if err == nil {
		t.Fatal("Error is nil, while it's supposed to be 'The account is not authenticated'")
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Error reading status from error, original error: %v", err)
	}

	if s.Code() != codes.Unauthenticated && s.Message() != ("The account is not authenticated") {
		t.Fatalf("Error supposed to be Unauthenticated: The account is not authenticated, but received %s: %s", s.Code().String(), s.Message())
	}
}

func TestFalseCredentialsType(t *testing.T) {
	t.Log("Creating Sample Account with unsupported Credentials")
	username := randomdata.SillyName()
	_, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: false,
		},
		Credentials: &accounts.Credentials{
			Type: "unsupported",
		},
	})
	if err == nil && err.Error() != "auth type is wrong" {
		t.Error("Create isn't returning error, despite Credentials type must be unsupported")
	}
}

func TestAuthorizeDisabledAccount(t *testing.T) {
	t.Log("Creating Sample Disabled Account")
	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	credentials := &accounts.Credentials{
		Type: "standard",
		Data: []string{username, password},
	}

	_, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: false,
		},
		Credentials: credentials,
	})
	if err != nil {
		t.Error("Error creating Account")
		return
	}

	_, err = ctrl.Token(context.TODO(), &pb.TokenRequest{
		Auth: credentials,
	})
	if err == nil {
		t.Error("Error is nil despite Account is disabled")
	} else if s, ok := status.FromError(err); !ok || (s.Code() != codes.PermissionDenied && s.Message() != "Account is disabled") {
		t.Errorf("Expected error 'Account is disabled', got: %v", err)
	}
}

func TestAuthorizeStandard(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation")
	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	credentials := &accounts.Credentials{
		Type: "standard",
		Data: []string{username, password},
	}

	crtRes, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: true,
		},
		Credentials: credentials,
	})
	if err != nil {
		t.Fatal("Error creating Account")
		return
	}

	res, err := ctrl.Token(context.TODO(), &pb.TokenRequest{
		Auth: credentials,
	})
	if err != nil {
		t.Fatalf("Unexpected error while getting Token: %v", err)
	}

	token, err := jwt.Parse(res.GetToken(), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing alg")
		}
		return ctrl.SIGNING_KEY, nil
	})
	if err != nil {
		t.Fatalf("Error parsing JWT: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Unexpected error while reading JWT Claims as Map")
	}

	account := claims[inf.INFINIMESH_ACCOUNT_CLAIM]
	if account == nil {
		t.Fatal("Account Claim is empty")
	}

	id, ok := account.(string)
	if !ok {
		t.Fatal("Error casting claim value to string")
	}

	if id != crtRes.Account.Uuid {
		t.Fatalf("Expected account in Claim to be %s, got: %s", crtRes.Account.Uuid, id)
	}
}

func TestAuthorizeStandardFail(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation with false Credentials")
	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)

	_, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: true,
		},
		Credentials: &accounts.Credentials{
			Type: "standard",
			Data: []string{username, password},
		},
	})
	if err != nil {
		t.Fatal("Error creating Account")
		return
	}

	_, err = ctrl.Token(context.TODO(), &pb.TokenRequest{
		Auth: &accounts.Credentials{
		Type: "standard",
		Data: []string{username, password + "blah"},
	},
	})
	if err == nil {
		t.Fatal("Token request supposed to fail, but it didn't")
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Can't parse Status from error, got: %v", err)
	}

	if s.Code() != codes.Unauthenticated || s.Message() != "Wrong credentials given" {
		t.Fatalf("Error supposed to be Unauthenticated: Wrong credentials given, but got %v: %v", s.Code().String(), s.Message())
	}
}

func TestUpdateAccount(t *testing.T) {
	t.Log("Creating sample account")

	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	this := &accounts.Account{
		Title: username, Enabled: false,
	}

	res, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: this,
		Credentials: &accounts.Credentials{
			Type: "standard",
			Data: []string{username, password},
		},
	})
	if err != nil {
		t.Fatalf("Error creating Account: %v", err)
		return
	}

	uuid := res.GetAccount().GetUuid()
	this.Uuid = uuid
	this.Title = username + "-new"
	this.Enabled = true
	
	that, err := ctrl.Update(rootCtx, this)
	if err != nil {
		t.Fatalf("Error udpating Account: %v", err)
	}
	if that != this {
		t.Fatal("Requested updates and updated accounts(from Response) aren't matching, this:", this, "that:", that)
	}

	_, err = ctrl.col.ReadDocument(rootCtx, uuid, that)
	if err != nil {
		t.Fatalf("Error reading Account in DB: %v", err)
	}
	if that != this {
		t.Fatal("Requested updates and updated accounts(from DB) aren't matching, this:", this, "that:", that)
	}
}