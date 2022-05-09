/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

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
package credentials

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

var (
	log *zap.Logger
	arangodbHost string
	arangodbCred string

	rootCtx context.Context

	db driver.Database
)

func init() {
	viper.AutomaticEnv()
	log = zap.NewExample()

	viper.SetDefault("DB_HOST", "db.infinimesh.local")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("INF_DEFAULT_ROOT_PASS", "infinimesh")

	arangodbHost = viper.GetString("DB_HOST")
	arangodbCred = viper.GetString("DB_CRED")
	db = schema.InitDB(log, arangodbHost, arangodbCred, "infinimesh", false)
	
	md := metadata.New(map[string]string{
		inf.INFINIMESH_ACCOUNT_CLAIM: schema.ROOT_ACCOUNT_KEY,
	})
	rootCtx = metadata.NewIncomingContext(context.Background(), md)
	rootCtx = context.WithValue(rootCtx, inf.InfinimeshAccountCtxKey, schema.ROOT_ACCOUNT_KEY)
}

// TODO: Automate this test
// func TestListCredentialsAndEdges(t *testing.T) {
// 	nodes, err := ListCredentialsAndEdges(rootCtx, log, db, driver.NewDocumentID("Accounts", "infinimesh"))
// 	if err != nil {
// 		t.Errorf("Failed to list credentials and edges: %v", err)
// 	}

// 	log.Info("Retrieved nodes", zap.Any("nodes", nodes))
// }