package graph_test

import (
	"testing"

	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	redis_mocks "github.com/infinimesh/infinimesh/mocks/github.com/go-redis/redis/v8"
	credentials_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/credentials"
	graph_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/graph"
	sessions_mocks "github.com/infinimesh/infinimesh/mocks/github.com/infinimesh/infinimesh/pkg/sessions"
	"github.com/infinimesh/infinimesh/pkg/graph"
	"github.com/infinimesh/proto/node/accounts"
	"go.uber.org/zap"
)

type accountsControllerFixture struct {
	repo *graph.AccountsController

	mocks struct {
		db *driver_mocks.MockDatabase

		col    *driver_mocks.MockCollection
		acc2ns *driver_mocks.MockCollection
		ns2acc *driver_mocks.MockCollection

		rdb *redis_mocks.MockCmdable

		sessions *sessions_mocks.MockSessionsHandler
		cred     *credentials_mocks.MockCredentialsController
		ica_repo *graph_mocks.MockInfinimeshCommonActionsRepo
		repo     *graph_mocks.MockInfinimeshGenericActionsRepo[*accounts.Account]
	}
}

func newAccountsControllerFixture(t *testing.T) {
	f := &accountsControllerFixture{}

	f.mocks.db = &driver_mocks.MockDatabase{}
	f.mocks.col = &driver_mocks.MockCollection{}
	f.mocks.acc2ns = &driver_mocks.MockCollection{}
	f.mocks.ns2acc = &driver_mocks.MockCollection{}
	f.mocks.rdb = redis_mocks.NewMockCmdable(t)
	f.mocks.sessions = &sessions_mocks.MockSessionsHandler{}
	f.mocks.cred = &credentials_mocks.MockCredentialsController{}
	f.mocks.ica_repo = &graph_mocks.MockInfinimeshCommonActionsRepo{}
	f.mocks.repo = &graph_mocks.MockInfinimeshGenericActionsRepo[*accounts.Account]{}

	f.repo = graph.NewAccountsController(
		zap.NewExample(),
		f.mocks.db, f.mocks.rdb, f.mocks.sessions,
		f.mocks.ica_repo, f.mocks.repo,
		f.mocks.cred,
	)
}

// Token
//
