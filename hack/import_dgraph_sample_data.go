package main

import (
	"context"
	"fmt"
	"time"

	"flag"

	retry "github.com/avast/retry-go"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"

	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
)

var (
	dgraphURL string
	drop      bool
)

func init() {
	flag.BoolVar(&drop, "drop", false, "Drop all data in dgraph before import")
	flag.StringVar(&dgraphURL, "host", "localhost:9080", "dgraph host and port")
}

func main() {
	flag.Parse()

	retry.Do(func() error {
		conn, _ := grpc.Dial(dgraphURL, grpc.WithInsecure())
		defer conn.Close()

		dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

		repo := dgraph.NewDGraphRepo(dg)

		if drop {
			err := dg.Alter(context.Background(), &api.Operation{DropAll: true})
			if err != nil {
				return err
			}
			fmt.Println("Dropped data")
		}

		_ = dgraph.ImportSchema(dg)
		fmt.Println("Imported schema")

		_, _, err := dgraph.ImportStandardSet(repo)
		if err != nil {
			return err
		}
		return nil
	}, retry.Delay(time.Second*2), retry.Attempts(40))

}
