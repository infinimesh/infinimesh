package dgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func isPermissionSufficient(required, actual string) bool {
	switch required {
	case "WRITE":
		return actual == "WRITE"
	case "READ":
		return actual == "WRITE" || actual == "READ"
	default:
		return false
	}
}

type DGraphRepo struct {
	Dg *dgo.Dgraph
}

func NewDGraphRepo(dg *dgo.Dgraph) node.Repo {
	return &DGraphRepo{Dg: dg}
}

func checkType(ctx context.Context, txn *dgo.Txn, uid, _type string) bool {
	q := `query object($_uid: string, $type: string) {
                object(func: uid($_uid)) @filter(eq(type, $type)) {
                  uid
                }
              }
             `
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$_uid": uid,
		"$type": _type,
	})
	if err != nil {
		return false
	}

	var result struct {
		Object []map[string]interface{} `json:"object"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false
	}

	return len(result.Object) > 0
}

func NameExists(ctx context.Context, txn *dgo.Txn, name, namespace, parent string) bool { //nolint
	var q string
	if parent == "" {
		q = `query object($name: string, $namespace: string, $parent: uid){
  object(func: eq(name, $name)) @cascade {
    uid
    name
    ~owns @filter(eq(name, $namespace)) {
      name
    }
  }
}
`
	} else {
		q = `query exists($name: string, $namespace: string, $parent: uid){
  exists(func: eq(name, $name)) @cascade {
    uid
    name
    ~owns @filter(eq(name, $namespace)) {
      name
    }
    ~children @filter(uid($parent)) {
      uid
      name
    }
  }
}
`

	}

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$parent":    parent,
		"$name":      name,
		"$namespace": namespace,
	})
	if err != nil {
		return false
	}

	var result struct {
		Object []map[string]interface{} `json:"object"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false
	}

	return len(result.Object) > 0
}

func CheckExists(ctx context.Context, txn *dgo.Txn, uid string) bool { //nolint
	q := `query object($_uid: string) {
                object(func: uid($_uid)) {
                  uid
                }
              }
             `
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$_uid": uid,
	})
	if err != nil {
		return false
	}

	var result struct {
		Object []map[string]interface{} `json:"object"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false
	}

	return len(result.Object) > 0
}

func (s *DGraphRepo) AuthorizeNamespace(ctx context.Context, account, namespace string, action nodepb.Action) (err error) {
	txn := s.Dg.NewTxn()

	if ok := checkType(ctx, txn, account, "account"); !ok {
		return errors.New("invalid account")
	}

	// TODO use internal method that runs within txn
	ns, err := s.GetNamespace(ctx, namespace)
	if err != nil {
		return err
	}

	in := Account{
		Node: Node{
			UID: account,
		},
	}

	in.AccessToNamespace = []*Namespace{
		&Namespace{
			Node: Node{
				UID: ns.GetId(),
			},
			AccessToPermission: action.String(),
		},
	}

	js, err := json.Marshal(&in)
	if err != nil {
		return err
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return errors.New("Failed to mutate")
	}
	return nil

}

func (s *DGraphRepo) Authenticate(ctx context.Context, username, password string) (success bool, uid string, defaultNamespace string, err error) {
	txn := s.Dg.NewReadOnlyTxn()

	const q = `query authenticate($username: string, $password: string){
  login(func: eq(username, $username)) @filter(eq(type, "credentials")) {
    uid
    checkpwd(password, $password)
    ~has.credentials {
      uid
      type
      default.namespace{
        uid
        name
      }
    }
  }
}
`

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$username": username, "$password": password})
	if err != nil {
		return false, "", "", err
	}

	var result struct {
		Login []*UsernameCredential `json:"login"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false, "", "", err
	}

	if len(result.Login) > 0 {
		login := result.Login[0]
		if login.CheckPwd {
			// Success
			if len(login.Account) > 0 {
				if len(login.Account[0].DefaultNamespace) > 0 {
					defaultNamespace = login.Account[0].DefaultNamespace[0].Name
				}
				return result.Login[0].CheckPwd, login.Account[0].UID, defaultNamespace, nil
			}
		}
	}
	return false, "", "", errors.New("Invalid credentials")
}

