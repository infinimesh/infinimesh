package node

import (
	"context"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

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

func ImportStandardSet(repo Repo) error {
	building, err := repo.CreateObject(context.Background(), "Angerstr 14", "")
	if err != nil {
		return err
	}

	first, err := repo.CreateObject(context.Background(), "First Floor", building)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Second Floor", building)
	if err != nil {
		return err
	}

	apartment1Right, err := repo.CreateObject(context.Background(), "Apartment right side", first)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Entrance", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bathroom", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Kitchen", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bedroom", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Nursery", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Walk-through room", apartment1Right)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Living room", apartment1Right)
	if err != nil {
		return err
	}

	user, err := repo.CreateAccount(context.Background(), "joex", "test123")
	if err != nil {
		return err
	}

	result := repo.Authorize(context.Background(), user, apartment1Right, "WRITE", true)

	_, _, nested, err := repo.ListForAccount(context.Background(), user)
	if err != nil {
		return nil
	}

	tools.PrettyPrint(nested)

	return result
}
