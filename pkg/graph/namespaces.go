/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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

	"github.com/arangodb/go-driver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	pb "github.com/infinimesh/proto/node"
	accpb "github.com/infinimesh/proto/node/accounts"
	"go.uber.org/zap"

	"github.com/infinimesh/proto/node/access"
	nspb "github.com/infinimesh/proto/node/namespaces"
)

type Namespace struct {
	*nspb.Namespace
	driver.DocumentMeta
}

func (o *Namespace) ID() driver.DocumentID {
	return o.DocumentMeta.ID
}

func (o *Namespace) SetAccessLevel(level access.Level) {
	if o.Access == nil {
		o.Access = &access.Access{
			Level: level,
		}
		return
	}
	o.Access.Level = level
}

func NewBlankNamespaceDocument(key string) *Namespace {
	return &Namespace{
		Namespace: &nspb.Namespace{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.NAMESPACES_COL, key),
	}
}

type NamespacesController struct {
	pb.UnimplementedNamespacesServiceServer
	log *zap.Logger

	col    driver.Collection // Namespaces Collection
	accs   driver.Collection // Accounts Collection
	acc2ns driver.Collection // Accounts to Namespaces permissions edge collection
	ns2acc driver.Collection // Namespaces to Accounts permissions edge collection

	db driver.Database
}

func NewNamespacesController(log *zap.Logger, db driver.Database) *NamespacesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.NAMESPACES_COL)
	accs, _ := db.Collection(ctx, schema.ACCOUNTS_COL)
	return &NamespacesController{
		log: log.Named("NamespacesController"), col: col, db: db, accs: accs,
		acc2ns: GetEdgeCol(ctx, db, schema.ACC2NS), ns2acc: GetEdgeCol(ctx, db, schema.NS2ACC),
	}
}

func (c *NamespacesController) Create(ctx context.Context, request *nspb.Namespace) (*nspb.Namespace, error) {
	log := c.log.Named("Create")
	log.Debug("Create request received", zap.Any("request", request), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	if request.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "Title is required")
	}

	if request.Uuid != "" {
		request.Uuid = ""
	}

	namespace := Namespace{Namespace: request}
	meta, err := c.col.CreateDocument(ctx, namespace)
	if err != nil {
		log.Warn("Error creating namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while creating namespace")
	}
	namespace.Uuid = meta.ID.Key()
	namespace.DocumentMeta = meta

	requestorAcc := NewBlankAccountDocument(requestor)
	err = Link(ctx, log, c.acc2ns,
		requestorAcc,
		&namespace, access.Level_ADMIN, access.Role_OWNER,
	)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, namespace.Uuid)
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return namespace.Namespace, nil
}

func (c *NamespacesController) Get(ctx context.Context, ns *nspb.Namespace) (res *nspb.Namespace, err error) {
	log := c.log.Named("Get")
	log.Debug("Get request received", zap.Any("request", ns))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	uuid := ns.GetUuid()
	// Getting Namespace from DB
	// and Check requestor access
	result := *NewBlankNamespaceDocument(uuid)
	err = AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &result)
	if err != nil {
		log.Warn("Failed to get Namespace and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if result.Access.Level < access.Level_READ {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	return result.Namespace, nil
}

func (c *NamespacesController) List(ctx context.Context, _ *pb.EmptyMessage) (*nspb.Namespaces, error) {
	log := c.log.Named("List")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := ListQuery(ctx, log, c.db, NewBlankAccountDocument(requestor), schema.NAMESPACES_COL, 4)
	if err != nil {
		return nil, err
	}
	defer cr.Close()

	var r []*nspb.Namespace
	for {
		var ns nspb.Namespace
		_, err := cr.ReadDocument(ctx, &ns)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		log.Debug("Got document", zap.Any("namespace", &ns))
		r = append(r, &ns)
	}

	return &nspb.Namespaces{
		Namespaces: r,
	}, nil
}

const listJoinsQuery = `
FOR node, edge, path IN 1 INBOUND @namespace
GRAPH Permissions
FILTER edge.role != 1 && edge.level > 0
RETURN MERGE(node, { uuid: node._key, access: { level: edge.level } })
`

func (c *NamespacesController) Joins(ctx context.Context, request *nspb.Namespace) (*accpb.Accounts, error) {
	log := c.log.Named("Joins")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetUuid())
	err := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}

	if ns.Access.Level < access.Level_ADMIN {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	cr, err := c.db.Query(ctx, listJoinsQuery, map[string]interface{}{
		"namespace": ns.ID(),
	})
	if err != nil {
		log.Warn("Error querying for joins", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error querying for joins")
	}
	defer cr.Close()

	var r []*accpb.Account
	for {
		var acc accpb.Account
		_, err := cr.ReadDocument(ctx, &acc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Warn("Error unmarshalling Document", zap.Error(err))
			return nil, status.Error(codes.Internal, "Couldn't execute query")
		}
		log.Debug("Got document", zap.Any("account", &acc))
		r = append(r, &acc)
	}

	return &accpb.Accounts{Accounts: r}, nil
}

func (c *NamespacesController) Join(ctx context.Context, request *pb.JoinRequest) (*accpb.Accounts, error) {
	log := c.log.Named("Join")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetNamespace())
	err := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	acc := *NewBlankAccountDocument(request.GetAccount())
	_, err = c.accs.ReadDocument(ctx, request.GetAccount(), &acc)
	if err != nil {
		log.Warn("Error getting Account", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Account not found")
	}

	if request.Access > ns.Access.Level {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights: can't grant higher access than current")
	}

	err = Link(ctx, log, c.acc2ns, &acc, &ns, request.Access, access.Role_UNSET)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return c.Joins(ctx, ns.Namespace)
}

func (c *NamespacesController) Deletables(ctx context.Context, request *nspb.Namespace) (*access.Nodes, error) {
	log := c.log.Named("Deletables")
	log.Debug("Deletables request received", zap.Any("request", request))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetUuid())
	err := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	nodes, err := ListOwnedDeep(ctx, log, c.db, &ns)
	if err != nil {
		log.Warn("Error getting owned nodes", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting owned nodes")
	}

	return nodes, nil
}

func (c *NamespacesController) Delete(ctx context.Context, request *nspb.Namespace) (*pb.DeleteResponse, error) {
	log := c.log.Named("Delete")
	log.Debug("Delete request received", zap.Any("request", request), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetUuid())
	err := AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	err = DeleteRecursive(ctx, log, c.db, &ns)
	if err != nil {
		log.Warn("Error deleting namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error deleting namespace")
	}

	return &pb.DeleteResponse{}, nil
}