func (s *DGraphRepo) ListAccounts(ctx context.Context) (accounts []*nodepb.Account, err error) {
	txn := s.Dg.NewReadOnlyTxn()

	const q = `query accounts{
                     accounts(func: eq(type, "account")) {
                       uid
                       name
                     }
                   }`

	res, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var result struct {
		Accounts []*Account `json:"accounts"`
	}

	if err := json.Unmarshal(res.Json, &result); err != nil {
		return nil, err
	}

	for _, account := range result.Accounts {
		accounts = append(accounts, &nodepb.Account{
			Uid:    account.UID,
			Name:   account.Name,
			IsRoot: account.IsRoot,
		})
	}

	return accounts, nil
}

func (s *DGraphRepo) CreateUserAccount(ctx context.Context, username, password string, isRoot bool) (uid string, err error) {
	defaultNs, err := s.CreateNamespace(ctx, username)
	if err != nil {
		return "", err
	}
	txn := s.Dg.NewTxn()

	q := `query userExists($name: string) {
                exists(func: eq(name, $name)) @filter(eq(type, "account")) {
                  uid
                }
              }
             `

	var result struct {
		Exists []map[string]interface{} `json:"exists"`
	}

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$name": username})
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return "", err
	}

	if len(result.Exists) == 0 {
		js, err := json.Marshal(&Account{
			Node: Node{
				Type: "account",
				UID:  "_:user",
			},
			Name:   username,
			IsRoot: isRoot,
			HasCredentials: &UsernameCredential{
				Node: Node{
					Type: "credentials",
				},
				Username: username,
				Password: password,
			},
			AccessToNamespace: []*Namespace{
				&Namespace{
					Node: Node{
						UID: defaultNs,
					},
					AccessToPermission: nodepb.Action_WRITE.String(),
				},
			},
			DefaultNamespace: []*Namespace{
				&Namespace{
					Node: Node{
						UID: defaultNs,
					},
				},
			},
		})
		if err != nil {
			return "", err
		}
		m := &api.Mutation{SetJson: js}
		a, err := txn.Mutate(ctx, m)
		if err != nil {
			return "", err
		}

		err = txn.Commit(ctx)
		if err != nil {
			return "", errors.New("Failed to commit")
		}
		userUID := a.GetUids()["user"]
		return userUID, nil

	}
	return "", errors.New("User exists already")
}

func (s *DGraphRepo) Authorize(ctx context.Context, account, node, action string, inherit bool) (err error) {
	txn := s.Dg.NewTxn()

	if ok := checkType(ctx, txn, account, "account"); !ok {
		return errors.New("invalid account")
	}

	var _type string
	if ok := checkType(ctx, txn, node, "namespace"); !ok {
		if ok := checkType(ctx, txn, node, "object"); !ok {
			return errors.New("resource does not exist")
		} else {
			_type = "object"
		}
	} else {
		_type = "namespace"
	}

	in := Account{
		Node: Node{
			UID: account,
		},
	}

	if _type == "namespace" {
		// ignore inherit flag; access to a namespace is always recursive
		in.AccessToNamespace = []*Namespace{
			&Namespace{
				Node: Node{
					UID: node,
				},
				AccessToPermission: action,
			},
		}
	} else if _type == "object" {
		in.AccessTo = []*Object{
			&Object{
				Node: Node{
					UID: node,
				},
				AccessToPermission: action,
				AccessToInherit:    inherit,
			},
		}
	}

	js, err := json.Marshal(&in)
	if err != nil {
		return err
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return errors.New("Failed to mutate")
	}
	return nil
}

func (s *DGraphRepo) GetAccount(ctx context.Context, name string) (account *nodepb.Account, err error) {
	txn := s.Dg.NewReadOnlyTxn()
	const q = `query accounts($account: string) {
                     accounts(func: uid($account)) @filter(eq(type, "account"))  {
                       uid
                       name
                       type
                       isRoot
                     }
                   }`

	response, err := txn.QueryWithVars(ctx, q, map[string]string{"$account": name})
	if err != nil {
		return nil, err
	}

	var result struct {
		Account []*Account `json:"accounts"`
	}

	err = json.Unmarshal(response.Json, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Account) == 0 {
		return nil, errors.New("Account not found")
	}

	return &nodepb.Account{
		Uid:    result.Account[0].UID,
		Name:   result.Account[0].Name,
		IsRoot: result.Account[0].IsRoot,
	}, err
}

