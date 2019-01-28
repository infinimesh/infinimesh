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
	DeleteObject(ctx context.Context, uid string) (err error)
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
	const q = `
query accounts($account: string) {
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

func (s *dGraphRepo) DeleteObject(ctx context.Context, uid string) (err error) {
	txn := s.dg.NewTxn()

	// Find target node
	const q = `
	query deleteObject($root: string){
	  object(func: uid($root)) {
	    uid
	    name
	    contains {
	      uid
	    }
	    ~contains { # Parent
	      uid
	    name
	    }
            ~access.to {
              uid
              name
            }
	  }
	}
	`

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$root": uid,
	})
	if err != nil {
		return err
	}

	var result struct {
		Objects []*ObjectList `json:"object"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return err
	}

	mu := &api.Mutation{}

	if len(result.Objects) == 0 {
		return errors.New("unexpected response from DB: 0 objects founds")
	}

	// Detect parent by ~contains edge
	toDelete := result.Objects[0]
	if len(toDelete.ContainedIn) > 0 {
		parent := toDelete.ContainedIn[0]
		mu.Del = append(mu.Del, &api.NQuad{
			Subject:   parent.UID,
			Predicate: "contains",
			ObjectId:  toDelete.UID,
		})

	}

	// Detect parent by ~access.to edge
	if len(toDelete.AccessedBy) > 0 {
		parent := toDelete.AccessedBy[0]
		mu.Del = append(mu.Del, &api.NQuad{
			Subject:   parent.UID,
			Predicate: "access.to",
			ObjectId:  toDelete.UID,
		})

	}

	// Find and delete all edges & nodes below this node, including this
	// node
	const qChilds = `query children($root: string){
			  object(func: uid($root)) @recurse {
			    uid
			    contains {
			    }
                            contains.device {
                            }
			  }
			}`

	res, err := txn.QueryWithVars(ctx, qChilds, map[string]string{
		"$root": toDelete.UID,
	})
	if err != nil {
		return err
	}

	var resultChildren struct {
		Object []ObjectList
	}

	err = json.Unmarshal(res.Json, &resultChildren)
	if err != nil {
		return err
	}

	addDeletesRecursively(mu, resultChildren.Object)

	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		return err
	}

	return txn.Commit(ctx)
}

func addDeletesRecursively(mu *api.Mutation, items []ObjectList) {
	for _, item := range items {
		dgo.DeleteEdges(mu, item.UID, "_STAR_ALL")
		for _, device := range item.ContainsDevice {
			dgo.DeleteEdges(mu, device.UID, "_STAR_ALL")
		}
		addDeletesRecursively(mu, item.Contains)
	}
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

	// TODO check if parent node is of correct type (node)
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
