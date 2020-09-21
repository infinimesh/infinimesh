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

package dgraph

import (
	"context"
	"testing"

	"os"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

var repo node.Repo

func init() {
	dgURL := os.Getenv("DGRAPH_URL")
	if dgURL == "" {
		dgURL = "localhost:9080"
	}
	conn, err := grpc.Dial(dgURL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	repo = NewDGraphRepo(dg)
}

//test for Authorize
func TestAuthorize(t *testing.T) {
	ctx := context.Background()
	_, err := repo.CreateNamespace(ctx, "default")
	require.NoError(t, err)

	account, err := repo.CreateUserAccount(ctx, randomdata.SillyName(), "password", false, true)
	require.NoError(t, err)

	node, err := repo.CreateObject(ctx, "sample-node", "", "asset", "default")
	require.NoError(t, err)

	err = repo.Authorize(ctx, account, node, "READ", true)
	require.NoError(t, err)

	decision, err := repo.IsAuthorized(ctx, node, account, "READ")
	require.NoError(t, err)
	require.True(t, decision)

	//Update the account
	err = repo.UpdateAccount(context.Background(), &nodepb.UpdateAccountRequest{
		Account: &nodepb.Account{
			Uid:     account,
			Enabled: false,
			IsRoot:  false,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Enabled", "IsRoot"},
		},
	})
	require.NoError(t, err)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)
}

func TestIsAuthorizedNamespace(t *testing.T) {
	ctx := context.Background()

	accountname := randomdata.SillyName()
	account, err := repo.CreateUserAccount(ctx, accountname, "password", false, true)
	require.NoError(t, err)

	ns, err := repo.GetNamespace(ctx, accountname)
	require.NoError(t, err)

	decision, err := repo.IsAuthorizedNamespace(ctx, ns.Id, account, nodepb.Action_WRITE)
	require.NoError(t, err)
	require.True(t, decision)

	//Update the account
	err = repo.UpdateAccount(context.Background(), &nodepb.UpdateAccountRequest{
		Account: &nodepb.Account{
			Uid:     account,
			Enabled: false,
			IsRoot:  false,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Enabled", "IsRoot"},
		},
	})
	require.NoError(t, err)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)
}

func TestListInNamespaceForAccount(t *testing.T) {
	ctx := context.Background()

	acc := randomdata.SillyName()

	// Create Account
	account, err := repo.CreateUserAccount(ctx, acc, "password", false, true)
	require.NoError(t, err)

	//Get Namespace
	nsName, err := repo.GetNamespace(ctx, acc)

	//Create Object
	newObj, err := repo.CreateObject(ctx, "sample-node", "", "asset", nsName.Name)
	require.NoError(t, err)

	err = repo.AuthorizeNamespace(ctx, account, nsName.Id, nodepb.Action_WRITE)
	require.NoError(t, err)

	objs, err := repo.ListForAccount(ctx, account, nsName.GetName(), true)
	require.NoError(t, err)

	// Assert
	require.Contains(t, objs, &nodepb.Object{Uid: newObj, Name: "sample-node", Kind: "asset", Objects: []*nodepb.Object{}})

	//Update the account
	err = repo.UpdateAccount(context.Background(), &nodepb.UpdateAccountRequest{
		Account: &nodepb.Account{
			Uid:     account,
			Enabled: false,
			IsRoot:  false,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Enabled", "IsRoot"},
		},
	})
	require.NoError(t, err)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)
}

