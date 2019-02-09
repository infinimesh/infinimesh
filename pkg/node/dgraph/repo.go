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
	"github.com/infinimesh/infinimesh/pkg/tools"
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

type dGraphRepo struct {
	dg *dgo.Dgraph
}

func NewDGraphRepo(dg *dgo.Dgraph) node.Repo {
	return &dGraphRepo{dg: dg}
}

func checkExists(ctx context.Context, txn *dgo.Txn, uid, _type string) bool {
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

func (s *dGraphRepo) Authenticate(ctx context.Context, username, password string) (success bool, uid string, err error) {
	txn := s.dg.NewReadOnlyTxn()

	const q = `query authenticate($username: string, $password: string){
  login(func: eq(username, $username)) @filter(eq(type, "credentials")) {
    uid
    checkpwd(password, $password)
    ~has.credentials {
      uid
      type
    }
  }
}
`

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$username": username, "$password": password})
	if err != nil {
		return false, "", err
	}

	var result struct {
		Login []*UsernameCredential `json:"login"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return false, "", err
	}

	if len(result.Login) > 0 {
		login := result.Login[0]
		if login.CheckPwd {
			// Success
			if len(login.Account) > 0 {
				return result.Login[0].CheckPwd, login.Account[0].UID, nil
			}
		}
	}
	return false, "", errors.New("Invalid credentials")
}

func (s *dGraphRepo) CreateAccount(ctx context.Context, username, password string) (uid string, err error) {
	txn := s.dg.NewTxn()

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
			Name: username,
			HasCredentials: &UsernameCredential{
				Node: Node{
					Type: "credentials",
				},
				Username: username,
				Password: password,
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

func (s *dGraphRepo) Authorize(ctx context.Context, account, node, action string, inherit bool) (err error) {
	txn := s.dg.NewTxn()

	if ok := checkExists(ctx, txn, account, "account"); !ok {
		return errors.New("account does not exist")
	}

	if ok := checkExists(ctx, txn, node, "object"); !ok {
		if ok := checkExists(ctx, txn, node, "device"); !ok {
			return errors.New("resource does not exist")
		}
	}

	in := Account{
		Node: Node{
			UID: account,
		},
		AccessTo: []*ObjectList{
			&ObjectList{
				Node: Node{
					UID: node,
				},
				AccessToPermission: action,
				AccessToInherit:    inherit,
			},
		},
	}

	fmt.Println("axx")
	tools.PrettyPrint(in)

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

func (s *dGraphRepo) GetAccount(ctx context.Context, name string) (account *nodepb.Account, err error) {
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

	// return result.Account[0], nil
	return &nodepb.Account{
		Uid:  result.Account[0].UID,
		Name: result.Account[0].Name,
	}, err
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
		Object []*ObjectList
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

func addDeletesRecursively(mu *api.Mutation, items []*ObjectList) {
	for _, item := range items {
		dgo.DeleteEdges(mu, item.UID, "_STAR_ALL")
		for _, object := range item.Contains {
			dgo.DeleteEdges(mu, object.UID, "_STAR_ALL")
		}
		addDeletesRecursively(mu, item.Contains)
	}
}

func (s *dGraphRepo) CreateObject(ctx context.Context, name, parent, kind, namespace string) (id string, err error) {
	var newObject *ObjectList
	if parent == "" {
		newObject = &ObjectList{
			Node: Node{
				UID:  "_:new",
				Type: kind,
			},
			Name: name,
		}
	} else {
		newObject = &ObjectList{
			Node: Node{
				UID: parent,
			},
			Contains: []*ObjectList{
				&ObjectList{
					Node: Node{
						UID:  "_:new",
						Type: kind,
					},
					Name: name,
				},
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

func (s *dGraphRepo) ListForAccount(ctx context.Context, account string) (directDevices []*nodepb.Device, directObjects []*nodepb.Object, inheritedObjects []*nodepb.Object, err error) {
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
                     }
                   }
                  }`

	var result struct {
		Inherited []ObjectList `json:"inherited"`
		Direct    []struct {
			AccessTo []ObjectList `json:"access.to"`
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

	var roots []ObjectList

	// Access grants
	for _, accessObject := range result.Inherited {

		var isRoot = true
		for _, other := range result.Inherited {
			if other.UID != accessObject.UID {
				if isSub := isSubtreeOf(&accessObject, &other); isSub {
					isRoot = false
				}

			}
		}

		if isRoot {
			fmt.Println(accessObject.Name, " is root")
			roots = append(roots, accessObject)
		}

	}

	if len(result.Direct) > 0 {
		for _, directObject := range result.Direct[0].AccessTo {
			directObjects = append(directObjects, mapObject(&directObject))
		}
	}

	// inheritedObjects = roots

	for _, root := range roots {
		inheritedObjects = append(inheritedObjects, mapObject(&root))
	}

	return
}

func mapObject(o *ObjectList) *nodepb.Object {
	objects := make([]*nodepb.Object, 0)
	if len(o.Contains) > 0 {
		for _, v := range o.Contains {
			object := mapObject(v)
			objects = append(objects, object)

		}
	}

	res := &nodepb.Object{
		Uid:     o.UID,
		Name:    o.Name,
		Objects: objects,
	}

	return res
}

func isSubtreeOf(tree, other *ObjectList) bool {
	if tree.UID == other.UID {
		return true
	}

	// We assume that it's sufficient to check if the root is contained in
	// the other tree. If this is the case, the subtree is being merged into
	// the detected enclosing tree
	for i := range other.Contains {
		otherChild := other.Contains[i]
		if sub := isSubtreeOf(tree, otherChild); sub {
			// we're part of the other tree -> merge into the other
			// (so data which is maybe only in this tree, but not
			// the target) and flag ourself as subree
			mergeInto(tree, otherChild)
			return true

		}
	}
	return false
}

func mergeInto(source, target *ObjectList) {
	targetMap := make(map[string]*ObjectList)
	for _, targetNode := range target.Contains {
		targetMap[target.UID] = targetNode
	}

	for _, sourceNode := range source.Contains {
		if _, exists := targetMap[sourceNode.UID]; exists {
			mergeInto(sourceNode, targetMap[sourceNode.UID])
		} else {
			target.Contains = append(target.Contains, sourceNode)
		}
	}
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
		Direct          []ObjectList `json:"direct"`
		DirectViaObject []ObjectList `json:"direct_via_one_object"`
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
                         }
                       }`

	const qRecursiveRead = `query recursive($user_id: string, $device_id: string){
                         shortest(from: $user_id, to: $device_id) {
                           access.to @facets(eq(inherit, true) AND (eq(permission,"WRITE") OR eq(permission, "READ"))) @filter(eq(type, "object"))
                           contains
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
