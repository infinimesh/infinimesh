/*
Copyright © 2021-2023 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/go-redis/redis/v8"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/node/access"
	accpb "github.com/infinimesh/proto/node/accounts"
	nspb "github.com/infinimesh/proto/node/namespaces"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Access struct {
	From  driver.DocumentID `json:"_from"`
	To    driver.DocumentID `json:"_to"`
	Level access.Level      `json:"level"`
	Role  access.Role       `json:"role,omitempty"`

	driver.DocumentMeta
}

type InfinimeshGraphNode interface {
	GetUuid() string
	ID() driver.DocumentID

	SetAccessLevel(level access.Level)
	GetAccess() *access.Access
}

type InfinimeshController interface {
	_DB() driver.Database
	_log() *zap.Logger
}

type InfinimeshBaseController struct {
	log *zap.Logger
	db  driver.Database
}

func (c *InfinimeshBaseController) _DB() driver.Database {
	return c.db
}

func (c *InfinimeshBaseController) _log() *zap.Logger {
	return c.log
}

type InfinimeshCommonActionsRepo interface {
	GetEdgeCol(ctx context.Context, name string) driver.Collection
	CheckLink(ctx context.Context, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode) bool
	Link(
		ctx context.Context, log *zap.Logger, edge driver.Collection,
		from InfinimeshGraphNode, to InfinimeshGraphNode,
		lvl access.Level, role access.Role,
	) error
	Move(ctx context.Context, c InfinimeshController, obj InfinimeshGraphNode, edge driver.Collection, ns string) error
	//
	EnsureRootExists(_log *zap.Logger, rdb *redis.Client, passwd string) (err error)
}

type infinimeshCommonActionsRepo struct {
	db driver.Database
}

func NewInfinimeshCommonActionsRepo(db driver.Database) InfinimeshCommonActionsRepo {
	return &infinimeshCommonActionsRepo{db: db}
}

func (r *infinimeshCommonActionsRepo) GetEdgeCol(ctx context.Context, name string) driver.Collection {
	g, _ := r.db.Graph(ctx, schema.PERMISSIONS_GRAPH.Name)
	col, _, _ := g.EdgeCollection(ctx, name)
	return col
}

func (r *infinimeshCommonActionsRepo) CheckLink(ctx context.Context, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode) bool {
	res, err := edge.DocumentExists(ctx, from.ID().Key()+"-"+to.ID().Key())
	return err == nil && res
}

func (r *infinimeshCommonActionsRepo) Link(ctx context.Context, log *zap.Logger, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode, lvl access.Level, role access.Role) error {
	log.Debug("Linking two nodes",
		zap.Any("from", from.ID()),
		zap.Any("to", to.ID()),
		zap.Any("level", lvl),
		zap.Any("role", role),
	)

	a := Access{
		From:  from.ID(),
		To:    to.ID(),
		Level: lvl,
		Role:  role,
		DocumentMeta: driver.DocumentMeta{
			Key: from.ID().Key() + "-" + to.ID().Key(),
		},
	}

	if a.Level == access.Level_NONE {
		_, err := edge.RemoveDocument(ctx, a.Key)
		return err
	}

	if _, err := edge.UpdateDocument(ctx, a.DocumentMeta.Key, a); err == nil {
		return nil
	}

	_, err := edge.CreateDocument(ctx, a)
	return err
}

func (r *infinimeshCommonActionsRepo) Move(ctx context.Context, c InfinimeshController, obj InfinimeshGraphNode, edge driver.Collection, ns string) error {
	log := c._log().Named("Move")
	log.Debug("Move request received", zap.Any("object", obj), zap.String("namespace", ns))

	requestor := NewBlankAccountDocument(
		ctx.Value(inf.InfinimeshAccountCtxKey).(string),
	)
	log.Debug("Requestor", zap.String("id", requestor.Key))

	err := AccessLevelAndGet(ctx, log, c._DB(), requestor, obj)
	if err != nil {
		return status.Error(codes.NotFound, "Object not found or not enough Access Rights")
	}
	role := obj.GetAccess().Role
	if role != access.Role_OWNER && obj.GetAccess().Level != access.Level_ROOT {
		return status.Error(codes.PermissionDenied, "Must be Owner or Root to perform Move")
	}
	if obj.GetAccess().Namespace == nil {
		return status.Error(codes.Internal, "Object is not under any Namespace, contact support")
	}

	old_namespace := NewBlankNamespaceDocument(*obj.GetAccess().Namespace)

	namespace := NewBlankNamespaceDocument(ns)
	err = AccessLevelAndGet(ctx, log, c._DB(), requestor, namespace)
	if err != nil {
		return status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if namespace.GetAccess().Role != access.Role_OWNER && namespace.GetAccess().Level != access.Level_ROOT {
		return status.Error(codes.PermissionDenied, "Must be Owner or Root to perform Move")
	}

	err = r.Link(ctx, log, edge, old_namespace, obj, access.Level_NONE, access.Role_UNSET)
	if err != nil {
		log.Warn("Error unlinking Object from Namespace",
			zap.String("object", obj.ID().String()),
			zap.String("namespace", old_namespace.Key),
			zap.Error(err))
		return status.Error(codes.Internal, "Couldn't unlink the object")
	}

	err = r.Link(ctx, log, edge, namespace, obj, access.Level_ADMIN, role)
	if err != nil {
		log.Warn("Error linking Object to Namespace",
			zap.String("object", obj.ID().String()),
			zap.String("namespace", namespace.Key),
			zap.Error(err))
		return status.Error(codes.Internal, "Couldn't link the object, contact support")
	}

	return nil
}

func NewBlankDocument(col string, key string) driver.DocumentMeta {
	return driver.DocumentMeta{
		Key: key,
		ID:  driver.NewDocumentID(col, key),
	}
}

const getWithAccessLevelRoleAndNS = `
FOR path IN OUTBOUND K_SHORTEST_PATHS @account TO @node
GRAPH @permissions SORT path.edges[0].level DESC
    LET perm = path.edges[0]
    LET last = path.edges[-1]
	RETURN MERGE(
	    path.vertices[-1], {
	        access: {
	            level: last.role == 2 ? last.level : perm.level,
	            role: last.role == 2 ? last.role : perm.role,
	            namespace: path.vertices[-2]._key
	        }
	    }
    )
`

func AccessLevelAndGet(ctx context.Context, log *zap.Logger, db driver.Database, account *Account, node InfinimeshGraphNode) error {
	vars := map[string]interface{}{
		"account":     account.ID(),
		"node":        node.ID(),
		"permissions": schema.PERMISSIONS_GRAPH.Name,
	}
	c, err := db.Query(ctx, getWithAccessLevelRoleAndNS, vars)
	if err != nil {
		log.Debug("Error while executing query", zap.Any("vars", vars), zap.Error(err))
		return err
	}
	defer c.Close()

	_, err = c.ReadDocument(ctx, &node)
	if err != nil {
		log.Debug("Error while reading node document", zap.Error(err))
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	if account.ID() == node.ID() {
		node.SetAccessLevel(access.Level_ROOT)
	}

	return nil
}

const listObjectsOfKind = `
FOR node, edge, path IN 0..@depth OUTBOUND @from
GRAPH @permissions_graph
OPTIONS {order: "bfs", uniqueVertices: "global"}
FILTER IS_SAME_COLLECTION(@@kind, node)
FILTER edge.level > 0
%s
    LET perm = path.edges[0]
    LET last = path.edges[-1]
	RETURN MERGE(node, {
	    uuid: node._key,
	    access: {
	        level: last.role == 2 ? last.level : perm.level,
	        role:  last.role == 2 ? last.role : perm.role,
	        namespace: last.role == 2 ? null : path.vertices[-2]._key
	     }
	    }
    )
`

// List children nodes
// ctx - context
// log - logger
// db - Database connection
// from - Graph node to start traversal from
// children - children type(collection name)
// depth
func ListQuery(ctx context.Context, log *zap.Logger, db driver.Database, from InfinimeshGraphNode, children string) (driver.Cursor, error) {
	bindVars := map[string]interface{}{
		"depth":             DepthValue(ctx),
		"from":              from.ID(),
		"permissions_graph": schema.PERMISSIONS_GRAPH.Name,
		"@kind":             children,
	}
	log.Debug("Ready to build query", zap.Any("bindVars", bindVars))

	filters := ""
	if ns := NSFilterValue(ctx); ns != "" {
		filters += fmt.Sprintf("FILTER path.vertices[-2]._key == \"%s\"\n", ns)
	}

	return db.Query(ctx, fmt.Sprintf(listObjectsOfKind, filters), bindVars)
}

const listOwnedQuery = `
FOR node, edge IN 0..100
OUTBOUND @from
GRAPH Permissions
OPTIONS { uniqueVertices: "path" }
FILTER !edge || edge.role == 1
    RETURN MERGE({ node: node._id }, edge ? { edge: edge._id, parent: edge._from } : { edge: null, parent: null })
`

func ListOwnedDeep(ctx context.Context, log *zap.Logger, db driver.Database, from InfinimeshGraphNode) (res *access.Nodes, err error) {
	c, err := db.Query(ctx, listOwnedQuery, map[string]interface{}{
		"from": from.ID(),
	})
	if err != nil {
		log.Debug("Error while executing query", zap.Error(err))
		return res, err
	}
	defer c.Close()

	var nodes []*access.Node
	for {
		var node access.Node
		_, err := c.ReadDocument(ctx, &node)
		if err != nil {
			if driver.IsNoMoreDocuments(err) {
				break
			}
			return res, err
		}
		nodes = append(nodes, &node)
	}

	return &access.Nodes{Nodes: nodes}, nil
}

func DeleteRecursive(ctx context.Context, log *zap.Logger, db driver.Database, from InfinimeshGraphNode) error {
	nodes, err := ListOwnedDeep(ctx, log, db, from)
	if err != nil {
		return err
	}

	cols := make(map[string]driver.Collection)
	for i := len(nodes.Nodes) - 1; i >= 0; i-- {
		node := nodes.Nodes[i]
		log.Debug("Deleting", zap.String("node", node.Node), zap.String("edge", node.Edge))

		if node.Node != "" {
			err := handleDeleteNodeInRecursion(ctx, log, db, node.Node, cols)
			if err != nil {
				if err.Error() == "ERR_ROOT_OBJECT_CANNOT_BE_DELETED" {
					continue
				}
				return err
			}
		}

		if node.Edge != "" {
			err := handleDeleteNodeInRecursion(ctx, log, db, node.Edge, cols)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SplitDocID(id string) (col, key string) {
	data := strings.SplitN(id, "/", 2)
	if len(data) != 2 {
		return "", ""
	}
	return data[0], data[1]
}

func handleDeleteNodeInRecursion(ctx context.Context, log *zap.Logger, db driver.Database, node string, cols map[string]driver.Collection) (err error) {
	log.Debug("Handling deletion", zap.String("node", node))

	id := strings.SplitN(node, "/", 2)
	log.Debug("Retrieving Collection", zap.String("collection", id[0]), zap.String("id", node))
	col, ok := cols[id[0]]
	if !ok {
		col, err = db.Collection(ctx, id[0])
		if err != nil {
			return err
		}
		cols[id[0]] = col
	}

	if id[0] == schema.ACCOUNTS_COL {
		if id[1] == schema.ROOT_ACCOUNT_KEY {
			log.Warn("Root account cannot be deleted")
			return errors.New("ERR_ROOT_OBJECT_CANNOT_BE_DELETED")
		}
		nodes, err := credentials.ListCredentialsAndEdges(ctx, log, col.Database(), driver.DocumentID(node))
		if err != nil {
			return err
		}
		for _, node := range nodes {
			err = handleDeleteNodeInRecursion(ctx, log, col.Database(), node, cols)
			if err != nil {
				return err
			}
		}
	}
	if id[0] == schema.NAMESPACES_COL && id[1] == schema.ROOT_NAMESPACE_KEY {
		log.Warn("Root namespace cannot be deleted")
		return errors.New("ERR_ROOT_OBJECT_CANNOT_BE_DELETED")
	}

	_, err = col.RemoveDocument(ctx, id[1])
	if e, ok := driver.AsArangoError(err); ok && e.Code == 404 {
		return nil
	}
	return err
}

const getWithAccessLevelQuery = `
FOR path IN OUTBOUND
K_SHORTEST_PATHS @requestor TO @node
GRAPH @permissions
    RETURN path.edges[-1].role == 2 ? path.edges[-1].level : path.edges[0].level
`

func AccessLevel(ctx context.Context, db driver.Database, requestor InfinimeshGraphNode, node InfinimeshGraphNode) (bool, access.Level) {
	if requestor.ID() == node.ID() {
		return true, access.Level_ROOT
	}
	c, err := db.Query(ctx, getWithAccessLevelQuery, map[string]interface{}{
		"requestor":   requestor.ID(),
		"node":        node.ID(),
		"permissions": schema.PERMISSIONS_GRAPH.Name,
	})
	if err != nil {
		return false, 0
	}
	defer c.Close()

	_access := access.Level_NONE
	for {
		var level access.Level
		_, err := c.ReadDocument(ctx, &level)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			continue
		}
		if level > _access {
			_access = level
		}
	}
	return _access > access.Level_NONE, _access
}

const toggleQuery = `
LET o = DOCUMENT(@node)
UPDATE o WITH {%[1]s: !o.%[1]s} IN @@col RETURN NEW 
`

func Toggle(ctx context.Context, db driver.Database, node InfinimeshGraphNode, field string) error {
	c, err := db.Query(ctx, fmt.Sprintf(toggleQuery, field), map[string]interface{}{
		"node": node.ID(),
		"@col": node.ID().Collection(),
	})
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.ReadDocument(ctx, &node)
	return err
}

func (r *infinimeshCommonActionsRepo) EnsureRootExists(_log *zap.Logger, rdb *redis.Client, passwd string) (err error) {

	ctx := context.TODO()
	log := _log.Named("EnsureRootExists")

	log.Debug("Checking Root Account exists")
	col, _ := r.db.Collection(ctx, schema.ACCOUNTS_COL)
	exists, err := col.DocumentExists(ctx, schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Warn("Error checking Root Account existance")
		return err
	}

	var meta driver.DocumentMeta
	if !exists {
		log.Debug("Root Account doesn't exist, creating")
		meta, err = col.CreateDocument(ctx, Account{
			Account: &accpb.Account{
				Title:   "infinimesh",
				Enabled: true,
			},
			DocumentMeta: driver.DocumentMeta{Key: schema.ROOT_ACCOUNT_KEY},
		})
		if err != nil {
			log.Warn("Error creating Root Account")
			return err
		}
		log.Debug("Created root Account", zap.Any("result", meta))
	}
	var acc accpb.Account
	meta, err = col.ReadDocument(ctx, schema.ROOT_ACCOUNT_KEY, &acc)
	if err != nil {
		log.Warn("Error reading Root Account")
		return err
	}
	root := &Account{
		Account:      &acc,
		DocumentMeta: meta,
	}

	ns_col, _ := r.db.Collection(ctx, schema.NAMESPACES_COL)
	exists, err = ns_col.DocumentExists(ctx, schema.ROOT_NAMESPACE_KEY)
	if err != nil || !exists {
		meta, err := ns_col.CreateDocument(ctx, Namespace{
			Namespace: &nspb.Namespace{
				Title: "infinimesh",
			},
			DocumentMeta: driver.DocumentMeta{Key: schema.ROOT_NAMESPACE_KEY},
		})
		if err != nil {
			log.Warn("Error creating Root Namespace")
			return err
		}
		log.Debug("Created root Namespace", zap.Any("result", meta))
	}

	var ns nspb.Namespace
	meta, err = ns_col.ReadDocument(ctx, schema.ROOT_NAMESPACE_KEY, &ns)
	if err != nil {
		log.Warn("Error reading Root Namespace")
		return err
	}
	rootNS := &Namespace{
		Namespace:    &ns,
		DocumentMeta: meta,
	}

	edge_col := r.GetEdgeCol(ctx, schema.ACC2NS)
	exists = r.CheckLink(ctx, edge_col, root, rootNS)
	if err != nil {
		log.Warn("Error checking link Root Account to Root Namespace", zap.Error(err))
		return err
	} else if !exists {
		err = r.Link(ctx, log, edge_col, root, rootNS, access.Level_ROOT, access.Role_OWNER)
		if err != nil {
			log.Warn("Error linking Root Account to Root Namespace")
			return err
		}
	}

	ctx = context.WithValue(ctx, schema.InfinimeshAccount, schema.ROOT_ACCOUNT_KEY)
	cred_edge_col, _ := r.db.Collection(ctx, schema.ACC2CRED)
	cred, err := credentials.NewStandardCredentials("infinimesh", passwd)
	if err != nil {
		log.Warn("Error creating Root Account Credentials")
		return err
	}

	ctrl := NewAccountsController(log, r.db, rdb)
	exists, err = cred_edge_col.DocumentExists(ctx, fmt.Sprintf("standard-%s", schema.ROOT_ACCOUNT_KEY))
	if err != nil || !exists {
		err = ctrl._SetCredentials(ctx, *root, cred_edge_col, cred)
		if err != nil {
			log.Warn("Error setting Root Account Credentials")
			return err
		}
	}
	_, res := ctrl.Authorize(ctx, "standard", "infinimesh", passwd)
	if !res {
		log.Warn("Error authorizing Root Account")
		return errors.New("cannot authorize infinimesh")
	}
	return nil
}
