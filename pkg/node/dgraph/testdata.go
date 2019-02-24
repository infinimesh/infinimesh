package dgraph

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func ImportSchema(dg *dgo.Dgraph) error {
	return dg.Alter(context.Background(), &api.Operation{
		Schema: `
  name: string @index(exact) .
  username: string @index(exact) .
  action: string @index(term) .
  type: string @index(exact) .
  access.to: uid @reverse .
  children: uid @reverse .
  owns: uid @reverse .
  has.credentials: uid @reverse .
  password: password .`,
	})

}

func ImportStandardSet(repo node.Repo) error {
	ns := "joe/default"
	_, err := repo.CreateNamespace(context.Background(), ns)
	if err != nil {
		return err
	}

	building, err := repo.CreateObject(context.Background(), "Angerstr 14", "", node.KindAsset, ns)
	if err != nil {
		return err
	}

	first, err := repo.CreateObject(context.Background(), "First Floor", building, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Second Floor", building, node.KindAsset, ns)
	if err != nil {
		return err
	}

	apartment1Right, err := repo.CreateObject(context.Background(), "Apartment right side", first, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Entrance", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bathroom", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Kitchen", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Bedroom", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Kinderzimmer", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Walk-through room", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	livingRoom, err := repo.CreateObject(context.Background(), "Living room", apartment1Right, node.KindAsset, ns)
	if err != nil {
		return err
	}

	_, err = repo.CreateObject(context.Background(), "Test-device", livingRoom, node.KindDevice, ns)
	if err != nil {
		return err
	}

	user, err := repo.CreateUserAccount(context.Background(), "joe", "test123", false)
	if err != nil {
		return err
	}

	fmt.Println("User: ", user)

	result := repo.Authorize(context.Background(), user, apartment1Right, "WRITE", true)
	err = repo.AuthorizeNamespace(context.Background(), user, ns, nodepb.Action_WRITE)
	if err != nil {
		return err
	}

	_, err = repo.ListForAccount(context.Background(), user)
	if err != nil {
		return err
	}

	admin, err := repo.CreateUserAccount(context.Background(), "admin", "admin123", true)
	if err != nil {
		return err
	}
	fmt.Println("Admin: ", admin)

	return result
}
