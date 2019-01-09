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
}

type dGraphRepo struct {
	dg *dgo.Dgraph
}

const ContextKeyAccount = "infinimesh/pkg/node/account"

func NewDGraphRepo(dg *dgo.Dgraph) Repo {
	return &dGraphRepo{dg: dg}
}

func (s *dGraphRepo) CreateObject(ctx context.Context, name, parent string) (id string, err error) {
	var newObject *Resource
	if parent == "" {
		newObject = &Resource{
			Node: Node{
				UID:  "_:new",
				Type: "object",
			},
			Name: name,
		}
	} else {
		newObject = &Resource{
			Node: Node{
				UID: parent,
			},
			Contains: &Resource{
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
		Direct          []Resource `json:"direct"`
		DirectViaObject []Resource `json:"direct_via_one_object"`
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
