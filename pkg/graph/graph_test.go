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
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"github.com/infinimesh/infinimesh/pkg/node/proto/access"
	"github.com/infinimesh/infinimesh/pkg/node/proto/accounts"
	"github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"github.com/infinimesh/infinimesh/pkg/node/proto/namespaces"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
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

	ctrl *AccountsController
	ns_ctrl *NamespacesController
	dev_ctrl *DevicesController

	rootCtx context.Context

	db driver.Database
)

func init() {
	viper.AutomaticEnv()
	log = zap.NewExample()

	viper.SetDefault("DB_HOST", "db.infinimesh.local")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("INF_DEFAULT_ROOT_PASS", "infinimesh")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	rootPass := viper.GetString("INF_DEFAULT_ROOT_PASS")
	db = schema.InitDB(log, arangodbHost, arangodbCred, "infinimesh", false)
	
	ctrl = NewAccountsController(log, db)
	err := EnsureRootExists(log, db, rootPass)
	if err != nil {
		panic(err)
	}

	ns_ctrl = NewNamespacesController(log, db)
	dev_ctrl = NewDevicesController(log, db)

	md := metadata.New(map[string]string{
		inf.INFINIMESH_ACCOUNT_CLAIM: schema.ROOT_ACCOUNT_KEY,
	})
	rootCtx = metadata.NewIncomingContext(context.Background(), md)
	rootCtx = context.WithValue(rootCtx, inf.InfinimeshAccountCtxKey, schema.ROOT_ACCOUNT_KEY)
}

