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
	"time"

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/infinimesh/pkg/sessions"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	pb "github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	accpb "github.com/infinimesh/proto/node/accounts"
	"github.com/infinimesh/proto/node/namespaces"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Account struct {
	*accpb.Account
	driver.DocumentMeta
}

func (o *Account) ID() driver.DocumentID {
	return o.DocumentMeta.ID
}

func (o *Account) SetAccessLevel(level access.Level) {
	if o.Access == nil {
		o.Access = &access.Access{
			Level: level,
		}
		return
	}
	o.Access.Level = level
}

func (o *Account) GetAccess() *access.Access {
	if o.Access == nil {
		return &access.Access{
			Level: access.Level_NONE,
		}
	}
	return o.Access
}

func NewBlankAccountDocument(key string) *Account {
	return &Account{
		Account: &accpb.Account{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.ACCOUNTS_COL, key),
	}
}

func NewAccountFromPB(acc *accpb.Account) (res *Account) {
	return &Account{
		Account:      acc,
		DocumentMeta: NewBlankDocument(schema.ACCOUNTS_COL, acc.Uuid),
	}
}

type AccountsController struct {
	InfinimeshBaseController

	col  driver.Collection // Accounts Collection
	cred driver.Collection

	rdb *redis.Client

	acc2ns driver.Collection // Accounts to Namespaces permissions edge collection
	ns2acc driver.Collection // Namespaces to Accounts permissions edge collection

	sessions sessions.SessionsHandler

	ica_repo InfinimeshCommonActionsRepo // Infinimesh Common Actions Repository

	SIGNING_KEY []byte
}

func NewAccountsController(log *zap.Logger, db driver.Database, rdb *redis.Client) *AccountsController {
	ctx := context.TODO()
	perm_graph, _ := db.Graph(ctx, schema.PERMISSIONS_GRAPH.Name)
	col, _ := perm_graph.VertexCollection(ctx, schema.ACCOUNTS_COL)

	cred_graph, _ := db.Graph(ctx, schema.CREDENTIALS_GRAPH.Name)
	cred, _ := cred_graph.VertexCollection(ctx, schema.CREDENTIALS_COL)

	ica := NewInfinimeshCommonActionsRepo(db)

	return &AccountsController{
		InfinimeshBaseController: InfinimeshBaseController{
			log: log.Named("AccountsController"), db: db,
		}, col: col, cred: cred, rdb: rdb,

		acc2ns: ica.GetEdgeCol(ctx, schema.ACC2NS),
		ns2acc: ica.GetEdgeCol(ctx, schema.NS2ACC),

		sessions: sessions.NewSessionsHandlerModule(rdb).Handler(),

		ica_repo: ica,

		SIGNING_KEY: []byte("just-an-init-thing-replace-me"),
	}
}

func (c *AccountsController) Token(ctx context.Context, _req *connect.Request[pb.TokenRequest]) (*connect.Response[pb.TokenResponse], error) {
	log := c.log.Named("Token")
	req := _req.Msg
	log.Debug("Token request received", zap.Any("request", req))

	var account Account
	var ok bool

	if requestor := ctx.Value(inf.InfinimeshAccountCtxKey); requestor != nil && req.Uuid != nil {
		account = *NewBlankAccountDocument(*req.Uuid)
		requestor := requestor.(string)
		if *req.Uuid == requestor {
			return nil, status.Error(codes.PermissionDenied, "You can't create such token for yourself")
		}
		err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &account)
		if err != nil {
			log.Warn("Failed to get Account and access level", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "Wrong credentials given")
		}
		if account.Access.Level < access.Level_ROOT && account.Access.Role != access.Role_OWNER {
			log.Warn("Super-Admin Token Request attempted", zap.String("requestor", requestor), zap.String("account", account.Uuid))
			return nil, status.Error(codes.Unauthenticated, "Wrong credentials given")
		}

		req.Exp = time.Now().Unix() + int64(time.Minute.Seconds())*5
	} else {
		account, ok = c.Authorize(ctx, req.Auth.Type, req.Auth.Data...)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Wrong credentials given")
		}
	}

	log.Debug("Authorized user", zap.String("ID", account.ID().String()))
	if !account.Enabled {
		return nil, status.Error(codes.PermissionDenied, "Account is disabled")
	}

	session := c.sessions.New(req.Exp, req.GetClient())
	if err := c.sessions.Store(account.Key, session); err != nil {
		log.Error("Failed to store session", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to issue token: session")
	}

	claims := jwt.MapClaims{}
	claims[inf.INFINIMESH_ACCOUNT_CLAIM] = account.Key
	claims[inf.INFINIMESH_SESSION_CLAIM] = session.Id
	claims["exp"] = req.Exp

	if req.Inf != nil && *req.Inf {
		ok, lvl := c.ica_repo.AccessLevel(ctx, &account, NewBlankNamespaceDocument(schema.ROOT_NAMESPACE_KEY))
		claims[inf.INFINIMESH_ROOT_CLAIM] = ok && lvl > access.Level_ADMIN
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString(c.SIGNING_KEY)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to issue token")
	}

	return connect.NewResponse(&pb.TokenResponse{Token: token_string}), nil
}

