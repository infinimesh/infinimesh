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
package schema

import (
	"context"
	"strings"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"go.uber.org/zap"
)

func CheckAndRegisterCollections(log *zap.Logger, db driver.Database, collections []string) {
	options := &driver.CreateCollectionOptions{
		KeyOptions: &driver.CollectionKeyOptions{AllowUserKeys: true, Type: "uuid"},
	}
	for _, col := range collections {
		log.Debug("Checking Collection existence", zap.String("collection", col))
		exists, err := db.CollectionExists(context.TODO(), col)
		if err != nil {
			log.Fatal("Failed to check collection", zap.Any(col, err))
		}
		log.Debug("Collection "+col, zap.Bool("Exists", exists))
		if !exists {
			log.Debug("Creating", zap.String("collection", col))
			_, err := db.CreateCollection(context.TODO(), col, options)
			if err != nil {
				log.Fatal("Failed to create collection", zap.Any(col, err))
			}
		}
	}
}

func CheckAndRegisterGraph(log *zap.Logger, db driver.Database, graph InfinimeshGraphSchema) {
	graphExists, err := db.GraphExists(context.TODO(), graph.Name)
	if err != nil {
		log.Fatal("Failed to check graph", zap.Any(graph.Name, err))
	}
	log.Debug("Graph Permissions", zap.Bool("Exists", graphExists))

	if graphExists {
		return
	}
	log.Debug("Creating", zap.String("graph", graph.Name))
	edges := make([]driver.EdgeDefinition, 0)
	for _, edge := range graph.Edges {
		edges = append(edges, driver.EdgeDefinition{
			Collection: strings.Join(edge, "2"),
			From:       []string{edge[0]}, To: []string{edge[1]},
		})
	}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = edges

	_, err = db.CreateGraph(context.TODO(), graph.Name, &options)
	if err != nil {
		log.Fatal("Failed to create Graph", zap.Any(graph.Name, err))
	}
}

func InitDB(log *zap.Logger, dbHost, dbCred, rootPass string, quick bool) driver.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://" + dbCred + "@" + dbHost},
	})
	if err != nil {
		log.Fatal("Error creating connection to DB", zap.Error(err))
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatal("Error creating driver instance for DB", zap.Error(err))
	}

	// Checking if DB exists and creating it if not
	log.Debug("Checking if DB exists")
	dbExists, err := c.DatabaseExists(context.TODO(), DB_NAME)
	if err != nil {
		log.Fatal("Error checking if DataBase exists", zap.Error(err))
	}
	log.Debug("DataBase", zap.Bool("Exists", dbExists))

	if dbExists && quick {
		return nil
	}

	var db driver.Database
	if !dbExists {
		_, err = c.CreateDatabase(context.TODO(), DB_NAME, nil)
		if err != nil {
			log.Fatal("Error creating DataBase", zap.Error(err))
		}
	}
	db, _ = c.Database(context.TODO(), DB_NAME)

	CheckAndRegisterCollections(log, db, COLLECTIONS)

	for _, graph := range GRAPHS_SCHEMAS {
		CheckAndRegisterGraph(log, db, graph)
	}

	return db
}
