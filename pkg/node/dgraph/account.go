//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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
	"encoding/json"
	"errors"
	"strings"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slntopp/infinimesh/pkg/node/nodepb"
)

//UserExists is a method to check if the user is present in the Dgraph DB
func (s *DGraphRepo) UserExists(ctx context.Context, account string) (exists bool, err error) {

	txn := s.Dg.NewReadOnlyTxn()

	q := `query userExists($accid: string) {
		exists(func: uid($accid)) @filter(eq(type, "account") and eq(enabled,"true"))  {
			uid
			enabled
		}
	}
	`

	var result struct {
		Exists []map[string]interface{} `json:"exists"`
	}

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$accid": account})
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false, err
	}

	if len(result.Exists) == 0 {
		return false, errors.New("Account not found")
	}

	return true, nil
}

//ListAccounts is a method to List details of all Account
func (s *DGraphRepo) ListAccounts(ctx context.Context) (accounts []*nodepb.Account, err error) {
	txn := s.Dg.NewReadOnlyTxn()

	const q = `query accounts{
                     accounts(func: eq(type, "account")) {
                       uid
                       name
                       enabled
					   isRoot
					   isAdmin
                     }
                   }`

	res, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var result struct {
		Accounts []*Account `json:"accounts"`
	}

	if err := json.Unmarshal(res.Json, &result); err != nil {
		return nil, err
	}

	for _, account := range result.Accounts {
		accounts = append(accounts, &nodepb.Account{
			Uid:     account.UID,
			Name:    account.Name,
			IsRoot:  account.IsRoot,
			Enabled: account.Enabled,
			IsAdmin: account.IsAdmin,
		})
	}

	return accounts, nil
}

//ListAccountsforAdmin is a method to List accounts owned by an Admin user
func (s *DGraphRepo) ListAccountsforAdmin(ctx context.Context, requestorID string) (accounts []*nodepb.Account, err error) {
	txn := s.Dg.NewReadOnlyTxn()

	const q = `query listaccountsforadmin($accountid: string) {
		accounts(func: uid($accountid))  {
			
		  owns @filter(eq(type, "account")) {
			uid
			name
			isRoot
			isAdmin
			enabled
		  }
		}
	  }`

	res, err := txn.QueryWithVars(ctx, q, map[string]string{"$accountid": requestorID})
	if err != nil {
		return nil, err
	}

	var result struct {
		Accounts []*Account `json:"accounts"`
	}

	if err := json.Unmarshal(res.Json, &result); err != nil {
		return nil, err
	}

	if len(result.Accounts) == 0 {
		return nil, errors.New("The Admin doesnot own any accounts")
	}

	for _, account := range result.Accounts[0].Owns {
		accounts = append(accounts, &nodepb.Account{
			Uid:     account.UID,
			Name:    account.Name,
			IsRoot:  account.IsRoot,
			IsAdmin: account.IsAdmin,
			Enabled: account.Enabled,
		})
	}

	return accounts, nil
}

//UpdateAccount is a method to Udpdate details of an Account
func (s *DGraphRepo) UpdateAccount(ctx context.Context, account *nodepb.UpdateAccountRequest, isself bool) (err error) {

	txn := s.Dg.NewTxn()

	//Get the data for the Account that is to be modified
	tempacc, err := s.GetAccount(ctx, account.Account.Uid)
	if err != nil {
		return err
	}

	//Build the account data in json for updating the account
	acc := &Account{
		Node: Node{
			Type: "account",
			UID:  account.Account.Uid,
		},
		Name:    tempacc.Name,
		IsRoot:  tempacc.IsRoot,
		IsAdmin: tempacc.IsAdmin,
		Enabled: tempacc.Enabled,
	}

	//Update the account based on the flag isself
	//isself if true means the requestor id updating its own account
	//Updating own account means only updating the name and the password
	for _, field := range account.FieldMask.Paths {
		switch strings.ToLower(field) {
		case "name":
			acc.Name = account.Account.Name
			//If the account name is updated make sure the default Namespace for the account is also updated
			err = s.UpdateNamespace(ctx, &nodepb.UpdateNamespaceRequest{
				Namespace: &nodepb.Namespace{
					Id:   tempacc.DefaultNamespace.Id,
					Name: account.Account.Name,
				},
				NamespaceMask: &field_mask.FieldMask{
					Paths: []string{"Name"},
				},
			})
		case "is_root":
			if !isself {
				acc.IsRoot = account.Account.IsRoot
			}
		case "is_admin":
			if !isself {
				acc.IsAdmin = account.Account.IsAdmin
			}
		case "enabled":
			if !acc.IsRoot && !isself {
				acc.Enabled = account.Account.Enabled
			}
		case "password":
			err = s.SetPassword(ctx, account.Account.Uid, account.Account.Password)
			if err != nil {
				return err
			}
		}
	}

	js, err := json.Marshal(acc)
	if err != nil {
		return err
	}

	m := &api.Mutation{SetJson: js}
	_, err = txn.Mutate(ctx, m)
	if err != nil {
		return err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return errors.New("Failed to commit")
	}

	return nil
}