func (c *AccountsController) Accessibles(ctx context.Context, req *connect.Request[namespaces.Namespace]) (*connect.Response[access.Nodes], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

func (c *AccountsController) Get(ctx context.Context, req *connect.Request[accpb.Account]) (res *connect.Response[accpb.Account], err error) {
	log := c.log.Named("Get")
	acc := req.Msg
	log.Debug("Get request received", zap.Any("request", acc))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	uuid := acc.GetUuid()
	if uuid == "me" {
		uuid = requestor
	}
	// Getting Account from DB
	// and Check requestor access
	result := *NewBlankAccountDocument(uuid)
	err = c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &result)
	if err != nil {
		log.Warn("Failed to get Account and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if result.Access.Level < access.Level_READ {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	return connect.NewResponse(result.Account), nil
}

func (c *AccountsController) List(ctx context.Context, _ *connect.Request[pb.EmptyMessage]) (*connect.Response[accpb.Accounts], error) {
	log := c.log.Named("List")

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	cr, err := c.ica_repo.ListQuery(ctx, log, NewBlankAccountDocument(requestor), schema.ACCOUNTS_COL)
	if err != nil {
		log.Warn("Error executing query", zap.Error(err))
		return nil, status.Error(codes.Internal, "Couldn't execute query")
	}
	defer cr.Close()

	var r []*accpb.Account
	for {
		var acc accpb.Account
		meta, err := cr.ReadDocument(ctx, &acc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Warn("Error unmarshalling Document", zap.Error(err))
			return nil, status.Error(codes.Internal, "Couldn't execute query")
		}
		acc.Uuid = meta.ID.Key()

		log.Debug("Got document", zap.Any("account", &acc))
		r = append(r, &acc)
	}

	return connect.NewResponse(&accpb.Accounts{
		Accounts: r,
	}), nil
}

func (c *AccountsController) Create(ctx context.Context, req *connect.Request[accpb.CreateRequest]) (*connect.Response[accpb.CreateResponse], error) {
	log := c.log.Named("Create")
	request := req.Msg
	log.Debug("Create request received", zap.Any("request", request), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	ns_id := request.GetNamespace()
	if ns_id == "" {
		ns_id = schema.ROOT_NAMESPACE_KEY
	}

	ok, level := c.ica_repo.AccessLevel(ctx, NewBlankAccountDocument(requestor), NewBlankNamespaceDocument(ns_id))
	if !ok || level < access.Level_ADMIN {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("no Access to Namespace %s", ns_id))
	}

	if request.Account.GetDefaultNamespace() == "" {
		request.Account.DefaultNamespace = ns_id
	}

	account := Account{Account: request.GetAccount()}
	meta, err := c.col.CreateDocument(ctx, account)
	if err != nil {
		log.Warn("Error creating Account", zap.Error(err))
		return nil, StatusFromString(connect.CodeInternal, "Error while creating Account")
	}
	account.Uuid = meta.ID.Key()
	account.DocumentMeta = meta

	ns := NewBlankNamespaceDocument(ns_id)
	err = c.ica_repo.Link(ctx, log, c.ns2acc, ns, &account, access.Level_ADMIN, access.Role_OWNER)
	if err != nil {
		defer c.col.RemoveDocument(ctx, meta.Key)
		log.Warn("Error Linking Namespace to Account", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	col, _ := c.db.Collection(ctx, schema.CREDENTIALS_EDGE_COL)
	cred, err := credentials.MakeCredentials(request.GetCredentials(), log)
	if err != nil {
		defer c.col.RemoveDocument(ctx, meta.Key)
		log.Warn("Error making Credentials for Account", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	err = c._SetCredentials(ctx, account, col, cred)
	if err != nil {
		defer c.col.RemoveDocument(ctx, meta.Key)
		log.Warn("Error setting Credentials for Account", zap.Error(err))
		return nil, err
	}
	return connect.NewResponse(
		&accpb.CreateResponse{Account: account.Account},
	), nil
}

func (c *AccountsController) Update(ctx context.Context, req *connect.Request[accpb.Account]) (*connect.Response[accpb.Account], error) {
	log := c.log.Named("Update")
	acc := req.Msg
	log.Debug("Update request received", zap.Any("request", acc), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))
	requestorAccount := NewBlankAccountDocument(requestor)

	old := *NewBlankAccountDocument(acc.GetUuid())
	err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, requestorAccount, &old)
	if err != nil || old.Access.Level < access.Level_ADMIN {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Account %s", acc.GetUuid())
	}

	if old.GetDefaultNamespace() != acc.GetDefaultNamespace() {
		ok, level := c.ica_repo.AccessLevel(ctx, &old, NewBlankNamespaceDocument(acc.GetDefaultNamespace()))
		if !ok || level < access.Level_READ {
			return nil, status.Errorf(codes.PermissionDenied, "Account has no Access to Namespace %s", acc.GetDefaultNamespace())
		}
	}

	_, err = c.col.UpdateDocument(ctx, acc.GetUuid(), acc)
	if err != nil {
		log.Warn("Internal error while updating Document", zap.Any("request", acc), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while updating Account")
	}

	return connect.NewResponse(acc), nil
}

func (c *AccountsController) Toggle(ctx context.Context, req *connect.Request[accpb.Account]) (*connect.Response[accpb.Account], error) {
	log := c.log.Named("Update")
	acc := req.Msg
	log.Debug("Update request received", zap.Any("account", acc), zap.Any("context", ctx))

	resp, err := c.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	curr := resp.Msg

	if curr.Access.Level < access.Level_MGMT {
		return nil, status.Errorf(codes.PermissionDenied, "No Access to Account %s", acc.Uuid)
	}

	res := NewAccountFromPB(curr)
	err = c.ica_repo.Toggle(ctx, res, "enabled")
	if err != nil {
		log.Warn("Error updating Account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while updating Account")
	}

	return connect.NewResponse(res.Account), nil
}

func (c *AccountsController) Deletables(ctx context.Context, req *connect.Request[accpb.Account]) (*connect.Response[access.Nodes], error) {
	log := c.log.Named("Deletables")
	request := req.Msg
	log.Debug("Deletables request received", zap.Any("request", request))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(request.GetUuid())
	err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &acc)
	if err != nil {
		log.Warn("Error getting Account and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if acc.Access.Role != access.Role_OWNER && acc.Access.Level < access.Level_ROOT {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	nodes, err := c.ica_repo.ListOwnedDeep(ctx, log, &acc)
	if err != nil {
		log.Warn("Error getting owned nodes", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting owned nodes")
	}

	return connect.NewResponse(nodes), nil
}

func (c *AccountsController) Delete(ctx context.Context, request *connect.Request[accpb.Account]) (*connect.Response[pb.DeleteResponse], error) {
	log := c.log.Named("Delete")
	req := request.Msg
	log.Debug("Delete request received", zap.Any("request", req), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(req.GetUuid())
	err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &acc)
	if err != nil {
		log.Warn("Error getting Account and access level", zap.Error(err))
		return nil, status.Error(codes.NotFound, "Account not found or not enough Access Rights")
	}
	if acc.Access.Role != access.Role_OWNER && acc.Access.Level < access.Level_ADMIN {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access Rights")
	}

	err = c.ica_repo.DeleteRecursive(ctx, log, &acc)
	if err != nil {
		log.Warn("Error deleting account", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error deleting account")
	}

	return connect.NewResponse(&pb.DeleteResponse{}), nil
}

// Helper Functions

func (ctrl *AccountsController) Authorize(ctx context.Context, auth_type string, args ...string) (Account, bool) {
	ctrl.log.Debug("Authorization request", zap.String("type", auth_type))

	credentials, err := credentials.Find(ctx, ctrl.col.Database(), ctrl.log, auth_type, args...)
	// Check if could authorize
	if err != nil {
		ctrl.log.Info("Coudn't authorize", zap.Error(err))
		return Account{}, false
	}

	account, ok := Authorisable(ctx, &credentials, ctrl.col.Database())
	ctrl.log.Debug("Authorized account", zap.Bool("result", ok), zap.Any("account", account))
	return account, ok
}

// Return Account authorisable by this Credentials
func Authorisable(ctx context.Context, cred *credentials.Credentials, db driver.Database) (Account, bool) {
	query := `FOR account IN 1 INBOUND @credentials GRAPH @credentials_graph RETURN account`
	c, err := db.Query(ctx, query, map[string]interface{}{
		"credentials":       cred,
		"credentials_graph": schema.CREDENTIALS_GRAPH.Name,
	})
	if err != nil {
		return Account{}, false
	}
	defer c.Close()

	var r Account
	_, err = c.ReadDocument(ctx, &r)
	return r, err == nil
}

func (c *AccountsController) GetCredentials(ctx context.Context, request *connect.Request[pb.GetCredentialsRequest]) (*connect.Response[pb.GetCredentialsResponse], error) {
	log := c.log.Named("GetCredentials")
	req := request.Msg
	log.Debug("Get Credentials request received", zap.String("account", req.GetUuid()))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(req.GetUuid())
	err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &acc)

	if err != nil {
		log.Warn("Error getting Account", zap.String("requestor", requestor), zap.String("account", req.GetUuid()), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting Account or not enough Access right to set credentials for this Account")
	}

	if acc.Access.Level < access.Level_ROOT && acc.Access.Role != access.Role_OWNER {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access right to set credentials for this Account. Only Owner and Super-Admin can do this")
	}

	linked, err := credentials.ListCredentials(ctx, log, c.db, acc.ID())
	if err != nil {
		return nil, status.Error(codes.Internal, "Error listing Account's Credentials")
	}

	var creds []*accpb.Credentials
	for _, res := range linked {
		listable, err := credentials.MakeListable(res)
		if err != nil {
			log.Warn("Couldn't make Listable", zap.Error(err))
			continue
		}
		creds = append(creds, &accpb.Credentials{
			Type: res.Type, Data: listable.Listable(),
		})
	}

	return connect.NewResponse(&pb.GetCredentialsResponse{Credentials: creds}), nil
}

// Set Account Credentials, ensure account has only one credentials document linked per credentials type
func (ctrl *AccountsController) _SetCredentials(ctx context.Context, acc Account, edge driver.Collection, c credentials.Credentials) error {
	key := c.Type() + "-" + acc.Key
	var oldLink credentials.Link
	meta, err := edge.ReadDocument(ctx, key, &oldLink)
	if err == nil {
		ctrl.log.Debug("Link exists", zap.Any("meta", meta))
		_, err = ctrl.cred.UpdateDocument(ctx, oldLink.To.Key(), c)
		if err != nil {
			ctrl.log.Warn("Error updating Credentials of type", zap.Error(err), zap.String("key", key))
			return status.Error(codes.InvalidArgument, "Error updating Credentials of type")
		}

		return nil
	}
	ctrl.log.Debug("Credentials either not created yet or failed to get them from DB, overwriting", zap.Error(err), zap.String("key", key))

	cred, err := ctrl.cred.CreateDocument(ctx, c)
	if err != nil {
		ctrl.log.Warn("Error creating Credentials Document", zap.String("type", c.Type()), zap.Error(err))
		return status.Error(codes.Internal, "Couldn't create credentials")
	}

	_, err = edge.CreateDocument(ctx, credentials.Link{
		From: acc.ID(),
		To:   cred.ID,
		Type: c.Type(),
		DocumentMeta: driver.DocumentMeta{
			Key: key, // Ensures only one credentials vertex per type
		},
	})
	if err != nil {
		ctrl.log.Warn("Error Linking Credentials to Account",
			zap.String("account", acc.Key), zap.String("type", c.Type()), zap.Error(err),
		)
		ctrl.cred.RemoveDocument(ctx, cred.Key)
		return status.Error(codes.Internal, "Couldn't assign credentials")
	}
	return nil
}

func (c *AccountsController) SetCredentials(ctx context.Context, _req *connect.Request[pb.SetCredentialsRequest]) (*connect.Response[pb.SetCredentialsResponse], error) {
	log := c.log.Named("SetCredentials")
	req := _req.Msg
	log.Debug("Set Credentials request received", zap.String("account", req.GetUuid()), zap.String("type", req.GetCredentials().GetType()), zap.Any("context", ctx))

	requestor := ctx.Value(inf.InfinimeshAccountCtxKey).(string)
	log.Debug("Requestor", zap.String("id", requestor))

	acc := *NewBlankAccountDocument(req.GetUuid())
	err := c.ica_repo.AccessLevelAndGet(ctx, log, c.db, NewBlankAccountDocument(requestor), &acc)

	if err != nil {
		log.Warn("Error getting Account", zap.String("requestor", requestor), zap.String("account", req.GetUuid()), zap.Error(err))
		return nil, status.Error(codes.Internal, "Error getting Account or not enough Access right to set credentials for this Account")
	}

	if acc.Access.Level < access.Level_ROOT && acc.Access.Role != access.Role_OWNER {
		return nil, status.Error(codes.PermissionDenied, "Not enough Access right to set credentials for this Account. Only Owner and Super-Admin can do this")
	}

	col, _ := c.db.Collection(ctx, schema.CREDENTIALS_EDGE_COL)
	cred, err := credentials.MakeCredentials(req.GetCredentials(), log)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = c._SetCredentials(ctx, acc, col, cred)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&pb.SetCredentialsResponse{}), nil
}

func (c *AccountsController) DelCredentials(context.Context, *connect.Request[pb.DeleteCredentialsRequest]) (*connect.Response[pb.DeleteResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
