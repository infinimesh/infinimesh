package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

type Repo interface {
	CreateAccount(ctx context.Context, username, password string) (uid string, err error)
	IsAuthorized(ctx context.Context, target, who, action string) (decision bool, err error)
	Authorize(ctx context.Context, account, node, action string, inherit bool) (err error)
	GetAccount(ctx context.Context, name string) (account *Account, err error)
	Authenticate(ctx context.Context, username, password string) (success bool, uid string, err error)

	CreateObject(ctx context.Context, name, parent string) (id string, err error)
	DeleteObject(ctx context.Context, uid string) (err error)
	ListForAccount(ctx context.Context, account string) (directDevices []Device, directObjects []ObjectList, inheritedObjects []ObjectList, err error)
}

type dGraphRepo struct {
	dg *dgo.Dgraph
}

func NewDGraphRepo(dg *dgo.Dgraph) Repo {
	return &dGraphRepo{dg: dg}
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
		AccessTo: &Object{
			Node: Node{
				UID: node,
			},
			AccessToPermission: action,
			AccessToInherit:    inherit,
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
		directDevices = result.Direct[0].AccessToDevice
		directObjects = result.Direct[0].AccessTo
	}

	inheritedObjects = roots

	return
}

func isSubtreeOf(tree, other *ObjectList) bool {
	if tree.UID == other.UID {
		return true
	}

	// We assume that it's sufficient to check if the root is contained in
	// the other tree. If this is the case, the subtree is being merged into
	// the detected enclosing tree
	for i := range other.Contains {
		otherChild := &other.Contains[i]
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
	targetDevices := make(map[string]*Device)
	for _, device := range target.ContainsDevice {
		targetDevices[device.UID] = &device
	}

	targetMap := make(map[string]*ObjectList)
	for _, targetNode := range target.Contains {
		targetMap[target.UID] = &targetNode
	}

	for _, sourceDevice := range source.ContainsDevice {
		if _, exists := targetDevices[sourceDevice.UID]; !exists {
			target.ContainsDevice = append(target.ContainsDevice, sourceDevice)
		}
	}

	for _, sourceNode := range source.Contains {
		if _, exists := targetMap[sourceNode.UID]; exists {
			mergeInto(&sourceNode, targetMap[sourceNode.UID])
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
                         direct_device(func: uid($user_id)) @normalize @cascade {
                           access.to.device  @filter(uid($device_id)) @facets(permission,inherit) {
                             type: type
                           }
                         }
                         direct_via_one_object(func: uid($user_id)) @normalize @cascade {
                           access.to @filter(eq(type, "object")) @facets(permission,inherit) {
                             contains @filter(uid($device_id)) {
                               uid
                               type: type
                             }
                           }
                         }
                         direct_device_via_one_object(func: uid($user_id)) @normalize @cascade {
                           access.to @filter(eq(type, "object")) @facets(permission,inherit) {
                             contains.device @filter(uid($device_id)) {
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
		Direct                []Object `json:"direct"`
		DirectDevice          []Device `json:"direct_device"`
		DirectViaObject       []Object `json:"direct_via_one_object"`
		DirectDeviceViaObject []Object `json:"direct_device_via_one_object"`
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

	if len(permissions.DirectDevice) > 0 {
		if isPermissionSufficient(action, permissions.DirectDevice[0].AccessToDevicePermission) {
			return true, nil
		}
	}

	if len(permissions.DirectViaObject) > 0 {
		if isPermissionSufficient(action, permissions.DirectViaObject[0].AccessToPermission) {
			return true, nil
		}
	}

	if len(permissions.DirectDeviceViaObject) > 0 {
		if isPermissionSufficient(action, permissions.DirectDeviceViaObject[0].AccessToPermission) {
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
