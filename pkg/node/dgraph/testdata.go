package dgraph

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/tools"
)

func ImportSchema(dg *dgo.Dgraph) error {
	return dg.Alter(context.Background(), &api.Operation{
		Schema: `
  name: string @index(exact) .
  username: string @index(exact) .
  action: string @index(term) .
  type: string @index(exact) .
  access.to: uid @reverse .
  contains: uid @reverse .
  has.credentials: uid @reverse .
  password: password .`,
	})

}

func ImportStandardSet(repo node.Repo) error {
	building, err := repo.CreateObject(context.Background(), "Angerstr 14", "", node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	first, err := repo.CreateObject(context.Background(), "First Floor", building, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Second Floor", building, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	apartment1Right, err := repo.CreateObject(context.Background(), "Apartment right side", first, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Entrance", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bathroom", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Kitchen", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bedroom", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Kinderzimmer", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Walk-through room", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Living room", apartment1Right, node.KindAsset, "joe/default")
	if err != nil {
		return err
	}

	user, err := repo.CreateAccount(context.Background(), "joe", "test123")
	if err != nil {
		return err
	}

	fmt.Println("User: ", user)

	result := repo.Authorize(context.Background(), user, apartment1Right, "WRITE", true)

	nested, err := repo.ListForAccount(context.Background(), user)
	if err != nil {
		return nil
	}

	tools.PrettyPrint(nested)

	return result
}