func TestChangePassword(t *testing.T) {
	ctx := context.Background()

	acc := randomdata.SillyName()

	// Create Account
	account, err := repo.CreateUserAccount(ctx, acc, "password", false, true)
	require.NoError(t, err)

	err = repo.SetPassword(ctx, account, "newpassword")
	require.NoError(t, err)

	ok, _, _, err := repo.Authenticate(ctx, acc, "newpassword")
	require.True(t, ok)

	//Update the account
	err = repo.UpdateAccount(context.Background(), &nodepb.UpdateAccountRequest{
		Account: &nodepb.Account{
			Uid:     account,
			Enabled: false,
			IsRoot:  false,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Enabled", "IsRoot"},
		},
	})
	require.NoError(t, err)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	ctx := context.Background()

	randomName := randomdata.SillyName()

	account, err := repo.CreateUserAccount(ctx, randomName, "password", true, true)
	require.NoError(t, err)

	//Set new values
	NewName := randomdata.SillyName()

	//Update the account
	err = repo.UpdateAccount(context.Background(), &nodepb.UpdateAccountRequest{
		Account: &nodepb.Account{
			Uid:     account,
			Name:    NewName,
			Enabled: false,
			IsRoot:  false,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"Name", "Enabled", "IsRoot"},
		},
	})
	require.NoError(t, err)

	//Get the updated Account Details
	respGet, err := repo.GetAccount(ctx, account)

	//Validate the updated Account
	require.NoError(t, err)
	require.EqualValues(t, NewName, respGet.Name)
	require.EqualValues(t, false, respGet.Enabled)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	ctx := context.Background()

	acc := randomdata.SillyName()

	// Create Account
	account, err := repo.CreateUserAccount(ctx, acc, "password", false, false)
	require.NoError(t, err)

	//Delete the Account created
	err = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})
	require.NoError(t, err)

	//Try to fetch the delete account
	_, err = repo.GetAccount(ctx, account)

	//Validation
	require.EqualValues(t, string(err.Error()), "The Account is not found")
}

func TestChangePasswordWithNoUser(t *testing.T) {
	ctx := context.Background()

	err := repo.SetPassword(ctx, "non-existing-user", "newpassword")
	require.Error(t, err)
}

func TestListPermissionsOnNamespace(t *testing.T) {
	ctx := context.Background()
	var accountID *nodepb.Account

	//Get All accounts
	accounts, err := repo.ListAccounts(ctx)

	//Find the required account
	for _, account := range accounts {
		if account.Name == "joe" {
			//Get Account
			accountID, err = repo.GetAccount(ctx, account.Uid)
			require.NoError(t, err)

		}
	}

	//Get Namespace
	nsID, err := repo.GetNamespace(ctx, "joe")

	err = repo.AuthorizeNamespace(ctx, accountID.Uid, nsID.Id, nodepb.Action_WRITE)
	require.NoError(t, err)

	permissions, err := repo.ListPermissionsInNamespace(ctx, nsID.Id)
	require.NoError(t, err)

	var namespaceFound bool
	for _, permission := range permissions {
		if permission.AccountName == nsID.Name {
			namespaceFound = true
		}
	}
	require.True(t, namespaceFound)
}

func TestDeletePermissionOnNamespace(t *testing.T) {
	ctx := context.Background()
	var accountID *nodepb.Account

	//Get All accounts
	accounts, err := repo.ListAccounts(ctx)

	//Find the required account
	for _, account := range accounts {
		if account.Name == "joe" {
			//Get Account
			accountID, err = repo.GetAccount(ctx, account.Uid)
			require.NoError(t, err)

		}
	}

	//Get Namespace
	nsID, err := repo.GetNamespace(ctx, "joe")

	err = repo.AuthorizeNamespace(ctx, accountID.Uid, nsID.Id, nodepb.Action_WRITE)
	require.NoError(t, err)

	err = repo.DeletePermissionInNamespace(ctx, nsID.Id, accountID.Uid)
	require.NoError(t, err)

	permissions, err := repo.ListPermissionsInNamespace(ctx, nsID.Id)
	require.NoError(t, err)
	require.Empty(t, permissions)

}

/*//Test to check API Endpoints

//Generic function to Perform HTTP request and return resopnse

func performRequest(r http.Handler, method, path string, body bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestUpdateAccountAPI(t *testing.T) {

	ctx := context.Background()
	acc := randomdata.SillyName()

	// Create Account
	account, err := repo.CreateUserAccount(ctx, acc, "password", false, true)
	require.NoError(t, err)

	//Set the JSON Body for the HTTP Request
	var jsonStr = []byte(`{"name":"Ankit"}`)

	//Set the request with the Method, path and the Json
	req, err := http.NewRequest("PATCH", "/accounts/"+account, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	//Set Http header
	req.Header.Set("Content-Type", "application/json")

	//Create the recorder for the request
	rr := httptest.NewRecorder()

	//Send the HTTP request to the endpoint
	handler := http.HandlerFunc(http.NewServeMux().ServeHTTP)
	handler.ServeHTTP(rr, req)

	//Delete the Account created
	_ = repo.DeleteAccount(ctx, &nodepb.DeleteAccountRequest{Uid: account})

	//assert.Equal(t, http.StatusBadRequest, w.Code)
	//assert.Equal(t, "{\"error\":\"Record not found!\"}", w.Body.String())
}
*/
