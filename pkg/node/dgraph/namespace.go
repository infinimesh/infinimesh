package dgraph

import (
	"context"
	"encoding/json"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func (s *DGraphRepo) ListNamespaces(ctx context.Context) (namespaces []*nodepb.Namespace, err error) {
	const q = `{
                     namespaces(func: eq(type, "namespace")) {
  	               uid
                       name
	             }
                   }`

	res, err := s.Dg.NewReadOnlyTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var resultSet struct {
		Namespaces []*Namespace `json:"namespaces"`
	}

	if err := json.Unmarshal(res.Json, &resultSet); err != nil {
		return nil, err
	}

	for _, namespace := range resultSet.Namespaces {
		namespaces = append(namespaces, &nodepb.Namespace{
			Id:   namespace.UID,
			Name: namespace.Name,
		})
	}

	return namespaces, nil
}

func (s *DGraphRepo) ListPermissionsInNamespace(ctx context.Context, namespace string) (permissions []*nodepb.Permission, err error) {
	const q = `{
  accounts(func: eq(name, "joe")) @filter(eq(type, "namespace")) @normalize @cascade  {
    ~access.to.namespace {
      uid: uid
      name: name
    } @facets(permission)
  }
}`

	res, err := s.Dg.NewReadOnlyTxn().QueryWithVars(ctx, q, map[string]string{"$namespace": namespace})
	if err != nil {
		return nil, err
	}

	var resultSet struct {
		Accounts []*struct {
			UID    string `json:"uid"`
			Name   string `json:"name"`
			Action string `json:"~access.to.namespace|permission"`
		} `json:"accounts"`
	}

	if err := json.Unmarshal(res.Json, &resultSet); err != nil {
		return nil, err
	}

	for _, account := range resultSet.Accounts {
		permissions = append(permissions, &nodepb.Permission{
			Namespace:   namespace,
			AccountId:   account.UID,
			AccountName: account.Name,
			Action:      nodepb.Action(nodepb.Action_value[account.Action]),
		})
	}

	return permissions, nil

}

func (s *DGraphRepo) ListNamespacesForAccount(ctx context.Context, accountID string) (namespaces []*nodepb.Namespace, err error) {
	const q = `query listNamespaces($account: string) {
                     namespaces(func: uid($account)) @normalize @cascade  {
                       access.to.namespace @filter(eq(type, "namespace")) {
                         uid : uid
                         name : name
                       }
	             }
                   }`

	res, err := s.Dg.NewReadOnlyTxn().QueryWithVars(ctx, q, map[string]string{"$account": accountID})
	if err != nil {
		return nil, err
	}

	var resultSet struct {
		Namespaces []*Namespace `json:"namespaces"`
	}

	if err := json.Unmarshal(res.Json, &resultSet); err != nil {
		return nil, err
	}

	for _, namespace := range resultSet.Namespaces {
		namespaces = append(namespaces, &nodepb.Namespace{
			Id:   namespace.UID,
			Name: namespace.Name,
		})
	}

	return namespaces, nil
}

func (s *DGraphRepo) IsAuthorizedNamespace(ctx context.Context, namespace, account string, action nodepb.Action) (decision bool, err error) {
	acc, err := s.GetAccount(ctx, account)
	if err != nil {
		return false, err
	}

	if acc.IsRoot {
		return true, nil
	}

	params := map[string]string{
		"$namespace": namespace,
		"$user_id":   account,
	}

	txn := s.Dg.NewReadOnlyTxn()

	const q = `query access($namespace: string, $user_id: string){
  access(func: uid($user_id)) @cascade {
    name
    uid
    access.to.namespace @filter(eq(name, "$namespace")) {
      uid
      name
      type
    }
  }
}
`

	res, err := txn.QueryWithVars(ctx, q, params)
	if err != nil {
		return false, err
	}
	var access struct {
		Access []Object `json:"access"`
	}

	err = json.Unmarshal(res.Json, &access)
	if err != nil {
		return false, err
	}

	return len(access.Access) > 0, nil
}
