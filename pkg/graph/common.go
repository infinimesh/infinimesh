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
	SetAccessLevel(level schema.InfinimeshAccessLevel)
}

func NewBlankDocument(col string, key string) (driver.DocumentMeta) {
	return driver.DocumentMeta{
		Key: key,
		ID: driver.NewDocumentID(col, key),
	}
}

func GetEdgeCol(ctx context.Context, db driver.Database, name string) (driver.Collection) {
	g, _ := db.Graph(ctx, schema.PERMISSIONS_GRAPH.Name)
	col, _, _ := g.EdgeCollection(ctx, name)
	return col
}

func Link(ctx context.Context, log *zap.Logger, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode, access schema.InfinimeshAccessLevel) error {
	log.Debug("Linking two nodes",
		zap.Any("from", from.ID()),
		zap.Any("to", to.ID()),
	)
	_, err := edge.CreateDocument(ctx, Access{
		From: from.ID(),
		To: to.ID(),
		Level: access,
		DocumentMeta: driver.DocumentMeta {
			Key: from.ID().Key() + "-" + to.ID().Key(),
		},
	})
	return err
}

func CheckLink(ctx context.Context, edge driver.Collection, from InfinimeshGraphNode, to InfinimeshGraphNode) (bool) {
	r, err := edge.DocumentExists(ctx, from.ID().Key() + "-" + to.ID().Key())
	return err == nil && r
}

const getWithAccessLevelAndNS = `
FOR path IN OUTBOUND K_SHORTEST_PATHS @account TO @node
GRAPH @permissions SORT path.edges[0].level
	RETURN MERGE(path.vertices[-1], { access_level: path.edges[0].level, namespace: path.vertices[-2]._key })
`
func AccessLevelAndGet(ctx context.Context, log *zap.Logger, db driver.Database, account *Account, node InfinimeshGraphNode) (error) {
	vars :=  map[string]interface{}{
		"account": account.ID(),
		"node": node.ID(),
		"permissions": schema.PERMISSIONS_GRAPH.Name,
	}
	c, err := db.Query(ctx, getWithAccessLevelAndNS, vars)
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

	if account.ID() == node.ID() {
		node.SetAccessLevel(schema.ROOT)
	}

	return nil
}

// List children nodes
// ctx - context
// log - logger
// db - Database connection
// from - Graph node to start traversal from
// children - children type(collection name)
// depth
func ListQuery(ctx context.Context, log *zap.Logger, db driver.Database, from InfinimeshGraphNode, children string, depth int) (driver.Cursor, error) {
	query := `
FOR node, edge, path IN 0..@depth OUTBOUND @from
	GRAPH @permissions_graph
	OPTIONS {order: "bfs", uniqueVertices: "global"}
	FILTER IS_SAME_COLLECTION(@@kind, node)
	FILTER edge.level > 0
	RETURN MERGE(node, { access_level: path.edges[0].level, namespace: path.vertices[-2]._key })`
	bindVars := map[string]interface{}{
		"depth": depth,
		"from": from.ID(),
		"permissions_graph": schema.PERMISSIONS_GRAPH.Name,
		"@kind": children,
	}

	log.Debug("Ready to build query", zap.Any("bindVars", bindVars))
	return db.Query(ctx, query, bindVars)
}


func AccessLevel(ctx context.Context, db driver.Database, account *Account, node InfinimeshGraphNode) (bool, int32) {
	if account.ID() == node.ID() {
		return true, int32(schema.ROOT)
	}
	query := `FOR path IN OUTBOUND K_SHORTEST_PATHS @account TO @node GRAPH @permissions RETURN path.edges[0].level`
	c, err := db.Query(ctx, query, map[string]interface{}{
		"account": account.ID(),
		"node": node.ID(),
		"permissions": schema.PERMISSIONS_GRAPH.Name,
	})
	if err != nil {
		return false, 0
	}
	defer c.Close()

	var access int32 = 0
	for {
		var level int32
		_, err := c.ReadDocument(ctx, &level)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			continue
		}
		if level > access {
			access = level
		}
	}
	return access > 0, access
}