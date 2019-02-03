package node

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/tools"
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

	err = repo.DeleteObject(context.Background(), "0xc376")

	require.NoError(t, err)

}

func TestList(t *testing.T) {
	t.SkipNow()
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	repo := &dGraphRepo{dg: dg}

	_, err = repo.CreateObject(context.Background(), "Johannes' Home", "")
	require.NoError(t, err)

	_, _, in, err := repo.ListForAccount(context.Background(), "0x13886")
	require.NoError(t, err)

	_ = in
	tools.PrettyPrint(in)
}

func TestStuff(t *testing.T) {
	t.SkipNow()
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	repo := &dGraphRepo{dg: dg}

	err = ImportStandardSet(repo)
	require.NoError(t, err)

}
