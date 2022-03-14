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
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"go.uber.org/zap"
)

type Access struct {
	From driver.DocumentID `json:"_from"`
	To driver.DocumentID `json:"_to"`
	Level schema.InfinimeshAccessLevel `json:"level"`

	driver.DocumentMeta
}

type InfinimeshGraphNode interface {
	GetUuid() string
	ID() driver.DocumentID
}

func NewBlankDocument(col string, key string) (driver.DocumentMeta) {
	return driver.DocumentMeta{
		Key: key,
		ID: driver.NewDocumentID(col, key),
	}
}

func GetEdgeCol(ctx context.Context, db driver.Database, name string) (driver.Collection) {
	col, _ := db.Collection(ctx, name)
	return col
}

func Link(ctx context.Context, log *zap.Logger, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode, access schema.InfinimeshAccessLevel) error {
	log.Debug("Linking two nodes",
		zap.Any("from", from),
		zap.Any("to", to),
	)
	_, err := edge.CreateDocument(ctx, Access{
		From: from.ID(),
		To: to.ID(),
		Level: access,
		DocumentMeta: driver.DocumentMeta {
			Key: from.GetUuid() + "-" + to.GetUuid(),
		},
	})
	return err
}

func CheckLink(ctx context.Context, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode) (bool) {
	r, err := edge.DocumentExists(ctx, from.GetUuid() + "-" + to.GetUuid())
	return err == nil && r
}