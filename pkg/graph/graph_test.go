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
	"github.com/arangodb/go-driver"
	"github.com/golang-jwt/jwt/v4"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/internal"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"github.com/infinimesh/infinimesh/pkg/node/proto/accounts"
	"github.com/infinimesh/infinimesh/pkg/node/proto/namespaces"
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
	ns_ctrl NamespacesController

	rootCtx context.Context

	db driver.Database
)

func init() {
	viper.AutomaticEnv()
	log = zap.NewExample()

	viper.SetDefault("DB_HOST", "localhost:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("ROOT_PASS", "infinimesh")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	rootPass := viper.GetString("ROOT_PASS")
	db = schema.InitDB(log, arangodbHost, arangodbCred, "infinimesh")
	
	ctrl = NewAccountsController(log, db)
	err := ctrl.EnsureRootExists(rootPass)
	if err != nil {
		panic(err)
	}

	ns_ctrl = NewNamespacesController(log, db)

	md := metadata.New(map[string]string{"requestorid": schema.ROOT_ACCOUNT_KEY})
	rootCtx = metadata.NewIncomingContext(context.Background(), md)
}

func CompareAccounts(a, b *accounts.Account) bool {
	return a.GetUuid() == b.GetUuid() &&
				 a.GetTitle() == b.GetTitle() &&
				 a.GetEnabled() == b.GetEnabled()
}

// AccountsController Tests

func TestNewBlankAccountDocument(t *testing.T) {
	uuid := randomdata.StringNumber(10, "-")
	uuidMeta := driver.NewDocumentID(schema.ACCOUNTS_COL, uuid)
	acc := NewBlankAccountDocument(uuid)
	if acc.ID() != uuidMeta {
		t.Fatalf("Blank document meta ID not equal to given. Comparing %v with %v", acc.ID(), uuidMeta)
	}
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

func TestAccountCreate_FalseCredentialsType(t *testing.T) {
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
		t.Fatal("Error creating Account")
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
	}

	uuid := res.GetAccount().GetUuid()
	this.Uuid = uuid
	this.Title = username + "-new"
	this.Enabled = true
	
	that, err := ctrl.Update(rootCtx, this)
	if err != nil {
		t.Fatalf("Error udpating Account: %v", err)
	}
	if !CompareAccounts(this, that) {
		t.Fatal("Requested updates and updated accounts(from Response) aren't matching, this:", this, "that:", that)
	}

	_, err = ctrl.col.ReadDocument(rootCtx, uuid, that)
	if err != nil {
		t.Fatalf("Error reading Account in DB: %v", err)
	}
	if !CompareAccounts(this, that) {
		t.Fatal("Requested updates and updated accounts(from DB) aren't matching, this:", this, "that:", that)
	}
}

func TestGetAccount(t *testing.T) {
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
	}

	uuid := res.GetAccount().GetUuid()
	this.Uuid = uuid
	that, err := ctrl.Get(rootCtx, &accounts.Account{Uuid: uuid})
	if err != nil {
		t.Fatalf("Error getting Account from API: %v", err)
	}
	if !CompareAccounts(this, that) {
		t.Fatal("Requested and created accounts(from API) aren't matching, this:", this, "that:", that)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	r, err := ctrl.Get(rootCtx, &accounts.Account{Uuid: randomdata.Alphanumeric(12)})
	if err == nil {
		t.Fatal("Get account received no error despite it should, response:", r)
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Can't parse Status from error, got: %v", err)
	}

	if s.Code() != codes.NotFound {
		t.Fatalf("Error supposed to be NotFound, but got %v", s.Code().String())
	}
}

func TestDeleteAccount(t *testing.T) {
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
	}

	this.Uuid = res.Account.GetUuid()
	_, err = ctrl.Delete(rootCtx, this)
	if err != nil {
		t.Fatalf("Unexpected error while deleting Account: %v", err)
	}

	r, err := ctrl.Get(rootCtx, this)
	if err == nil {
		t.Fatal("Get account received no error despite it should, response:", r)
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Can't parse Status from error, got: %v", err)
	}

	if s.Code() != codes.NotFound {
		t.Fatalf("Error supposed to be NotFound, but got %v", s.Code().String())
	}
}

func TestSetCredentialsStandard(t *testing.T) {
	t.Log("Creating sample account")

	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	this := &accounts.Account{
		Title: username, Enabled: true,
	}

	crtRes, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: this,
		Credentials: &accounts.Credentials{
			Type: "standard",
			Data: []string{username, password},
		},
	})
	if err != nil {
		t.Fatalf("Error creating Account: %v", err)
	}
	this.Uuid = crtRes.GetAccount().GetUuid()

	_, err = ctrl.SetCredentials(rootCtx, &pb.SetCredentialsRequest{
		Uuid: this.GetUuid(), Credentials: &accounts.Credentials{
			Type: "standard",
			Data: []string{username, password + "-addon"},
		},
	})
	if err != nil {
		t.Fatalf("Error Setting New Credentials: %v", err)
	}

	res, err := ctrl.Token(context.TODO(), &pb.TokenRequest{
			Auth: &accounts.Credentials{
			Type: "standard",
			Data: []string{username, password + "-addon"},
		},
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

	if id != this.GetUuid() {
		t.Fatalf("Expected account in Claim to be %s, got: %s", this.GetUuid(), id)
	}
}

// NamespacesController Tests

func TestCreateNamespace(t *testing.T) {
	title := randomdata.SillyName()
	nspb, err := ns_ctrl.Create(rootCtx, &namespaces.Namespace{
		Title: title,
	})
	if err != nil {
		t.Fatalf("Couldn't create Namespace: %v", err)
	}

	ok, err := ns_ctrl.col.DocumentExists(rootCtx, nspb.Uuid)
	if err != nil {
		t.Fatalf("Error testing Namespace existance: %v", err)
	} else if !ok {
		t.Fatalf("Namespace doesn't exist in DB")
	}

	edge := GetEdgeCol(rootCtx, db, schema.ACC2NS)
	var access Access
	_, err = edge.ReadDocument(rootCtx, schema.ROOT_ACCOUNT_KEY + "-" + nspb.Uuid, &access)
	if err != nil {
		t.Fatalf("Can't read edge document or it doesn't exist: %v", err)
	}

	if access.Level < 3 {
		t.Fatalf("Access level incorrect(%d), must be: %d", access.Level, schema.ADMIN)
	}
}
