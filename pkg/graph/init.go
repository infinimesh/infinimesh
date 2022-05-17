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
	"errors"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/proto/node/access"
	accpb "github.com/infinimesh/proto/node/accounts"
	nspb "github.com/infinimesh/proto/node/namespaces"
	"go.uber.org/zap"
)

func EnsureRootExists(_log *zap.Logger, db driver.Database, passwd string) (err error) {

	ctx := context.TODO()
	log := _log.Named("EnsureRootExists")

	log.Debug("Checking Root Account exists")
	col, _ := db.Collection(ctx, schema.ACCOUNTS_COL)
	exists, err := col.DocumentExists(ctx, schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Error("Error checking Root Account existance")
		return err
	}

	var meta driver.DocumentMeta
	if !exists {
		log.Debug("Root Account doesn't exist, creating")
		meta, err = col.CreateDocument(ctx, Account{
			Account: &accpb.Account{
				Title:   "infinimesh",
				Enabled: true,
			},
			DocumentMeta: driver.DocumentMeta{Key: schema.ROOT_ACCOUNT_KEY},
		})
		if err != nil {
			log.Error("Error creating Root Account")
			return err
		}
		log.Debug("Created root Account", zap.Any("result", meta))
	}
	var acc accpb.Account
	meta, err = col.ReadDocument(ctx, schema.ROOT_ACCOUNT_KEY, &acc)
	if err != nil {
		log.Error("Error reading Root Account")
		return err
	}
	root := &Account{
		Account:      &acc,
		DocumentMeta: meta,
	}

	ns_col, _ := db.Collection(ctx, schema.NAMESPACES_COL)
	exists, err = ns_col.DocumentExists(ctx, schema.ROOT_NAMESPACE_KEY)
	if err != nil || !exists {
		meta, err := ns_col.CreateDocument(ctx, Namespace{
			Namespace: &nspb.Namespace{
				Title: "infinimesh",
			},
			DocumentMeta: driver.DocumentMeta{Key: schema.ROOT_NAMESPACE_KEY},
		})
		if err != nil {
			log.Error("Error creating Root Namespace")
			return err
		}
		log.Debug("Created root Namespace", zap.Any("result", meta))
	}

	var ns nspb.Namespace
	meta, err = ns_col.ReadDocument(ctx, schema.ROOT_NAMESPACE_KEY, &ns)
	if err != nil {
		log.Error("Error reading Root Namespace")
		return err
	}
	rootNS := &Namespace{
		Namespace:    &ns,
		DocumentMeta: meta,
	}

	edge_col := GetEdgeCol(ctx, db, schema.ACC2NS)
	exists = CheckLink(ctx, edge_col, root, rootNS)
	if err != nil {
		log.Error("Error checking link Root Account to Root Namespace", zap.Error(err))
		return err
	} else if !exists {
		err = Link(ctx, log, edge_col, root, rootNS, access.Level_ROOT, access.Role_OWNER)
		if err != nil {
			log.Error("Error linking Root Account to Root Namespace")
			return err
		}
	}

	ctx = context.WithValue(ctx, schema.InfinimeshAccount, schema.ROOT_ACCOUNT_KEY)
	cred_edge_col, _ := db.Collection(ctx, schema.ACC2CRED)
	cred, err := credentials.NewStandardCredentials("infinimesh", passwd)
	if err != nil {
		log.Error("Error creating Root Account Credentials")
		return err
	}

	ctrl := NewAccountsController(log, db)
	exists, err = cred_edge_col.DocumentExists(ctx, fmt.Sprintf("standard-%s", schema.ROOT_ACCOUNT_KEY))
	if err != nil || !exists {
		err = ctrl.SetCredentialsCtrl(ctx, *root, cred_edge_col, cred)
		if err != nil {
			log.Error("Error setting Root Account Credentials")
			return err
		}
	}
	_, r := ctrl.Authorize(ctx, "standard", "infinimesh", passwd)
	if !r {
		log.Error("Error authorizing Root Account")
		return errors.New("cannot authorize infinimesh")
	}
	return nil
}
