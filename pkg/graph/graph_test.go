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
		t.Fatalf("Error creating Account: %v", err)
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
		t.Fatalf("Error creating Account: %v", err)
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
		t.Fatalf("Error creating Account: %v", err)
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

func TestList(t *testing.T) {
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

	t.Logf("Created account: %s", res.GetAccount().GetUuid())

	pool, err := ctrl.List(rootCtx, nil)
	if err != nil {
		t.Fatalf("Error listing Account: %v", err)
	}

	if len(pool.Accounts) < 1 {
		t.Fatalf("Pool is empty, length: %d", len(pool.Accounts))
	}

	r := false
	for _, acc := range pool.Accounts {
		if acc.Uuid == res.GetAccount().GetUuid() {
			r = true
			break
		}
	}

	if !r {
		t.Fatalf("Account not found in pool")
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

func TestListNamespaces(t *testing.T) {
	title := randomdata.SillyName()
	nspb, err := ns_ctrl.Create(rootCtx, &namespaces.Namespace{
		Title: title,
	})
	if err != nil {
		t.Fatalf("Couldn't create Namespace: %v", err)
	}

	pool, err := ns_ctrl.List(rootCtx, &pb.EmptyMessage{})
	if err != nil {
		t.Fatalf("Couldn't list Namespace: %v", err)
	}

	rootFound, createdFound := false, false
	for _, ns := range pool.GetNamespaces() {
		if ns.GetUuid() == schema.ROOT_NAMESPACE_KEY {
			rootFound = true
		} else if ns.GetUuid() == nspb.GetUuid() {
			createdFound = true
			if ns.GetTitle() != nspb.GetTitle() {
				t.Logf("[WARNING]: namespaces titles don't match. Listed: %s; Created: %s", ns.GetTitle(), nspb.GetTitle())
			}
		}
	}
	if !rootFound {
		t.Fatal("Root Namespace not listed")
	}
	if !createdFound {
		t.Fatal("Created Namespace not listed")
	}
}

// Permissions Tests

func TestNewAccountNoNamespaceGiven(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation")
	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	credentials := &accounts.Credentials{
		Type: "standard",
		Data: []string{username, password},
	}

	accpb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: true,
		},
		Credentials: credentials,
	})
	if err != nil {
		t.Fatalf("Failed to create Account: %v", err)
	}
	acc := NewAccountFromPB(accpb.Account)

	edge := GetEdgeCol(rootCtx, db, schema.ACC2NS)
	ok := CheckLink(rootCtx, edge, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), acc)
	if !ok {
		t.Fatal("Account has to be under platform Namespace byt default")
	}
}

func TestNewAccountAccessToRoot(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation")
	username := randomdata.SillyName()
	password := randomdata.Alphanumeric(12)
	credentials := &accounts.Credentials{
		Type: "standard",
		Data: []string{username, password},
	}

	accPb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username, Enabled: true,
		},
		Credentials: credentials,
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatal("Error creating Account")
	}
	acc := NewAccountFromPB(accPb.Account)

	// Checking Account access to Root Account
	ok, level := AccessLevel(rootCtx, db, acc, NewBlankAccountDocument(schema.ROOT_ACCOUNT_KEY))
	if ok {
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", level, schema.NONE)
	}
}

/*
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░█████░░░░░░░░░░░░░░░░░█████░░░░░░░░░░
░░░██A██────────V────────██1██░░░░░░░░░░
░░░█████░░░░░░░░░░░░░░░░░█████░░░░░░░░░░
░░░░░│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░░│░░░░░░░ A - infinimesh NS ░░░░░░░░
░░░░░X░░░░░░░ 1 - User 1 ░░░░░░░░░░░░░░░
░░░░░│░░░░░░░ 2 - User 2 ░░░░░░░░░░░░░░░
░░░░░│░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░██2██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░█████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
*/

