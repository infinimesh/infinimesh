/*
Copyright © 2018-2024 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

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

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	proto_eventbus "github.com/infinimesh/proto/eventbus"

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

	ica  InfinimeshCommonActionsRepo
	repo InfinimeshGenericActionsRepo[*nspb.Namespace]

	bus EventBusService

	db driver.Database
}

func NewNamespacesController(log *zap.Logger, db driver.Database, bus EventBusService, ica InfinimeshCommonActionsRepo, repo InfinimeshGenericActionsRepo[*nspb.Namespace]) *NamespacesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.NAMESPACES_COL)
	accs, _ := db.Collection(ctx, schema.ACCOUNTS_COL)

	return &NamespacesController{
		log: log.Named("NamespacesController"), col: col, db: db, accs: accs,
		acc2ns: ica.GetEdgeCol(ctx, schema.ACC2NS), ns2acc: ica.GetEdgeCol(ctx, schema.NS2ACC),
		ica:  ica,
		repo: repo,
		bus:  bus,
	}
}

func (c *NamespacesController) Create(ctx context.Context, req *connect.Request[nspb.Namespace]) (*connect.Response[nspb.Namespace], error) {
	log := c.log.Named("Create")
	request := req.Msg
	log.Debug("Create request received", zap.Any("request", request), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	if request.Title == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("Title is required"))
	}

	if request.Uuid != "" {
		request.Uuid = ""
	}

	namespace := Namespace{Namespace: request}
	meta, err := c.col.CreateDocument(ctx, namespace)
	if err != nil {
		log.Warn("Error creating namespace", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error while creating namespace"))
	}
	namespace.Uuid = meta.ID.Key()
	namespace.DocumentMeta = meta

	requestorAcc := NewBlankAccountDocument(requestor)
	err = c.ica.Link(ctx, c.acc2ns,
		requestorAcc,
		&namespace, access.Level_ADMIN, access.Role_OWNER,
	)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, namespace.Uuid)
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error creating Permission"))
	}

	notifier, err := c.bus.Notify(ctx, &proto_eventbus.Event{
		EventKind: proto_eventbus.EventKind_NAMESPACE_CREATE,
		Entity:    &proto_eventbus.Event_Namespace{Namespace: namespace.Namespace},
	})

	if err == nil {
		err = notifier()
		if err != nil {
			log.Warn("Failed to notify", zap.Error(err))
		}
	} else {
		log.Warn("Failed to create notifier", zap.Error(err))
	}

	return connect.NewResponse(namespace.Namespace), nil
}

func (c *NamespacesController) Get(ctx context.Context, ns *connect.Request[nspb.Namespace]) (res *connect.Response[nspb.Namespace], err error) {
	log := c.log.Named("Get")
	log.Debug("Get request received", zap.Any("request", ns))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	uuid := ns.Msg.GetUuid()
	// Getting Namespace from DB
	// and Check requestor access
	result := *NewBlankNamespaceDocument(uuid)
	err = c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &result)
	if err != nil {
		log.Warn("Failed to get Namespace and access level", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found or not enough Access Rights"))
	}
	if result.Access.Level < access.Level_READ {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	return connect.NewResponse(result.Namespace), nil
}

func (c *NamespacesController) Update(ctx context.Context, req *connect.Request[nspb.Namespace]) (res *connect.Response[nspb.Namespace], err error) {
	log := c.log.Named("Update")
	ns := req.Msg
	log.Debug("Request received", zap.Any("namespace", ns))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	curr := *NewBlankNamespaceDocument(ns.Uuid)
	err = c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &curr)

	if err != nil {
		log.Warn("Can't get Namespaces from DB", zap.String("namespace", ns.Uuid), zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found"))
	}

	if curr.Access.Level < 3 {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	changed := false
	if ns.Title != "" {
		curr.Title = ns.Title
		changed = true
	}

	if ns.Plugin != nil {
		curr.Plugin = ns.Plugin
		changed = true
	}

	if ns.Config != nil {
		curr.Config = ns.Config
		changed = true
	}

	if changed {
		_, err := c.col.ReplaceDocument(ctx, curr.Uuid, curr)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("Error while updating Namespace in DB: %v", err))
		}

		notifier, err := c.bus.Notify(ctx, &proto_eventbus.Event{
			EventKind: proto_eventbus.EventKind_NAMESPACE_UPDATE,
			Entity:    &proto_eventbus.Event_Namespace{Namespace: ns},
		})

		if err == nil {
			err = notifier()
			if err != nil {
				log.Warn("Failed to notify", zap.Error(err))
			}
		} else {
			log.Warn("Failed to create notifier", zap.Error(err))
		}
	}

	return connect.NewResponse(curr.Namespace), nil
}

func (c *NamespacesController) List(ctx context.Context, _ *connect.Request[pb.EmptyMessage]) (*connect.Response[nspb.Namespaces], error) {
	log := c.log.Named("List")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	result, err := c.repo.ListQuery(ctx, log, NewBlankAccountDocument(requestor))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&nspb.Namespaces{
		Namespaces: result.Result,
	}), nil
}

const listJoinsQuery = `
FOR node, edge, path IN 1 INBOUND @namespace
GRAPH Permissions
FILTER edge.role != 1 && edge.level > 0
RETURN MERGE(node, { uuid: node._key, access: { level: edge.level } })
`

func (c *NamespacesController) Joins(ctx context.Context, req *connect.Request[nspb.Namespace]) (*connect.Response[accpb.Accounts], error) {
	log := c.log.Named("Joins")
	request := req.Msg

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetUuid())
	err := c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found or not enough Access Rights"))
	}

	if ns.Access.Level < access.Level_ADMIN {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	cr, err := c.db.Query(ctx, listJoinsQuery, map[string]interface{}{
		"namespace": ns.ID(),
	})
	if err != nil {
		log.Warn("Error querying for joins", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error querying for joins"))
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
			return nil, connect.NewError(connect.CodeInternal, errors.New("Couldn't execute query"))
		}
		log.Debug("Got document", zap.Any("account", &acc))
		r = append(r, &acc)
	}

	return connect.NewResponse(&accpb.Accounts{Accounts: r}), nil
}

func (c *NamespacesController) Join(ctx context.Context, req *connect.Request[pb.JoinRequest]) (*connect.Response[accpb.Accounts], error) {
	log := c.log.Named("Join")
	request := req.Msg

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.GetNamespace())
	err := c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found or not enough Access Rights"))
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	acc := *NewBlankAccountDocument(request.GetAccount())
	_, err = c.accs.ReadDocument(ctx, request.GetAccount(), &acc)
	if err != nil {
		log.Warn("Error getting Account", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Account not found"))
	}

	if request.Access > ns.Access.Level {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights: can't grant higher access than current"))
	}

	err = c.ica.Link(ctx, c.acc2ns, &acc, &ns, request.Access, access.Role_UNSET)
	if err != nil {
		log.Warn("Error creating edge", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error creating Permission"))
	}

	return c.Joins(ctx, connect.NewRequest(ns.Namespace))
}

func (c *NamespacesController) Deletables(ctx context.Context, request *connect.Request[nspb.Namespace]) (*connect.Response[access.Nodes], error) {
	log := c.log.Named("Deletables")
	log.Debug("Deletables request received", zap.Any("request", request))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.Msg.GetUuid())
	err := c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found or not enough Access Rights"))
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	nodes, err := c.ica.ListOwnedDeep(ctx, &ns)
	if err != nil {
		log.Warn("Error getting owned nodes", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error getting owned nodes"))
	}

	return connect.NewResponse(nodes), nil
}

func (c *NamespacesController) Delete(ctx context.Context, request *connect.Request[nspb.Namespace]) (*connect.Response[pb.DeleteResponse], error) {
	log := c.log.Named("Delete")
	log.Debug("Delete request received", zap.Any("request", request), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns := *NewBlankNamespaceDocument(request.Msg.GetUuid())
	err := c.ica.AccessLevelAndGet(ctx, NewBlankAccountDocument(requestor), &ns)
	if err != nil {
		log.Warn("Error getting Namespace and access level", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, errors.New("Namespace not found or not enough Access Rights"))
	}
	if ns.Access.Role != access.Role_OWNER && ns.Access.Level < access.Level_ROOT {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("Not enough Access Rights"))
	}

	notifier, notify_err := c.bus.Notify(ctx, &proto_eventbus.Event{
		EventKind: proto_eventbus.EventKind_NAMESPACE_DELETE,
		Entity:    &proto_eventbus.Event_Namespace{Namespace: ns.Namespace},
	})

	err = c.ica.DeleteRecursive(ctx, &ns)
	if err != nil {
		log.Warn("Error deleting namespace", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("Error deleting namespace"))
	}

	if notify_err == nil {
		err = notifier()
		if err != nil {
			log.Warn("Failed to notify", zap.Error(err))
		}
	} else {
		log.Warn("Failed to create notifier", zap.Error(notify_err))
	}

	return connect.NewResponse(&pb.DeleteResponse{}), nil
}

func (c *NamespacesController) Accessibles(context.Context, *connect.Request[nspb.Namespace]) (*connect.Response[access.Nodes], error) {
	return nil, StatusFromString(connect.CodeUnimplemented, "Not implemented")
}
