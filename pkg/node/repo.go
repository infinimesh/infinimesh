package node

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

type Repo interface {
	IsAuthorized(ctx context.Context, target, who, action string) (decision bool, err error)
	CreateObject(ctx context.Context, name, parent string) (id string, err error)
	ListForAccount(ctx context.Context, account string) (directDevices []Device, directObjects []ObjectList, inheritedObjects []ObjectList, err error)
	GetAccount(ctx context.Context, name string) (account *Account, err error)
}

type dGraphRepo struct {
	dg *dgo.Dgraph
}

func NewDGraphRepo(dg *dgo.Dgraph) Repo {
	return &dGraphRepo{dg: dg}
}

func (s *dGraphRepo) GetAccount(ctx context.Context, name string) (account *Account, err error) {
	txn := s.dg.NewReadOnlyTxn()
	const q = `query accounts($account: string) {
  accounts(func: eq(name, $account)) @filter(eq(type, "account"))  {
    uid
    name
    type
  }
}
`
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

	return result.Account[0], nil
}

func (s *dGraphRepo) CreateObject(ctx context.Context, name, parent string) (id string, err error) {
	var newObject *Object
	if parent == "" {
		newObject = &Object{
			Node: Node{
				UID:  "_:new",
				Type: "object",
			},
			Name: name,
		}
	} else {
		newObject = &Object{
			Node: Node{
				UID: parent,
			},
			Contains: &Object{
				Node: Node{
					UID:  "_:new",
					Type: "object",
				},
				Name: name,
			},
		}
	}

	js, err := json.Marshal(&newObject)
	if err != nil {
		return "", err
	}

	a, err := s.dg.NewTxn().Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return "", err
	}

	return a.GetUids()["new"], nil

}

func (s *dGraphRepo) ListForAccount(ctx context.Context, account string) (directDevices []Device, directObjects []ObjectList, inheritedObjects []ObjectList, err error) {
	txn := s.dg.NewReadOnlyTxn()

	const q = `query list($account: string) {
                   var(func: uid($account)) {
                     access.to @facets(eq(inherit,true)) {
                       OBJS as uid
                       name
                     }
                   }

                   inherited(func: uid(OBJS)) @recurse {
                     contains{} 
                     contains.device{}
                     uid
                     type
                     name
                   }

                   direct(func: uid($account)) {
                   # Via enclosing object
                     access.to @facets(eq(inherit,false)) {
                       uid
                       name
                       type
                       contains.device @filter(eq(type, "device")) {
                         uid
                         name
                         type
                       }
                     }

                     # Via direct permission on device
                     access.to.device @filter(eq(type, "device")) {
                       uid
                       name
                       type
                     }
                   }
                  }`

	var result struct {
		Inherited []ObjectList `json:"inherited"`
		Direct    []struct {
			AccessTo       []ObjectList `json:"access.to"`
			AccessToDevice []Device     `json:"access.to.device"`
		} `json:"direct"`
	}

	params := map[string]string{
		"$account": account,
	}

	res, err := txn.QueryWithVars(ctx, q, params)
	if err != nil {
		return nil, nil, nil, err
	}

	err = json.Unmarshal(res.Json, &result)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(result.Direct) > 0 {
		directDevices = result.Direct[0].AccessToDevice
		directObjects = result.Direct[0].AccessTo
	}

	inheritedObjects = result.Inherited

	return
}

func (s *dGraphRepo) IsAuthorized(ctx context.Context, node, account, action string) (decision bool, err error) {
	if node == account {
		return true, nil
	}

	params := map[string]string{
		"$device_id": node,
		"$user_id":   account,
	}

	txn := s.dg.NewReadOnlyTxn()

	const qDirect = `query direct_access($device_id: string, $user_id: string){
                         direct(func: uid($user_id)) @normalize @cascade {
                           access.to  @filter(uid($device_id) AND eq(type, "device")) @facets(permission,inherit) {
                             type: type
                           }
                         }
                         direct_via_one_object(func: uid($user_id)) @normalize @cascade {
                           access.to @filter(eq(type, "object")) @facets(permission,inherit) {
                             contains @filter(uid($device_id) AND eq(type, "device")) {
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
                           access.to @facets(eq(inherit, true) AND eq(permission,"WRITE")) @filter(eq(type, "object"))
                           contains @filter(eq(type, "object"))
                           contains.device @filter(eq(type, "device"))
                         }
                       }`

	const qRecursiveRead = `query recursive($user_id: string, $device_id: string){
                         shortest(from: $user_id, to: $device_id) {
                           access.to @facets(eq(inherit, true) AND (eq(permission,"WRITE") OR eq(permission, "READ"))) @filter(eq(type, "object"))
                           contains @filter(eq(type, "object"))
                           contains.device @filter(eq(type, "device"))
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
