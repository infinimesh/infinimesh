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
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	"go.uber.org/zap"

	// accpb "github.com/infinimesh/infinimesh/pkg/node/proto/accounts"
	nspb "github.com/infinimesh/infinimesh/pkg/node/proto/namespaces"
)
type Namespace struct {
	*nspb.Namespace
	driver.DocumentMeta
}

func (o *Namespace) ID() (driver.DocumentID) {
	return o.DocumentMeta.ID
}

type NamespacesController struct {
	pb.UnimplementedNamespacesServiceServer
	log *zap.Logger

	col driver.Collection // Namespaces Collection
	acc2ns driver.Collection // Accounts to Namespaces permissions edge collection
	ns2acc driver.Collection // Namespaces to Accounts permissions edge collection

	db driver.Database
	SIGNING_KEY []byte
}

func NewNamespacesController(log *zap.Logger, db driver.Database) NamespacesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.NAMESPACES_COL)
	return NamespacesController{
		log: log.Named("NamespacesController"), col: col, db: db,
		acc2ns: GetEdgeCol(ctx, db, schema.ACC2NS), ns2acc: GetEdgeCol(ctx, db, schema.NS2ACC),
		SIGNING_KEY: []byte("just-an-init-thing-replace-me")}
}

func (c *NamespacesController) Create(ctx context.Context, request *nspb.Namespace) (*nspb.Namespace, error) {
	log := c.log.Named("Create")
	log.Debug("Create request received", zap.Any("request", request), zap.Any("context", ctx))

	//Get metadata from context and perform validation
	_, requestor, err := Validate(ctx, log)
	if err != nil {
		return nil, err
	}
	log.Debug("Requestor", zap.String("id", requestor))

	namespace := Namespace{Namespace: request}
	meta, err := c.col.CreateDocument(ctx, namespace)
	if err != nil {
		c.log.Debug("Error creating namespace", zap.Error(err))
		return nil, status.Error(codes.Internal, "Error while creating namespace")
	}
	namespace.Uuid = meta.ID.Key()
	namespace.DocumentMeta = meta


	err = Link(ctx, c.acc2ns,
		NewBlankAccountDocument(requestor),
		&namespace, schema.ADMIN,
	)
	if err != nil {
		log.Error("Error creating edge", zap.Error(err))
		c.col.RemoveDocument(ctx, namespace.GetUuid())
		return nil, status.Error(codes.Internal, "error creating Permission")
	}

	return namespace.Namespace, nil
}

