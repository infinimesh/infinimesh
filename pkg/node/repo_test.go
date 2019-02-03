package node

import (
	"context"
	"testing"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var repo Repo

func init() {
	conn, err := grpc.Dial("server:9080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	repo = &dGraphRepo{dg: dg}
	err = ImportSchema(dg)
	if err != nil {
		panic(err)
	}
}

func TestAuthorize(t *testing.T) {
	ctx := context.Background()
	account, err := repo.CreateAccount(ctx, randomdata.SillyName(), "password")
	require.NoError(t, err)

	node, err := repo.CreateObject(ctx, "sample-node", "")
	require.NoError(t, err)

	err = repo.Authorize(ctx, account, node, "READ", true)
	require.NoError(t, err)

	decision, err := repo.IsAuthorized(ctx, node, account, "READ")
	require.NoError(t, err)
	require.True(t, decision)
}
