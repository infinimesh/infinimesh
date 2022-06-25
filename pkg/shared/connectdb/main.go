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
package connectdb

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"go.uber.org/zap"
)

func MakeDBConnection(log *zap.Logger, host, cred string) driver.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://" + cred + "@" + host},
	})
	if err != nil {
		log.Fatal("Error creating connection to DB", zap.Error(err))
	}
	log.Debug("Instantiated DB connection", zap.Any("conn", conn))

	log.Info("Setting up DB client")
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatal("Error creating driver instance for DB", zap.Error(err))
	}
	log.Debug("Instantiated DB client", zap.Any("client", c))

	db_connect_attempts := 0
db_connect:
	log.Info("Trying to connect to DB")
	db, err := c.Database(context.TODO(), schema.DB_NAME)
	if err != nil {
		db_connect_attempts++
		log.Error("Failed to connect DB", zap.Error(err), zap.Int("attempts", db_connect_attempts), zap.Int("next_attempt", db_connect_attempts*5))
		time.Sleep(time.Duration(db_connect_attempts*5) * time.Second)
		goto db_connect
	}
	return db
}
