package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"flag"

	retry "github.com/avast/retry-go"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"

	"os"

	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
)

var (
	dgraphURL = "localhost:9080"
	drop      bool
)

func init() {
	flag.BoolVar(&drop, "drop", false, "Drop all data in dgraph before import")
	envURL := os.Getenv("DGRAPH_URL")
	if envURL != "" {
		dgraphURL = envURL
	}
}

func main() {
	flag.Parse()
	counter := 0

	err := retry.Do(func() error {
		conn, _ := grpc.Dial(dgraphURL, grpc.WithInsecure())
		defer conn.Close()

		dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

		repo := dgraph.NewDGraphRepo(dg)

		counter++
		fmt.Println("----------- Attempt " + strconv.Itoa(counter) + " -----------")

		if drop {
			err := dg.Alter(context.Background(), &api.Operation{DropAll: true})
			if err != nil {
				return err
			}
			fmt.Println("Dgraph Data Drop Successful")
		}

		err := dgraph.ImportSchema(dg, true)
		if err != nil {
			fmt.Println("Import failed with error" + err.Error())
			return err
		}
		fmt.Println("Dgraph Schema Import Successful")

		_, _, err = dgraph.ImportStandardSet(repo)
		if err != nil {
			fmt.Println("Import Standard set failed with error" + err.Error())
			return err
		}
		fmt.Println("Dgraph Test Data Import Successful")
		return nil
	}, retry.Delay(time.Second*2), retry.Attempts(5))

	if err != nil {
		panic(err)
	}
}