//CreateUserAccount is a method to Create User Account for Infinimesh
func (s *DGraphRepo) CreateUserAccount(ctx context.Context, username, password string, isRoot, isAdmin, enabled bool) (uid string, err error) {
	txn := s.Dg.NewTxn()

	q := `query userExists($name: string) {
                exists(func: eq(name, $name)) @filter(eq(type, "account")) {
                  uid
                }
              }
             `

	var result struct {
		Exists []map[string]interface{} `json:"exists"`
	}

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$name": username})
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return "", err
	}

	//If the user doesnt exists then create user
	if len(result.Exists) == 0 {

		//Create Default namespace for the new user being created
		defaultNs, err := s.CreateNamespace(ctx, username)
		if err != nil {
			return "", err
		}

		//Build the json data structure to create the user for DGraph
		js, err := json.Marshal(&Account{
			Node: Node{
				Type: "account",
				UID:  "_:user",
			},
			Name:    username,
			IsRoot:  isRoot,
			IsAdmin: isAdmin,
			Enabled: enabled,
			HasCredentials: []*UsernameCredential{
				{
					Node: Node{
						Type: "credentials",
					},
					Username: username,
					Password: password,
				},
			},
			AccessToNamespace: []*Namespace{
				&Namespace{
					Node: Node{
						UID: defaultNs,
					},
					AccessToPermission: nodepb.Action_WRITE.String(),
				},
			},
			DefaultNamespace: []*Namespace{
				&Namespace{
					Node: Node{
						UID: defaultNs,
					},
				},
			},
		})
		if err != nil {
			return "", err
		}

		//Create the new user in the Dgraph DB
		m := &api.Mutation{SetJson: js}
		a, err := txn.Mutate(ctx, m)
		if err != nil {
			return "", err
		}

		err = txn.Commit(ctx)
		if err != nil {
			return "", err
		}
		userUID := a.GetUids()["user"]
		return userUID, nil

	}

	//If the user exists then return error
	return "", errors.New("User exists already")
}

//GetAccount is a method to Get details of an Account
func (s *DGraphRepo) GetAccount(ctx context.Context, name string) (account *nodepb.Account, err error) {
	txn := s.Dg.NewReadOnlyTxn()
	const q = `query accounts($account: string) {
                     accounts(func: uid($account)) @filter(eq(type, "account"))  {
                       uid
                       name
                       type
					   isRoot
					   isAdmin
                       enabled
                       default.namespace {
                         name
                         uid
					   }
					   has.credentials{
						uid
						username
					   }
					   ~owns{
						uid
					   }
                     }
                   }`

	response, err := txn.QueryWithVars(ctx, q, map[string]string{"$account": name})
	if err != nil {
		return nil, err
	}

	var result struct {
		Account []*Account `json:"accounts"`
	}

	err = json.Unmarshal(response.Json, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Account) == 0 {
		return nil, errors.New("The Account is not found")
	}

	account = &nodepb.Account{
		Uid:     result.Account[0].UID,
		Name:    result.Account[0].Name,
		IsRoot:  result.Account[0].IsRoot,
		IsAdmin: result.Account[0].IsAdmin,
		Enabled: result.Account[0].Enabled,
	}

	//Added default namespace details is present
	if len(result.Account[0].DefaultNamespace) == 1 {
		account.DefaultNamespace = &nodepb.Namespace{
			Name:            result.Account[0].DefaultNamespace[0].Name,
			Id:              result.Account[0].DefaultNamespace[0].UID,
			Markfordeletion: false,
		}
	}

	//Added credetnials details is present
	if len(result.Account[0].HasCredentials) == 1 {
		account.Username = result.Account[0].HasCredentials[0].Username
	}

	//Added owner details is present
	if len(result.Account[0].OwnedBy) == 1 {
		account.Owner = result.Account[0].OwnedBy[0].UID
	}

	return account, err
}

