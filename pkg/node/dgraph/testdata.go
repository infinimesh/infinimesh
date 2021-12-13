//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package dgraph

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/slntopp/infinimesh/pkg/node"
	"github.com/slntopp/infinimesh/pkg/node/nodepb"
)

//ImportSchema is a method to import the schema in Dgraph DB for tests
func ImportSchema(dg *dgo.Dgraph, drop bool) error {
	if drop {
		err := dg.Alter(context.Background(), &api.Operation{DropAll: drop})
		if err != nil {
			return err
		}
	}
	schema := `
  tags: [string] .
  name: string @index(exact) .
  username: string @index(exact) .
  enabled: bool @index(bool) .
  isRoot: bool @index(bool) .
  isAdmin: bool @index(bool) .
  markfordeletion: bool @index(bool) .
  deleteinitiationtime: datetime @index(day) .
  retentionperiod: int @index(int) .
  action: string @index(term) .
  type: string @index(exact) .
  access.to: uid @reverse .
  children: uid @reverse .
  owns: uid @reverse .
  kind: string @index(exact) .
  has.credentials: uid @reverse .
  access.to.namespace: uid @reverse .
  fingerprint: string @index(exact) .
  certificates: uid @reverse .
  password: password .`
	fmt.Println("Apply Dgraph schema", schema)
	return dg.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})

}

//ImportStandardSet is a method to impor the test data for tests
func ImportStandardSet(repo node.Repo) (userID string, adminID string, err error) {
	// careful,  currently when referencing a namespace, the name of it has to be used, not the id (0x...)
	sharedNs := "shared-project"
	namespace, err := repo.CreateNamespace(context.Background(), sharedNs)
	if err != nil {
		fmt.Println("Create Namespace failed", err)
		return "", "", err
	}
	fmt.Println("Namespace: ", namespace)

	//Create user Joe
	joe, err := repo.CreateUserAccount(context.Background(), "joe", "test123", false, false, true)
	if err != nil {
		fmt.Println("Create Account failed for joe", err)
		return "", "", err
	}
	fmt.Println("User joe: ", joe)

	//get namespace for joe
	ns, err := repo.GetNamespace(context.Background(), "joe")
	if err != nil {
		fmt.Println("Get Namespace failed", err)
		return "", "", err
	}

	//Create user hanswurst
	hanswurst, err := repo.CreateUserAccount(context.Background(), "hanswurst", "hanswurst", false, false, true)
	if err != nil {
		fmt.Println("Create Account failed for hans", err)
		return "", "", err
	}
	fmt.Println("User hanswurst: ", hanswurst)

	//Create root and admin user
	admin, err := repo.CreateUserAccount(context.Background(), "admin", "admin123", true, true, true)
	if err != nil {
		fmt.Println("Create Account failed", err)
		return "", "", err
	}
	fmt.Println("Admin: ", admin)

	// Authorize both users on a shared project
	{
		err = repo.AuthorizeNamespace(context.Background(), joe, namespace, nodepb.Action_WRITE)
		if err != nil {
			fmt.Println("Authorize Namespace failed for joe", err)
			return "", "", err
		}

		err = repo.AuthorizeNamespace(context.Background(), hanswurst, namespace, nodepb.Action_WRITE)
		if err != nil {
			fmt.Println("Authorize Namespace failed for hans", err)
			return "", "", err
		}
	}

	building, err := repo.CreateObject(context.Background(), "Angerstr 14", "", node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	first, err := repo.CreateObject(context.Background(), "First Floor", building, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Second Floor", building, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	apartment1Right, err := repo.CreateObject(context.Background(), "Apartment right side", first, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Entrance", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Bathroom", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Kitchen", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Bedroom", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Kinderzimmer", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Walk-through room", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	livingRoom, err := repo.CreateObject(context.Background(), "Living room", apartment1Right, node.KindAsset, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Test-device", livingRoom, node.KindDevice, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	_, err = repo.CreateObject(context.Background(), "Test-device-no-parent", "", node.KindDevice, ns.Id)
	if err != nil {
		fmt.Println("Create Object failed", err)
		return "", "", err
	}

	return joe, admin, err
}
