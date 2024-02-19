package graph_test

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	credentials_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/credentials"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	sessions_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/sessions"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	"github.com/infinimesh/proto/node/accounts"
	"github.com/infinimesh/proto/node/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type accountsControllerFixture struct {
	repo *graph.AccountsController

	mocks struct {
		db *driver_mocks.MockDatabase

		col      *driver_mocks.MockCollection
		acc2ns   *driver_mocks.MockCollection
		ns2acc   *driver_mocks.MockCollection
		cred_col *driver_mocks.MockCollection

		rdb *redis_mocks.MockCmdable

		sessions *sessions_mocks.MockSessionsHandler
		cred     *credentials_mocks.MockCredentialsController
		ica_repo *graph_mocks.MockInfinimeshCommonActionsRepo
		repo     *graph_mocks.MockInfinimeshGenericActionsRepo[*accounts.Account]
	}
	data struct {
		ctx              context.Context
		ctx_no_requestor context.Context

		auth_data []string
		account   graph.Account
		session   *sessions.Session

		create_request *accounts.CreateRequest
	}
}

func newAccountsControllerFixture(t *testing.T) accountsControllerFixture {
	t.Parallel()

	f := accountsControllerFixture{}

	f.mocks.db = &driver_mocks.MockDatabase{}
	f.mocks.col = &driver_mocks.MockCollection{}
	f.mocks.acc2ns = &driver_mocks.MockCollection{}
	f.mocks.ns2acc = &driver_mocks.MockCollection{}
	f.mocks.cred_col = &driver_mocks.MockCollection{}
	f.mocks.rdb = redis_mocks.NewMockCmdable(t)
	f.mocks.sessions = &sessions_mocks.MockSessionsHandler{}
	f.mocks.cred = &credentials_mocks.MockCredentialsController{}
	f.mocks.ica_repo = &graph_mocks.MockInfinimeshCommonActionsRepo{}
	f.mocks.repo = &graph_mocks.MockInfinimeshGenericActionsRepo[*accounts.Account]{}

	f.data.ctx = context.WithValue(context.TODO(), inf.InfinimeshAccountCtxKey, uuid.New().String())
	f.data.ctx_no_requestor = context.TODO()

	f.mocks.ica_repo.EXPECT().GetVertexCol(f.data.ctx_no_requestor, schema.PERMISSIONS_GRAPH.Name, schema.ACCOUNTS_COL).
		Return(f.mocks.col)

	f.mocks.ica_repo.EXPECT().GetVertexCol(f.data.ctx_no_requestor, schema.CREDENTIALS_GRAPH.Name, schema.CREDENTIALS_COL).
		Return(f.mocks.cred_col)

	f.mocks.ica_repo.EXPECT().GetEdgeCol(f.data.ctx_no_requestor, schema.ACC2NS).
		Return(f.mocks.acc2ns)
	f.mocks.ica_repo.EXPECT().GetEdgeCol(f.data.ctx_no_requestor, schema.NS2ACC).
		Return(f.mocks.ns2acc)

	f.repo = graph.NewAccountsController(
		zap.NewExample(),
		f.mocks.db, f.mocks.rdb, f.mocks.sessions,
		f.mocks.ica_repo, f.mocks.repo,
		f.mocks.cred,
	)

	f.data.auth_data = []string{"username", "password"}
	f.data.account = graph.Account{
		Account: &accounts.Account{
			Uuid:             uuid.New().String(),
			Enabled:          true,
			DefaultNamespace: "infinimesh",
		},
	}
	f.data.account.DocumentMeta = driver.DocumentMeta{
		Key: f.data.account.Uuid,
		ID:  driver.NewDocumentID(schema.ACCOUNTS_COL, f.data.account.Uuid),
	}
	f.data.session = &sessions.Session{
		Id:     "session_id",
		Client: "",
	}

	f.data.create_request = &accounts.CreateRequest{
		Account: f.data.account.Account,
		Credentials: &accounts.Credentials{
			Type: "standard",
			Data: f.data.auth_data,
		},
	}

	return f
}

// Token
//

func TestToken_FailsOn_WrongCredentials(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Authorize(context.TODO(), "standard", f.data.auth_data[0], f.data.auth_data[1]).
		Return(nil, false)

	_, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Wrong credentials given")
}