func TestPermissionsRootNamespace(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation")
	username1 := randomdata.SillyName()
	credentials1 := &accounts.Credentials{
		Type: "standard",
		Data: []string{username1, randomdata.Alphanumeric(12)},
	}

	// Create Account 1 under platform Namespace
	acc1pb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username1, Enabled: true,
		},
		Credentials: credentials1,
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatal("Error creating Account 1")
	}
	acc1 := NewAccountFromPB(acc1pb.Account)

	username2 := randomdata.SillyName()
	credentials2 := &accounts.Credentials{
		Type: "standard",
		Data: []string{username2, randomdata.Alphanumeric(12)},
	}

	// Create Account 2 under platform Namespace
	acc2pb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username2, Enabled: true,
		},
		Credentials: credentials2,
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatal("Error creating Account 2")
	}
	acc2 := NewAccountFromPB(acc2pb.Account)

	// Giving Account 1 Management access(MGMT) to Platform
	edge := GetEdgeCol(rootCtx, db, schema.ACC2NS)
	err = Link(rootCtx, log, edge, acc1, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), schema.MGMT)
	if err != nil {
		t.Fatalf("Error linking Account 1 to platform Namespace: %v", err)
	}

	// Checking Account 1 access to Account 2
	ok, level := AccessLevel(rootCtx, db, acc1, acc2)
	if !ok {
		t.Fatalf("Error checking Access or Access Level is 0(none)")
	}

	if level > int32(schema.MGMT) {
		t.Fatalf("Account 1 has higher access level than expected: %d(should be %d)", level, schema.MGMT)
	}
	if level < int32(schema.MGMT) {
		t.Fatalf("Account 1 has lower access level than expected: %d(should be %d)", level, schema.MGMT)
	}

	// Checking Account 2 access to Account 1
	ok, level = AccessLevel(rootCtx, db, acc2, acc1)
	if ok {
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", level, schema.NONE)
	}
}

func TestPermissionsRootNamespaceAccessAndGet(t *testing.T) {
	t.Log("Creating Sample Account and testing Authorisation")
	username1 := randomdata.SillyName()
	credentials1 := &accounts.Credentials{
		Type: "standard",
		Data: []string{username1, randomdata.Alphanumeric(12)},
	}

	// Create Account 1 under platform Namespace
	acc1pb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username1, Enabled: true,
		},
		Credentials: credentials1,
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatal("Error creating Account 1")
	}
	acc1 := NewAccountFromPB(acc1pb.Account)

	username2 := randomdata.SillyName()
	credentials2 := &accounts.Credentials{
		Type: "standard",
		Data: []string{username2, randomdata.Alphanumeric(12)},
	}

	// Create Account 2 under platform Namespace
	acc2pb, err := ctrl.Create(rootCtx, &accounts.CreateRequest{
		Account: &accounts.Account{
			Title: username2, Enabled: true,
		},
		Credentials: credentials2,
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatal("Error creating Account 2")
	}
	acc2 := NewAccountFromPB(acc2pb.Account)

	// Giving Account 1 Management access(MGMT) to Platform
	edge := GetEdgeCol(rootCtx, db, schema.ACC2NS)
	err = Link(rootCtx, log, edge, acc1, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), schema.MGMT)
	if err != nil {
		t.Fatalf("Error linking Account 1 to platform Namespace: %v", err)
	}

	nacc1 := *NewBlankAccountDocument(acc1.Key)
	nacc2 := *NewBlankAccountDocument(acc2.Key)
	// Checking Account 1 access to Account 2
	ok, level := AccessLevelAndGet(rootCtx, log, db, acc1, &nacc2)
	if !ok {
		t.Fatalf("Error checking Access or Access Level is 0(none)")
	}

	if level > int32(schema.MGMT) {
		t.Fatalf("Account 1 has higher access level than expected: %d(should be %d)", level, schema.MGMT)
	}
	if level < int32(schema.MGMT) {
		t.Fatalf("Account 1 has lower access level than expected: %d(should be %d)", level, schema.MGMT)
	}

	// Checking Account 2 access to Account 1
	ok, level = AccessLevelAndGet(rootCtx, log, db, &nacc2, &nacc1)
	if ok {
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", level, schema.NONE)
	}
}

func TestAccessLevelAndGetUnexistingAccountAndNode(t *testing.T) {
	acc1 := *NewBlankAccountDocument(randomdata.SillyName())
	acc2 := *NewBlankAccountDocument(randomdata.SillyName())

	ok, level := AccessLevelAndGet(rootCtx, log, db, &acc1, &acc2)
	if ok {
		t.Fatalf("Has to be error but it's not: %d", level)
	}
}