//DeleteAccount is a method to delete the Account
func (s *DGraphRepo) DeleteAccount(ctx context.Context, request *nodepb.DeleteAccountRequest) (err error) {
	txn := s.Dg.NewTxn()
	m := &api.Mutation{CommitNow: true}

	//Query to get the Account to be deleted with all the related edges
	//conidition for deleting an account
	//1. Account should not be named root
	//2. Account should have isRoot flag as false
	//3. Account should have enabled flag as false.
	const q = `query delete($userid: string){
		account(func: uid($userid)) @filter(eq(type, "account") and  eq(isRoot,false) and Not eq(name,"root") and eq(enabled,false)) {
			uid
		  has.credentials {
			uid
		  }
		  default.namespace {
			uid
		  }
		 ~access.to.namespace {
			uid
		  }
		}
	  }`

	res, err := txn.QueryWithVars(ctx, q, map[string]string{"$userid": request.Uid})
	if err != nil {
		return status.Error(codes.Internal, "Failed to delete Account: "+err.Error())
	}

	var result struct {
		//Get the Device edge details from the query response and build JSON
		Accounts []*Account `json:"account"`
	}

	err = json.Unmarshal(res.Json, &result)
	if err != nil {
		return err
	}

	if len(result.Accounts) != 1 {
		return status.Error(codes.NotFound, "The Account is not found")
	}

	//Append edge if there is a default.namespace edge
	if len(result.Accounts[0].DefaultNamespace) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:   request.Uid,
			Predicate: "default.namespace",
			ObjectId:  result.Accounts[0].DefaultNamespace[0].UID,
		})
	}

	//Append edge if there is a has.credentials edge
	if len(result.Accounts[0].HasCredentials) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:   request.Uid,
			Predicate: "has.credentials",
			ObjectId:  result.Accounts[0].HasCredentials[0].UID,
		})
	}

	//Append edge if there is a reverse edge access.to.namespace edge
	if len(result.Accounts[0].AccessToNamespace) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:   request.Uid,
			Predicate: "access.to.namespace",
			ObjectId:  result.Accounts[0].AccessToNamespace[0].UID,
		})
	}

	//Delete all the edges appended in mutation m
	dgo.DeleteEdges(m, request.Uid, "_STAR_ALL")

	//Append node for the user account
	if len(result.Accounts[0].Node.UID) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:     result.Accounts[0].Node.UID,
			Predicate:   "_STAR_ALL",
			ObjectId:    "_STAR_ALL",
			ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		})
	}

	//Append node related to has.credentials edge
	if len(result.Accounts[0].HasCredentials) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:     result.Accounts[0].HasCredentials[0].UID,
			Predicate:   "_STAR_ALL",
			ObjectId:    "_STAR_ALL",
			ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		})
	}

	//Append node related to default.namespace edge
	if len(result.Accounts[0].DefaultNamespace) == 1 {
		m.Del = append(m.Del, &api.NQuad{
			Subject:     result.Accounts[0].DefaultNamespace[0].UID,
			Predicate:   "_STAR_ALL",
			ObjectId:    "_STAR_ALL",
			ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		})
	}

	_, err = txn.Mutate(context.Background(), m)
	if err != nil {
		return err
	}

	return nil
}

//AssignOwner is a method to delete the Account
func (s *DGraphRepo) AssignOwner(ctx context.Context, ownerID, acccountID string) (err error) {

	txn := s.Dg.NewTxn()
	m := &api.Mutation{CommitNow: true}

	//Added the owns predicate in teh mutation
	m.Set = append(m.Set, &api.NQuad{
		Subject:   ownerID,
		Predicate: "owns",
		ObjectId:  acccountID,
	})

	_, err = txn.Mutate(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

//RemoveOwner is a method to delete the Account.
func (s *DGraphRepo) RemoveOwner(ctx context.Context, ownerID, acccountID string) (err error) {

	txn := s.Dg.NewTxn()
	m := &api.Mutation{CommitNow: true}

	//Added the owns predicate in teh mutation
	m.Del = append(m.Del, &api.NQuad{
		Subject:   ownerID,
		Predicate: "owns",
		ObjectId:  acccountID,
	})

	_, err = txn.Mutate(ctx, m)
	if err != nil {
		return err
	}

	return nil
}
