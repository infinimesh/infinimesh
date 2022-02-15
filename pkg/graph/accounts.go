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
package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	pb "github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"go.uber.org/zap"
)

type Account struct {
	pb.Account
	driver.DocumentMeta
}

type Credentials struct {
	
}

type AccountsController struct {
	pb.UnimplementedAccountServiceServer
	log *zap.Logger

	col driver.Collection // Accounts Collection
	db driver.Database
}

func NewAccountsController(log *zap.Logger, db driver.Database) AccountsController {
	col, _ := db.Collection(context.TODO(), schema.ACCOUNTS_COL)
	return AccountsController{
		log: log.Named("AccountsController"), col: col, db: db,
	}
}