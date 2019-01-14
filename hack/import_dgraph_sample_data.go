package main

import (
	"context"
	"encoding/json"
	"fmt"

	"flag"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/infinimesh/infinimesh/pkg/node"
	"google.golang.org/grpc"
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

	conn, _ := grpc.Dial(dgraphURL, grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	if drop {
		err := dg.Alter(context.Background(), &api.Operation{DropAll: true})
		if err != nil {
			panic(err)
		}
		fmt.Println("Dropped data")
	}
	err := dg.Alter(context.Background(), &api.Operation{
		Schema: `
  name: string @index(exact) .
  username: string @index(exact) .
  action: string @index(term) .
  type: string @index(exact) .
  access.to: uid @reverse .
  has.credentials: uid @reverse .
  password: password .`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Imported schema")

	r := &node.Object{
		Node: node.Node{
			UID:  "_:home",
			Type: "object",
		},
		Name: "Johannes' Home",
		Contains: &node.Object{
			Node: node.Node{
				UID:  "_:first-floor",
				Type: "object",
			},
			Name: "First Floor",
			Contains: &node.Object{
				Node: node.Node{
					UID:  "_:living-room",
					Type: "object",
				},
				Name: "Living Room",
				ContainsDevice: &node.Device{
					Node: node.Node{
						UID:  "_:PC",
						Type: "device",
					},
					Name: "PC",
				},
			},
		},
	}

	bytes, _ := json.Marshal(&r)

	a, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(a.GetUids())

	u := &node.Account{
		Node: node.Node{
			UID:  "_:user",
			Type: "account",
		},
		Name: "joe",
		AccessTo: &node.Object{
			Node: node.Node{
				UID: a.GetUids()["home"],
			},
			AccessToInherit:    true,
			AccessToPermission: "WRITE",
		},
		HasCredentials: &node.UsernameCredential{
			Node: node.Node{
				UID:  "_:creds",
				Type: "credentials",
			},
			Username: "joe",
			Password: "test123",
		},
	}

	bytes, _ = json.Marshal(&u)
	a1, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(a1.GetUids())

	// Direct access
	u2 := &node.Account{
		Node: node.Node{
			UID: a1.GetUids()["user"],
		},
		AccessToDevice: &node.Device{
			Node: node.Node{
				UID:  "_:device2",
				Type: "device",
			},
			Name:                     "some device",
			AccessToDevicePermission: "WRITE",
		},
	}

	bytes, _ = json.Marshal(&u2)
	_, err = dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	r2 := &node.Object{
		Node: node.Node{
			UID: a.GetUids()["home"],
		},
		Contains: &node.Object{
			Node: node.Node{
				UID:  "_:second-floor",
				Type: "object",
			},
			Name: "Second Floor",
		},
	}

	bytes, _ = json.Marshal(&r2)
	a2, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(a2.GetUids())

	r3 := &node.Object{
		Node: node.Node{
			UID: a.GetUids()["home"],
		},
		ContainsDevice: &node.Device{
			Node: node.Node{
				UID:  "_:lamp1",
				Type: "device",
			},
			Name: "le lamp",
		},
	}

	bytes, _ = json.Marshal(&r3)
	a3, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(a3.GetUids())

	// Direct access
	u4 := &node.Account{
		Node: node.Node{
			UID: a1.GetUids()["user"],
		},
		AccessTo: &node.Object{
			Node: node.Node{
				UID:  "_:enclosingroom",
				Type: "object",
			},
			Name: "Enclosing Room",
			ContainsDevice: &node.Device{
				Node: node.Node{
					UID:  "_:enclosedobject",
					Type: "device",
				},
				Name: "Enclosing-room-device",
			},
		},
	}

	bytes, _ = json.Marshal(&u4)
	a6, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(a6.GetUids())
}
