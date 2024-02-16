package graph_test

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	credentials_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/credentials"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	sessions_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/sessions"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/accounts"
	"github.com/stretchr/testify/assert"
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
		ctx context.Context
	}
}

func newAccountsControllerFixture(t *testing.T) accountsControllerFixture {
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

	f.data.ctx = context.TODO()

	f.mocks.ica_repo.EXPECT().GetVertexCol(f.data.ctx, schema.PERMISSIONS_GRAPH.Name, schema.ACCOUNTS_COL).
		Return(f.mocks.col)

	f.mocks.ica_repo.EXPECT().GetVertexCol(f.data.ctx, schema.CREDENTIALS_GRAPH.Name, schema.CREDENTIALS_COL).
		Return(f.mocks.cred_col)

	f.mocks.ica_repo.EXPECT().GetEdgeCol(f.data.ctx, schema.ACC2NS).
		Return(f.mocks.acc2ns)
	f.mocks.ica_repo.EXPECT().GetEdgeCol(f.data.ctx, schema.NS2ACC).
		Return(f.mocks.ns2acc)

	f.repo = graph.NewAccountsController(
		zap.NewExample(),
		f.mocks.db, f.mocks.rdb, f.mocks.sessions,
		f.mocks.ica_repo, f.mocks.repo,
		f.mocks.cred,
	)

	return f
}

// Token
//

func TestToken_FailsOn_WrongCredentials(t *testing.T) {
	f := newAccountsControllerFixture(t)

	f.mocks.cred.EXPECT().Find(context.TODO(), "standard", "username", "password").
		Return(nil, fmt.Errorf("not found"))

	_, err := f.repo.Token(context.TODO(), &connect.Request[node.TokenRequest]{
		Msg: &node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "standard", Data: []string{"username", "password"},
			},
		},
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = Wrong credentials given")
}
