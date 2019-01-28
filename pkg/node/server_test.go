package node

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var dgraphURL = "localhost:9080"

func TestDelete(t *testing.T) {
	t.SkipNow()
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	repo := &dGraphRepo{dg: dg}

	err = repo.DeleteObject(context.Background(), "0x4f7c")

	require.NoError(t, err)

}