func TestToken_FailsOn_AccountDisabled(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Authorize(context.TODO(), "standard", f.data.auth_data[0], f.data.auth_data[1]).
		Return(&accounts.Account{Enabled: false}, true)

	_, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = Account is disabled")
}

func TestToken_FailsOn_Session(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Authorize(context.TODO(), "standard", f.data.auth_data[0], f.data.auth_data[1]).
		Return(f.data.account.Account, true)

	f.mocks.sessions.EXPECT().New(int64(0), "").Return(f.data.session)
	f.mocks.sessions.EXPECT().Store(f.data.account.Uuid, f.data.session).
		Return(assert.AnError)

	_, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = Failed to issue token: session")
}

func TestToken_User_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Authorize(context.TODO(), "standard", f.data.auth_data[0], f.data.auth_data[1]).
		Return(f.data.account.Account, true)

	f.mocks.sessions.EXPECT().New(int64(0), "").Return(f.data.session)
	f.mocks.sessions.EXPECT().Store(f.data.account.Uuid, f.data.session).
		Return(nil)

	res, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
		},
	})

	assert.NoError(t, err)

	token, _, err := new(jwt.Parser).ParseUnverified(res.Msg.GetToken(), jwt.MapClaims{})
	assert.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, f.data.account.Uuid, claims[inf.INFINIMESH_ACCOUNT_CLAIM])
	assert.Equal(t, f.data.session.Id, claims[inf.INFINIMESH_SESSION_CLAIM])
}

func TestToken_Root_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Authorize(context.TODO(), "standard", f.data.auth_data[0], f.data.auth_data[1]).
		Return(f.data.account.Account, true)

	f.mocks.sessions.EXPECT().New(int64(0), "").Return(f.data.session)
	f.mocks.sessions.EXPECT().Store(f.data.account.Uuid, f.data.session).
		Return(nil)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx_no_requestor, mock.Anything, mock.Anything).
		Return(true, access.Level_ROOT)

	root := true
	res, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
			Inf: &root,
		},
	})

	assert.NoError(t, err)

	token, _, err := new(jwt.Parser).ParseUnverified(res.Msg.GetToken(), jwt.MapClaims{})
	assert.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, f.data.account.Uuid, claims[inf.INFINIMESH_ACCOUNT_CLAIM])
	assert.Equal(t, f.data.session.Id, claims[inf.INFINIMESH_SESSION_CLAIM])
	assert.True(t, claims[inf.INFINIMESH_ROOT_CLAIM].(bool))
}

