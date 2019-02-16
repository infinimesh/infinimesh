package dgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

func checkKind(ctx context.Context, txn *dgo.Txn, uid, _type string) bool { //nolint
	q := `query object($_uid: string) {
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
		Objects []*Object `json:"object"`
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
	if len(toDelete.Parent) > 0 {
		parent := toDelete.Parent[0]
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
			  }
			}`

	res, err := txn.QueryWithVars(ctx, qChilds, map[string]string{
		"$root": toDelete.UID,
	})
	if err != nil {
		return err
	}

	var resultChildren struct {
		Object []*Object
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

func addDeletesRecursively(mu *api.Mutation, items []*Object) {
	for _, item := range items {
		dgo.DeleteEdges(mu, item.UID, "_STAR_ALL")
		for _, object := range item.Children {
			dgo.DeleteEdges(mu, object.UID, "_STAR_ALL")
		}
		addDeletesRecursively(mu, item.Children)
	}
}

func (s *dGraphRepo) CreateObject(ctx context.Context, name, parentID, kind, namespaceID string) (id string, err error) {
	txn := s.dg.NewTxn()

	if ok := checkType(ctx, txn, namespaceID, "namespace"); !ok {
		return "", errors.New("Invalid namespace")
	}

	newObject := &Object{
		Node: Node{
			UID:  "_:new",
			Type: "object",
		},
		Name: name,
		Kind: kind,
	}

	var object *Object
	if parentID == "" {
		object = newObject
	} else {
		if ok := checkType(ctx, txn, parentID, "object"); !ok {
			return "", errors.New("Invalid parent")
		}

		object = &Object{
			Node: Node{
				UID: parentID,
			},
			Children: []*Object{
				newObject,
			},
		}
	}

	js, err := json.Marshal(&object)
	if err != nil {
		return "", err
	}

	a, err := txn.Mutate(ctx, &api.Mutation{
		SetJson: js,
	})
	if err != nil {
		return "", err
	}

	newUID := a.GetUids()["new"]

	ns := &api.NQuad{
		Subject:   namespaceID,
		Predicate: "owns",
		ObjectId:  newUID,
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		Set: []*api.NQuad{
			ns,
		},
	})
	if err != nil {
		return "", err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return "", err
	}

	return a.GetUids()["new"], nil
}

func (s *dGraphRepo) ListForAccount(ctx context.Context, account string) (inheritedObjects []*nodepb.Object, err error) {
	txn := s.dg.NewReadOnlyTxn()

	const q = `query list($account: string) {
                   var(func: uid($account)) {
                     access.to @facets(eq(inherit,true)) {
                       OBJS as uid
                       name
                     }
                   }

                   inherited(func: uid(OBJS)) @recurse {
                     children{} 
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
		Inherited []Object `json:"inherited"`
		Direct    []struct {
			AccessTo []Object `json:"access.to"`
		} `json:"direct"`
	}

	params := map[string]string{
		"$account": account,
	}

	res, err := txn.QueryWithVars(ctx, q, params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Json, &result)
	if err != nil {
		return nil, err
	}

	var roots []Object

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

	// if len(result.Direct) > 0 {
	// 	for _, directObject := range result.Direct[0].AccessTo {
	// 		directObjects = append(directObjects, mapObject(&directObject))
	// 	}
	// }

	for _, root := range roots {
		inheritedObjects = append(inheritedObjects, mapObject(&root))
	}

	return inheritedObjects, nil
}

func mapObject(o *Object) *nodepb.Object {
	objects := make([]*nodepb.Object, 0)
	if len(o.Children) > 0 {
		for _, v := range o.Children {
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

func isSubtreeOf(tree, other *Object) bool {
	if tree.UID == other.UID {
		return true
	}

	// We assume that it's sufficient to check if the root is contained in
	// the other tree. If this is the case, the subtree is being merged into
	// the detected enclosing tree
	for i := range other.Children {
		otherChild := other.Children[i]
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

func mergeInto(source, target *Object) {
	targetMap := make(map[string]*Object)
	for _, targetNode := range target.Children {
		targetMap[target.UID] = targetNode
	}

	for _, sourceNode := range source.Children {
		if _, exists := targetMap[sourceNode.UID]; exists {
			mergeInto(sourceNode, targetMap[sourceNode.UID])
		} else {
			target.Children = append(target.Children, sourceNode)
		}
	}
}