func CompareAccounts(a, b *accounts.Account) bool {
	return a.GetUuid() == b.GetUuid() &&
				 a.GetTitle() == b.GetTitle() &&
				 a.GetEnabled() == b.GetEnabled() &&
				 a.GetDefaultNamespace() == b.GetDefaultNamespace()
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

func TestUpdateAccountDefaultNS(t *testing.T) {
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
	this.DefaultNamespace = schema.ROOT_NAMESPACE_KEY
	
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

	this.DefaultNamespace = "notexistent"
	_, err = ctrl.Update(rootCtx, this)
	if err == nil {
		t.Fatalf("Error supposed to be raised, but it didn't")
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Can't parse Status from error, got: %v", err)
	}
	if s.Code() != codes.PermissionDenied {
		t.Fatalf("Error supposed to be PermissionDenied, but it was %v", s.Code().String())
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
	var _access Access
	_, err = edge.ReadDocument(rootCtx, schema.ROOT_ACCOUNT_KEY + "-" + nspb.Uuid, &_access)
	if err != nil {
		t.Fatalf("Can't read edge document or it doesn't exist: %v", err)
	}

	if _access.Level < access.Level_ADMIN {
		t.Fatalf("Access level incorrect(%d), must be: %d", _access.Level, access.Level_ADMIN)
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

	edge := GetEdgeCol(rootCtx, db, schema.NS2ACC)
	ok := CheckLink(rootCtx, edge, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), acc)
	if !ok {
		t.Fatal("Account has to be under platform Namespace by default")
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
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", level, access.Level_NONE)
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
	err = Link(rootCtx, log, edge, acc1, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), access.Level_MGMT, access.Role_UNSET)
	if err != nil {
		t.Fatalf("Error linking Account 1 to platform Namespace: %v", err)
	}

	// Checking Account 1 access to Account 2
	ok, level := AccessLevel(rootCtx, db, acc1, acc2)
	if !ok {
		t.Fatalf("Error checking Access or Access Level is 0(none)")
	}

	if level > access.Level_MGMT {
		t.Fatalf("Account 1 has higher access level than expected: %d(should be %d)", level, access.Level_MGMT)
	}
	if level < access.Level_MGMT {
		t.Fatalf("Account 1 has lower access level than expected: %d(should be %d)", level, access.Level_MGMT)
	}

	// Checking Account 2 access to Account 1
	ok, level = AccessLevel(rootCtx, db, acc2, acc1)
	if ok {
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", level, access.Level_NONE)
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
	err = Link(rootCtx, log, edge, acc1, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY), access.Level_MGMT, access.Role_UNSET)
	if err != nil {
		t.Fatalf("Error linking Account 1 to platform Namespace: %v", err)
	}

	nacc1 := *NewBlankAccountDocument(acc1.Key)
	nacc2 := *NewBlankAccountDocument(acc2.Key)
	// Checking Account 1 access to Account 2
	err = AccessLevelAndGet(rootCtx, log, db, acc1, &nacc2)
	if err != nil {
		t.Fatalf("Error checking Access or Access Level is 0(none)")
	}

	if nacc2.Access.Level > access.Level_MGMT {
		t.Fatalf("Account 1 has higher access level than expected: %d(should be %d)", nacc2.Access.Level, access.Level_MGMT)
	}
	if nacc2.Access.Level < access.Level_MGMT {
		t.Fatalf("Account 1 has lower access level than expected: %d(should be %d)", nacc2.Access.Level, access.Level_MGMT)
	}

	// Checking Account 2 access to Account 1
	err = AccessLevelAndGet(rootCtx, log, db, &nacc2, &nacc1)
	if err == nil && nacc1.Access.Level > access.Level_NONE {
		t.Fatalf("Account 2 has higher access level than expected: %d(should be %d)", nacc1.Access.Level, access.Level_NONE)
	}
}

func TestAccessLevelAndGetUnexistingAccountAndNode(t *testing.T) {
	acc1 := *NewBlankAccountDocument(randomdata.SillyName())
	acc2 := *NewBlankAccountDocument(randomdata.SillyName())

	err := AccessLevelAndGet(rootCtx, log, db, &acc1, &acc2)
	if err == nil {
		t.Fatalf("Has to be error but it's not: %v", err)
	}
}

// Devices Tests

func TestCreateGetAndDelete(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUAhdObDJnggMc26EDvd11ItBaPewwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0xOTAzMzAyMTE2MjhaFw0yMDAz
MjkyMTE2MjhaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQCsqC6OCmA9bG8qA/qo8W8afBYANo4YKrcohAY9Mg+e
rabrPxS3fVCDUEn4A7EKlHIzRuobxVvUSOK25rGC0SftQmaRhqzhtV5oL84EvKzy
7Q83gCDZbocsf+GqdZxEayFm89WNQCCMr1BIRwVY7w8pElgAPVujRisjE6MV2//z
nQwigxuNGA06fFmT+IjCy7p3owhPanC26yn0KaFSPWqsxvl5/gD21rAJpYKtOf6V
Xl85K4tjXCqJ00oPZa6y5y1eSyEb9UE225LsBFoARgcoNCcwVpRdkiXt8TD1KI7K
LNhqOxeRi3y2i9e0f1vgo6BX6xCUx14tRs5nJyfEp7Nknyg7yK8a7FZ5eazRuRlS
CrySP3BjxEZTVySW9wxDaEY250vxtcFJanccD5/C/BAHwAN6eLczgCm4JXWA6HNq
G0LNouAONj9u5XY6gfJRz3Iyh4bHwMmXH4gN9Ep6MeeXlN13B8zQ9udG15UwNELl
FJ24KhRj+hAX8uEHeBLZN4hUKzwpkBeN9r/aotZsdVFAp9Gl1RrpOKwaShuh1fY3
nT+P0CECB4ugoIkuY/IgIGNwGzdBeAZ+pCxUuq95eGKFg0wdA7aJtV11IdcxdpYU
ljJlFlluc/CSY6G7I626nrYumLK8Lt1GCmAKu7fnZ02gKGJvbk7VEjj7/M7BywHm
XwIDAQABo1MwUTAdBgNVHQ4EFgQUybknD8HpuXXis6Jq6nT2H3R6I90wHwYDVR0j
BBgwFoAUybknD8HpuXXis6Jq6nT2H3R6I90wDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEAUHlDMCMjifD8M+7MTejca9dYOXK3kujgFjKnSig2Fa4q
97n9/9I+r3AAGz9YFBPw1bmHlnE7ieV1NJuCSe+L3z1s2Cq5ulQZV1qiP3NFdjUd
J1P8TBI18aA38huX4RwAUn1sA8B0ApiwRe8f3JxgvmCx3T9JSWaUavrMNUKM8aXf
gv6FrtrxKcE5MfdS72GYqY2+zy8HQQC0aVNhkoZk6GhlVLTG1HrZ22b9zQpn+pRc
adwWZ752ZQMSUrntPYA6BK/4aVl65m47dA0k1EnK8Hz5DSEr5W+NP1m8xsw3yFKo
1NQB0RULPvMiHqb3czMKY6ORMNZOvYqyUlbErjCcAc9jbvgOzmT91GFf6u3mUgmm
z8QvMTFzuAtmhmyKSjtiLmHZvhxQBXr2fX+7Zryz139EGijSCZOspwFHZkJFkg/U
Sc0vud9uKcbklt5GluM7HaxZWPTdze03Z4wbvU1nhKyGdf4bTvM3HeqAK2sURN3Y
1uxEV4k4S++gj2UUcfxtHGRiF3YGEfzt7pmigVo/Cr73ZXzPmr3qhxDyBSvyAARl
O57+cPhsWazbxf8tHDvSJEO+bVtFY4JRoLi2TpNbug75LmBFHojgxqgI/3qogJ2Q
lKUKOFnVNnDJhVjLh5DeNYbyyU/f+xFqTzQjSyuw+4FegKbzfq7oliRnDeT0Wgs=
-----END CERTIFICATE-----
`

	thisR, err := dev_ctrl.Create(rootCtx, &devices.CreateRequest{
		Device: &devices.Device{
			Title: randomdata.SillyName(),
			Enabled: true,
			Certificate: &devices.Certificate{
				PemData: cert,
			},
		},
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}
	this := thisR.Device

	t.Logf("Device created: %s", this.GetUuid())

	that, err := dev_ctrl.Get(rootCtx, this)
	if err != nil {
		t.Fatalf("Error getting device: %v", err)
	}

	if this.Uuid != that.Uuid	{
			t.Fatalf("Devices aren't same. %s != %s", this.Uuid, that.Uuid)
	}
	if this.Title != that.Title {
			t.Fatalf("Devices aren't same. %s != %s", this.Title, that.Title)
	}
	if this.Enabled != that.Enabled {
			t.Fatalf("Devices aren't same. %t != %t", this.Enabled, that.Enabled)
	}
	thisc := string(this.Certificate.Fingerprint)
	thatc := string(that.Certificate.Fingerprint)
	if thisc != thatc {
			t.Fatalf("Devices aren't same. %s != %s", thisc, thatc)
	}

	_, err = dev_ctrl.Delete(rootCtx, this)
	if err != nil {
		t.Fatalf("Error deleting device: %v", err)
	}
}

func TestCreateAndList(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIIEljCCAn4CCQC7oNynkLPhTjANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJk
ZTAeFw0yMTA2MTYxMTMyNDRaFw0yMjA2MTYxMTMyNDRaMA0xCzAJBgNVBAYTAmRl
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0hk6i+PxRW7XAy21QAsR
Dlyz60ojkDU5q2BfXzmo5GPGaAXuEwwT+AJGFAgIvSIWh7SBDY3re75YbShfbLEP
biHDtNKzr0v+RmNiZ66qZy7lVPyTcDe4Aj9iOsdAiocKXBECgpdvDPM2SPVsL915
oajg2RAp/VmvtHdENBjgD0e7xVXV4hKwn2UDMQbw1KBfIXVj6n7fwMvouovcmdc+
A107+HTudDqvhrkevAJXDmxTWRKz3anoU/dCcV4d1aHLys29L/vnlF0q29KEfSLJ
Ov9H/9mX/NjcmMqr4tsqjmu5ZepORhtGqq0Rmcg++FbCA4f68OchTPopvYKz7ExN
CPzgxufqduBdThIwNzdtXctm0othphQ3ADxnxCqDfAhqr02w7qaCd/c1KBK6EKvJ
uIWiqaVV3ipqre+T98AuzJ7il+mhIsRsXpBt3o7LBCgyl8rri+ZLEDRj3hOu3UN5
pS71R0xm62P8psKY0xtDneReUQ1CGObQS7XZDCJ0qlHDGUMTBwvGbcqrTwpA1udu
cP1GGDhRsdlx0NgJEemSojEiMKmSc1McNsubczfJCZAZRNNvR7pn4MyyS20aMNnd
1rRkX6ikyvRA96dJD0M4iI2f6asNpGe8SplwPJweNv/avwYiWKFVO5neuVEdiAcw
XjFL9u8OK0ID8Uid3TWV4psCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEALKx4BlYg
dizAl5jVICrswgVlS/Ec8dw3hTmuDodhA5jP5NLFIrzWHp6voythjhFIdXHI+8nW
y0V1TVviW73qFP9ib5LnLn30QVajwFRjBIOt4qsrIvMFDvwtQ940pUgR1iVGphV4
ahlCwNeZStdxMV8M4/5o78wP7uvyhleIaYrF7dLfFoszT4PfyRC2UEXtTknz1hH8
kOFwiZCio5sIzWNsAzHlOKbf2Rl0WtC9YWcKpdS1MrWi6E/jAJQ1/GyhUOEZHE/Z
fY1heN2YXPacYtFQTRmkp/oPzsIvwgfx6OKJe8RGa7EErQUVGTMYkZue7lpIOyJD
8m37TUVNizW2+OrQb/NUK9uwEBkGlpavTdK7eKAw0+KnlPqMpmQx7Vs5oE0ejy7y
GuMpc8AeJXUX9lHMJIT+lwkKzrVReC+jgyvO0QyRN7PTwRW8+9SNOeHRiC9Fj7Zg
fLCCa/hdALN6ECHn3JsQGiAbY6JS8LOdiLpnlR+cOQSQ3HnaBkpPeBmWfRvlvGeU
r+vyP3YimFBE9AbM5GgfUHGRBJBpC40aVaE7HtHapE4JJNit4NfBvfDotNUs6shJ
6Y893NPueYB4PfvC+1kgZFjXFEMDURaGUeEwl481Zn/rGXM4ev5qGPQgJ4fhmI68
cgSqKFgDFRxlHXLo9TZnxyBrIvN/siE+ZQI=
-----END CERTIFICATE-----`

	thisR, err := dev_ctrl.Create(rootCtx, &devices.CreateRequest{
		Device: &devices.Device{
			Title: randomdata.SillyName(),
			Enabled: true,
			Certificate: &devices.Certificate{
				PemData: cert,
			},
		},
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}
	this := thisR.Device

	t.Logf("Device created: %s", this.GetUuid())

	pool, err := dev_ctrl.List(rootCtx, nil)
	if err != nil {
		t.Fatalf("Error listing devices: %v", err)
	}

	this_found := false
	for _, dev := range pool.GetDevices() {
		if dev.Uuid == this.Uuid {
			this_found = true
			break
		}
	}

	if !this_found {
		t.Fatalf("Created device not found in Pool")
	}
}

func TestCreateFingByFingerprintAndDelete(t *testing.T) {
	cert := `-----BEGIN CERTIFICATE-----
MIIEmDCCAoACCQDLeCKlPBA5IzANBgkqhkiG9w0BAQsFADAOMQwwCgYDVQQDDANk
ZXYwHhcNMjEwNjE2MTEzMzE4WhcNMjIwNjE2MTEzMzE4WjAOMQwwCgYDVQQDDANk
ZXYwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC8XcvgKyPfTnGqP1Y0
GXpKj63/A285MMudwNHk59X38RSwrBl4IpWd43w6BBSbVBYYJ4lDsRDnQLjTelEa
BFlLrM1ZPtlh6qsYwcPRgHpujw7ufifhuKbtMPCz2IomyzzGFFKY3d+oJ7hRQe8m
fDcIpqrLCiuc2zuGLjVTEueFStJBdXciDRNeY9ILTHCpnZ7XNx8EsDylli3h5WOt
IsINB7osLmcnhsuvD594IQ2CBLUjOfeQDrkAWGxQ1DvaN1u7HCP84SXJ7nWdfma3
aYfesrDo+mlD70maucD1CemklY8yzgNxuFFIzQ+L4Klx7cjujqPd6XZmJjm1LTj7
ITHbbM9uzawx6+0591uOZivEgPB02b+92iAx5x1yYC5tjnTMy2P3TYrMMcKpI+uh
0KruEYWPKNuFOjcV6svoZEVeyTOLRjj63xkJ50HFrdE2Xcj70nPuB482B2jmcuOG
17ThaHTWd7paT/I5pgyQgRVTHmyRcqR6MfR6ofTXI76U5RBAKe6vrfVdNKNQ5J/T
uEygaIZCP+/Lq3ydS4/UNP0E7NgrkODG9i8FD1DF4wxneX2c8gozh5rVKuP2yS0r
+qgIH5AK9pKDNSL7J6P2ZIGeuwc2B800wGQhRCNMQUiz/llAVypqrJWUxBtN9v8k
CDMYxHQ4EBLfOOXLw8Youf/SpQIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQAGe3jr
tvC4auIpVFspS1KnRhaKxWootuMFrfo7yZQsSDYpCV3iM4RSwovsyn5xZvRJWO8z
7Gj5h4ZbmEFYsyNo9tbKeKWQjn3QVK7UlwMCjwZTYsxpxCioQ66XqjnSjzfTKFcm
CFOpP8nkR6mgxjDyQcqsDQX4vrUrt8PtwSag7r7+xl9MJ1pBKUaAmgGJAtjLxDyq
LDSnXI4/gTTKyXCxHzxgcioVz4j1gNtyFjPeDNkgOCuLhxBk1ewmR3m6swwMReiL
APdZLak2EPunZCYTG5648xYUowwkBSINQmGWS3YbuC0xncy/EhEuBS4mbsd7uO5w
m0HNT/FPHfoZnYS7eUOj42ER1q0JmPJYkgtJMwrNylF6+djrXbVLMeimh6ME2mA9
oROBNInt3vw6Ssd3kyQBMurh8ETu01Dj9MSzFcoX8293FsYbj9H/FndBf9I/a/UK
+iEPErgUiy8x+5qMRGeZrqHtfOdcHuliSJ0pS207nVdmMHUmHXm1LW3v+cScF+13
EqW8wfH8nuubsLAgpxx4s6hin9wjs9a27fAPEUPzFNmXs5SZF6+dGTbUtmd0Zp84
a+5z88Oa1aswXQBRt+4JTHJsc5KE2/pWuZY6+CL738hzWmDYpr3JHV1HdAN3dHU1
UWjgQjqXqHAguCY1KKG8lyzY3Q9pkmJcoy0HiA==
-----END CERTIFICATE-----`

	thisR, err := dev_ctrl.Create(rootCtx, &devices.CreateRequest{
		Device: &devices.Device{
			Title: randomdata.SillyName(),
			Enabled: true,
			Certificate: &devices.Certificate{
				PemData: cert,
			},
		},
		Namespace: schema.ROOT_NAMESPACE_KEY,
	})
	if err != nil {
		t.Fatalf("Error creating device: %v", err)
	}
	this := thisR.Device
	t.Logf("Device created: %s", this.GetUuid())

	that, err := dev_ctrl.GetByFingerprint(rootCtx, &devices.GetByFingerprintRequest{
		Fingerprint: this.Certificate.Fingerprint,
	})

	if err != nil {
		t.Fatalf("Error getting device: %v", err)
	}

	if this.Uuid != that.Uuid	{
			t.Fatalf("Devices aren't same. %s != %s", this.Uuid, that.Uuid)
	}
	if this.Title != that.Title {
			t.Fatalf("Devices aren't same. %s != %s", this.Title, that.Title)
	}
	if this.Enabled != that.Enabled {
			t.Fatalf("Devices aren't same. %t != %t", this.Enabled, that.Enabled)
	}
	thisc := string(this.Certificate.Fingerprint)
	thatc := string(that.Certificate.Fingerprint)
	if thisc != thatc {
			t.Fatalf("Devices aren't same. %s != %s", thisc, thatc)
	}

	_, err = dev_ctrl.Delete(rootCtx, this)
	if err != nil {
		t.Fatalf("Error deleting device: %v", err)
	}
}

func TestFingByFingerprintNotFound(t *testing.T) {
	_, err := dev_ctrl.GetByFingerprint(rootCtx, &devices.GetByFingerprintRequest{
		Fingerprint: []byte("notfound"),
	})

	if err == nil {
		t.Fatalf("Expected error")
	}

	s, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Error reading status from error, original error: %v", err)
	}

	if s.Code() != codes.NotFound && s.Message() != ("Device not found") {
		t.Fatalf("Error supposed to be NotFound: The device does not exist, but received %s: %s", s.Code().String(), s.Message())
	}
}

// TODO: Automate this test
// func TestListOwnedDeep(t *testing.T) {
// 	r, err := ListOwnedDeep(context.Background(), log, db, NewBlankAccountDocument("infinimesh"))
// 	if err != nil {
// 		t.Fatalf("Error listing nodes: %v", err)
// 	}

// 	log.Info("Retrieved nodes", zap.Any("nodes", r))
// }

// TODO: Automate this test
// func TestDeleteRecursive(t *testing.T) {
// 	err := DeleteRecursive(context.Background(), log, db, NewBlankAccountDocument("infinimesh"))
// 	if err != nil {
// 		t.Fatalf("Error deleting nodes: %v", err)
// 	}
// }