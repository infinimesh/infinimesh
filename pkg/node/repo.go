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

package node

import (
	"context"

	"github.com/slntopp/infinimesh/pkg/node/nodepb"
)

//KindAsset and KindDevice are constants for creating Object Heirarchy
const (
	KindAsset  = "asset"
	KindDevice = "device"
)

//Repo Interface to expose methods
type Repo interface {
	//Accounts
	CreateUserAccount(ctx context.Context, username, password string, isRoot, isAdmin, enabled bool) (uid string, err error)
	ListAccounts(ctx context.Context) (accounts []*nodepb.Account, err error)
	ListAccountsforAdmin(ctx context.Context, requestorID string) (accounts []*nodepb.Account, err error)
	UpdateAccount(ctx context.Context, account *nodepb.UpdateAccountRequest, isself bool) (err error)
	GetAccount(ctx context.Context, accountID string) (account *nodepb.Account, err error)
	SetPassword(ctx context.Context, account, password string) error
	DeleteAccount(ctx context.Context, account *nodepb.DeleteAccountRequest) (err error)
	AssignOwner(ctx context.Context, ownerID, accountID string) (err error)
	RemoveOwner(ctx context.Context, ownerID, accountID string) (err error)
	UserExists(ctx context.Context, account string) (exists bool, err error)

	//Authorizations
	IsAuthorized(ctx context.Context, target, who, action string) (decision bool, err error)
	IsAuthorizedNamespace(ctx context.Context, namespaceid, account string, action nodepb.Action) (decision bool, err error)
	Authorize(ctx context.Context, account, node, action string, inherit bool) (err error)
	AuthorizeNamespace(ctx context.Context, account, namespaceID string, action nodepb.Action) (err error)
	Authenticate(ctx context.Context, username, password string) (success bool, uid string, defaultNamespace string, err error)

	//Objects
	CreateObject(ctx context.Context, name, parentID, kind, namespaceID string) (id string, err error)
	DeleteObject(ctx context.Context, uid string) (err error)
	ListForAccount(ctx context.Context, account string, namespaceID string, recurse bool) (inheritedObjects []*nodepb.Object, err error)

	//Namespaces
	CreateNamespace(ctx context.Context, name string) (id string, err error)
	GetNamespace(ctx context.Context, uid string) (namespace *nodepb.Namespace, err error)
	GetNamespaceID(ctx context.Context, uid string) (namespace *nodepb.Namespace, err error)
	ListNamespaces(ctx context.Context) (namespaces []*nodepb.Namespace, err error)
	ListNamespacesForAccount(ctx context.Context, accountID string) (namespaces []*nodepb.Namespace, err error)
	ListPermissionsInNamespace(ctx context.Context, namespaceID string) (permissions []*nodepb.Permission, err error)
	DeletePermissionInNamespace(ctx context.Context, namespaceID, accountID string) (err error)
	SoftDeleteNamespace(ctx context.Context, namespaceID string) (err error)
	HardDeleteNamespace(ctx context.Context, datecondition string, retentionperiod string) (err error)
	UpdateNamespace(ctx context.Context, namespace *nodepb.UpdateNamespaceRequest) (err error)
	GetRetentionPeriods(ctx context.Context) (retentionperiod []int, err error)
}