func TestToken_LoginAs_FailsOn_RecursiveToken(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.data.ctx = context.WithValue(f.data.ctx, inf.InfinimeshAccountCtxKey, f.data.account.Uuid)
	_, err := f.repo.Token(f.data.ctx, &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
			Uuid: &f.data.account.Uuid,
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = You can't create such token for yourself")
}

func TestToken_LoginAs_FailsOn_AccessLevelAndGet(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.Anything,
	).Return(assert.AnError)

	_, err := f.repo.Token(f.data.ctx, &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
			Uuid: &f.data.account.Uuid,
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Account not found")
}

func TestToken_LoginAs_FailsOn_NotEnoughAccess(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Access = &access.Access{
				Level: access.Level_NONE,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	_, err := f.repo.Token(f.data.ctx, &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: f.data.auth_data,
			},
			Uuid: &f.data.account.Uuid,
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Wrong credentials given")
}

// Get
//

func TestGet_FailsOn_AccessLevelAndGet(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.Anything,
	).Return(assert.AnError)

	_, err := f.repo.Get(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: &accounts.Account{
			Uuid: f.data.account.Uuid,
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Account not found or not enough Access Rights")
}

func TestGet_FailsOn_NotEnoughAccess(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Access = &access.Access{
				Level: access.Level_NONE,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	_, err := f.repo.Get(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: &accounts.Account{
			Uuid: f.data.account.Uuid,
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = Not enough Access Rights")
}

func TestGet_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Account = f.data.account.Account
			acc.Access = &access.Access{
				Level: access.Level_READ,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	res, err := f.repo.Get(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: &accounts.Account{
			Uuid: "me",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, f.data.account.Uuid, res.Msg.GetUuid())
	assert.Equal(t, f.data.account.Title, res.Msg.GetTitle())
}

// List
//

func TestListAccounts_FailsOn_ListQuery(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.repo.EXPECT().ListQuery(f.data.ctx, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	_, err := f.repo.List(f.data.ctx, &connect.Request[node.EmptyMessage]{
		Msg: &node.EmptyMessage{},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = Failed to list accounts")
}

func TestListAccounts_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	accs := []*accounts.Account{
		{Uuid: "1", Title: "1"},
		{Uuid: "2", Title: "2"},
	}

	f.mocks.repo.EXPECT().ListQuery(f.data.ctx, mock.Anything, mock.Anything).
		Return(&graph.ListQueryResult[*accounts.Account]{
			Result: accs,
			Count:  2,
		}, nil)

	res, err := f.repo.List(f.data.ctx, &connect.Request[node.EmptyMessage]{
		Msg: &node.EmptyMessage{},
	})

	assert.NoError(t, err)
	assert.Len(t, res.Msg.GetAccounts(), 2)

	for i, acc := range res.Msg.GetAccounts() {
		assert.Equal(t, accs[i].Uuid, acc.GetUuid())
		assert.Equal(t, accs[i].Title, acc.GetTitle())
	}
}

// Create
//

func TestAccountCreate_FailsOn_AccessLevel(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(false, access.Level_NONE)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "permission_denied: no Access to Namespace infinimesh")
}

func TestAccountCreate_FailsOn_CreateDocument(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, assert.AnError)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error while creating Account")
}

func TestAccountCreate_FailsOn_Link(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(f.data.account.DocumentMeta, nil)

	f.mocks.ica_repo.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
		return ns.Key == "infinimesh"
	}), mock.MatchedBy(func(acc *graph.Account) bool {
		return acc.Uuid == f.data.account.Uuid
	}), access.Level_ADMIN, access.Role_OWNER).
		Return(assert.AnError)

	// Deferred cleanup
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.account.DocumentMeta.Key).
		Return(driver.DocumentMeta{}, nil)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error while creating Account")
}

func TestAccountCreate_FailsOn_MakeCredentials(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(f.data.account.DocumentMeta, nil)

	f.mocks.ica_repo.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
		return ns.Key == "infinimesh"
	}), mock.MatchedBy(func(acc *graph.Account) bool {
		return acc.Uuid == f.data.account.Uuid
	}), access.Level_ADMIN, access.Role_OWNER).
		Return(nil)

	f.mocks.cred.EXPECT().MakeCredentials(f.data.create_request.Credentials).
		Return(nil, assert.AnError)

	// Deferred cleanup
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.account.DocumentMeta.Key).
		Return(driver.DocumentMeta{}, nil)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error while creating Account: assert.AnError general error for testing")
}

func TestAccountCreate_FailsOn_SetCredentials(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(f.data.account.DocumentMeta, nil)

	f.mocks.ica_repo.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
		return ns.Key == "infinimesh"
	}), mock.MatchedBy(func(acc *graph.Account) bool {
		return acc.Uuid == f.data.account.Uuid
	}), access.Level_ADMIN, access.Role_OWNER).
		Return(nil)

	f.mocks.cred.EXPECT().MakeCredentials(f.data.create_request.Credentials).
		Return(&credentials.StandardCredentials{}, nil)

	f.mocks.cred.EXPECT().SetCredentials(f.data.ctx, f.data.account.DocumentMeta.ID, &credentials.StandardCredentials{}).
		Return(assert.AnError)

	// Deferred cleanup
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, f.data.account.DocumentMeta.Key).
		Return(driver.DocumentMeta{}, nil)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Error while creating Account: assert.AnError general error for testing")
}

func TestAccountCreate_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevel(f.data.ctx, mock.Anything, mock.Anything).
		Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(f.data.account.DocumentMeta, nil)

	f.mocks.ica_repo.EXPECT().Link(f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
		return ns.Key == "infinimesh"
	}), mock.MatchedBy(func(acc *graph.Account) bool {
		return acc.Uuid == f.data.account.Uuid
	}), access.Level_ADMIN, access.Role_OWNER).
		Return(nil)

	f.mocks.cred.EXPECT().MakeCredentials(f.data.create_request.Credentials).
		Return(&credentials.StandardCredentials{}, nil)

	f.mocks.cred.EXPECT().SetCredentials(f.data.ctx, f.data.account.DocumentMeta.ID, &credentials.StandardCredentials{}).
		Return(nil)

	_, err := f.repo.Create(f.data.ctx, &connect.Request[accounts.CreateRequest]{
		Msg: f.data.create_request,
	})

	assert.NoError(t, err)
}

