package dgraph

import (
	"context"
	"testing"

	"os"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/node"
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
	err = ImportSchema(dg)
	if err != nil {
		panic(err)
	}
}

func TestAuthorize(t *testing.T) {
	ctx := context.Background()
	account, err := repo.CreateAccount(ctx, randomdata.SillyName(), "password")
	require.NoError(t, err)

	node, err := repo.CreateObject(ctx, "sample-node", "", "asset", "default")
	require.NoError(t, err)

	err = repo.Authorize(ctx, account, node, "READ", true)
	require.NoError(t, err)

	decision, err := repo.IsAuthorized(ctx, node, account, "READ")
	require.NoError(t, err)
	require.True(t, decision)
}