func (s *DGraphRepo) CreateNamespace(ctx context.Context, name string) (id string, err error) {
	ns := &Namespace{
		Node: Node{
			Type: "namespace",
			UID:  "_:namespace",
		},
		Name: name,
	}

	txn := s.Dg.NewTxn()
	js, err := json.Marshal(&ns)
	if err != nil {
		return "", err
	}

	assigned, err := txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}

	return assigned.GetUids()["namespace"], nil
}

func (s *DGraphRepo) GetNamespace(ctx context.Context, namespaceID string) (namespace *nodepb.Namespace, err error) {
	const q = `query getNamespaces($namespace: string) {
                     namespaces(func: eq(name, $namespace)) @filter(eq(type, "namespace"))  {
	               uid
                       name
	             }
                   }`

	res, err := s.Dg.NewReadOnlyTxn().QueryWithVars(ctx, q, map[string]string{"$namespace": namespaceID})
	if err != nil {
		return nil, err
	}

	var resultSet struct {
		Namespaces []*Namespace `json:"namespaces"`
	}

	if err := json.Unmarshal(res.Json, &resultSet); err != nil {
		return nil, err
	}

	if len(resultSet.Namespaces) > 0 {
		return &nodepb.Namespace{
			Id:   resultSet.Namespaces[0].UID,
			Name: resultSet.Namespaces[0].Name,
		}, nil
	}

	return nil, errors.New("Namespace not found")
}

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

	fmt.Println("pars", params)

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

	fmt.Println("acc", len(access.Access) > 0)

	return len(access.Access) > 0, nil
}

func (s *DGraphRepo) IsAuthorized(ctx context.Context, node, account, action string) (decision bool, err error) {
	fmt.Println("Isa")
	if node == account {
		return true, nil
	}

	params := map[string]string{
		"$device_id": node,
		"$user_id":   account,
	}

	txn := s.Dg.NewReadOnlyTxn()

	const qDirect = `query direct_access($device_id: string, $user_id: string){
                         direct(func: uid($user_id)) @normalize @cascade {
                           access.to  @filter(uid($device_id)) @facets(permission,inherit) {
                             type: type
                           }
                         }
                         direct_via_one_object(func: uid($user_id)) @normalize @cascade {
                           access.to @facets(permission,inherit) {
                             contains @filter(uid($device_id)) {
                               uid
                               type: type
                             }
                           }
                         }
                        }`

	res, err := txn.QueryWithVars(ctx, qDirect, params)
	if err != nil {
		return false, err
	}

	var permissions struct {
		Direct          []Object `json:"direct"`
		DirectViaObject []Object `json:"direct_via_one_object"`
	}

	err = json.Unmarshal(res.Json, &permissions)
	if err != nil {
		return false, err
	}

	if len(permissions.Direct) > 0 {
		if isPermissionSufficient(action, permissions.Direct[0].AccessToPermission) {
			return true, nil
		}
	}

	if len(permissions.DirectViaObject) > 0 {
		if isPermissionSufficient(action, permissions.DirectViaObject[0].AccessToPermission) {
			return true, nil
		}
	}

	const qRecursiveWrite = `query recursive($user_id: string, $device_id: string){
                         shortest(from: $user_id, to: $device_id) {
                           access.to @facets(eq(inherit, true) AND eq(permission,"WRITE"))
                           access.to.namespace @facets(eq(permission,"WRITE"))
                           owns
                           children
                         }
                       }`

	const qRecursiveRead = `query recursive($user_id: string, $device_id: string){
                         shortest(from: $user_id, to: $device_id) {
                           access.to @facets(eq(inherit, true) AND (eq(permission,"WRITE") OR eq(permission, "READ")))
                           access.to.namespace @facets((eq(permission,"WRITE") OR eq(permission, "READ")))
                           owns
                           children
                         }
                       }`

	var qRecursive string

	switch action {
	case "READ":
		qRecursive = qRecursiveRead
	case "WRITE":
		qRecursive = qRecursiveWrite
	default:
		return false, errors.New("Invalid action")
	}

	res, err = txn.QueryWithVars(ctx, qRecursive, params)
	if err != nil {
		return false, err
	}

	var rez struct {
		Path []map[string]interface{} `json:"_path_"`
	}

	if err = json.Unmarshal(res.Json, &rez); err != nil {
		return false, err
	}

	if len(rez.Path) > 0 {
		return true, nil
	}

	return false, nil
}