// Update
//

func TestAccountUpdate_FailsOn_AccessLevelAndGet(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.Anything,
	).Return(assert.AnError)

	_, err := f.repo.Update(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("rpc error: code = PermissionDenied desc = No Access to Account %s", f.data.account.Uuid))
}

func TestAccountUpdate_FailsOn_NewNS_AccessLevel(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			buf := *f.data.account.Account
			acc.Account = &buf
			acc.Account.DefaultNamespace = "different_ns"
			acc.Access = &access.Access{
				Level: access.Level_ADMIN,
			}
			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	f.mocks.ica_repo.EXPECT().AccessLevel(
		f.data.ctx, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
			return ns.Key == f.data.account.DefaultNamespace
		}),
	).Return(false, access.Level_NONE)

	_, err := f.repo.Update(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = Account has no Access to Namespace infinimesh")
}

func TestAccountUpdate_FailsOn_UpdateDocument(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			buf := *f.data.account.Account
			acc.Account = &buf
			acc.Account.DefaultNamespace = "different_ns"
			acc.Access = &access.Access{
				Level: access.Level_ADMIN,
			}
			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	f.mocks.ica_repo.EXPECT().AccessLevel(
		f.data.ctx, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
			return ns.Key == f.data.account.DefaultNamespace
		}),
	).Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().UpdateDocument(f.data.ctx, mock.Anything, mock.Anything).
		Return(driver.DocumentMeta{}, assert.AnError)

	_, err := f.repo.Update(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = Error while updating Account")
}

func TestAccountUpdate_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			buf := *f.data.account.Account
			acc.Account = &buf
			acc.Account.DefaultNamespace = "different_ns"
			acc.Access = &access.Access{
				Level: access.Level_ADMIN,
			}
			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	f.mocks.ica_repo.EXPECT().AccessLevel(
		f.data.ctx, mock.Anything, mock.MatchedBy(func(ns *graph.Namespace) bool {
			return ns.Key == f.data.account.DefaultNamespace
		}),
	).Return(true, access.Level_ADMIN)

	f.mocks.col.EXPECT().UpdateDocument(f.data.ctx, mock.Anything, mock.Anything).
		Return(driver.DocumentMeta{}, nil)

	_, err := f.repo.Update(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.NoError(t, err)
}

// Toggle
//

func TestAccountToggle_FailsOn_Get(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.Anything,
	).Return(assert.AnError)

	_, err := f.repo.Toggle(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = Account not found or not enough Access Rights")
}

func TestAccountToggle_FailsOn_NotEnoughAccess(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Account = f.data.account.Account
			acc.Access = &access.Access{
				Level: access.Level_READ,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	_, err := f.repo.Toggle(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("rpc error: code = PermissionDenied desc = No Access to Account %s", f.data.account.Uuid))
}

func TestAccountToggle_FailsOn_RepoToggle(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Account = f.data.account.Account
			acc.Access = &access.Access{
				Level: access.Level_MGMT,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	f.mocks.ica_repo.EXPECT().Toggle(
		f.data.ctx, mock.Anything, mock.Anything,
	).Return(assert.AnError)

	_, err := f.repo.Toggle(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Internal desc = Error while updating Account")
}

func TestAccountToggle_Success(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.ica_repo.EXPECT().AccessLevelAndGet(
		f.data.ctx, mock.Anything, mock.Anything, mock.MatchedBy(func(acc *graph.Account) bool {
			acc.Account = f.data.account.Account
			acc.Access = &access.Access{
				Level: access.Level_MGMT,
			}

			return acc.Uuid == f.data.account.Uuid
		}),
	).Return(nil)

	f.mocks.ica_repo.EXPECT().Toggle(
		f.data.ctx, mock.Anything, mock.Anything,
	).Return(nil)

	_, err := f.repo.Toggle(f.data.ctx, &connect.Request[accounts.Account]{
		Msg: f.data.account.Account,
	})

	assert.NoError(t, err)
}
