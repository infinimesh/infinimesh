package node

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	dgraphURL   = "localhost:9080"
	grpcAddress = ":8082"
)

func TestStuff(t *testing.T) {
	conn, _ := grpc.Dial(dgraphURL, grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	r := &Object{
		Node: Node{
			UID:  "_:home",
			Type: "object",
		},
		Name: "Johannes' Home",
		Contains: &Object{
			Node: Node{
				UID:  "_:first-floor",
				Type: "object",
			},
			Name: "First Floor",
			Contains: &Object{
				Node: Node{
					UID:  "_:living-room",
					Type: "object",
				},
				Name: "Living Room",
				ContainsDevice: &Device{
					Node: Node{
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
	require.NoError(t, err)

	fmt.Println(a.GetUids())

	u := &Account{
		Node: Node{
			UID:  "_:user",
			Type: "user",
		},
		Name: "joe",
		AccessTo: &Object{
			Node: Node{
				UID: a.GetUids()["home"],
			},
			AccessToInherit:    true,
			AccessToPermission: "WRITE",
		},
	}

	bytes, _ = json.Marshal(&u)
	a1, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	require.NoError(t, err)

	fmt.Println(a1.GetUids())

	// Direct access
	u2 := &Account{
		Node: Node{
			UID: a1.GetUids()["user"],
		},
		AccessToDevice: &Device{
			Node: Node{
				UID:  "_:device2",
				Type: "device",
			},
			Name:                     "some device",
			AccessToDevicePermission: "WRITE",
		},
	}

	bytes, _ = json.Marshal(&u2)
	_, err = dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	require.NoError(t, err)

	r2 := &Object{
		Node: Node{
			UID: a.GetUids()["home"],
		},
		Contains: &Object{
			Node: Node{
				UID:  "_:second-floor",
				Type: "object",
			},
			Name: "Second Floor",
		},
	}

	bytes, _ = json.Marshal(&r2)
	a2, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	require.NoError(t, err)

	fmt.Println(a2.GetUids())

	r3 := &Object{
		Node: Node{
			UID: a.GetUids()["home"],
		},
		ContainsDevice: &Device{
			Node: Node{
				UID:  "_:lamp1",
				Type: "device",
			},
			Name: "le lamp",
		},
	}

	bytes, _ = json.Marshal(&r3)
	a3, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	require.NoError(t, err)

	fmt.Println(a3.GetUids())

	// Direct access
	u4 := &Account{
		Node: Node{
			UID: a1.GetUids()["user"],
		},
		AccessTo: &Object{
			Node: Node{
				UID:  "_:enclosingroom",
				Type: "object",
			},
			Name: "Enclosing Room",
			ContainsDevice: &Device{
				Node: Node{
					UID:  "_:enclosedobject",
					Type: "device",
				},
				Name: "Enclosing-room-device",
			},
		},
	}

	bytes, _ = json.Marshal(&u4)
	a6, err := dg.NewTxn().Mutate(context.Background(), &api.Mutation{SetJson: bytes, CommitNow: true})
	require.NoError(t, err)

	fmt.Println(a6.GetUids())

}
