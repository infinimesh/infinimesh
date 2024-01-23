/*
Copyright © 2021-2023 Infinite Devices GmbH

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
	"errors"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	accountspb "github.com/infinimesh/proto/node/accounts"
	"go.uber.org/zap"
)

type Link struct {
	From driver.DocumentID `json:"_from"`
	To   driver.DocumentID `json:"_to"`
	Type string            `json:"type"`

	driver.DocumentMeta
}

type Credentials interface {
	// Check if given authorization data are mapped
	// to existent Credentials
	Authorize(...string) bool
	// Return Credentials type
	Type() string
	// Return Credentials Document Key
	Key() string

	// Find Credentials in database by authorisation data and Unmarshall it's data into struct
	Find(context.Context, driver.Database) bool
	// Find Credentials in database by document key and Unmarshall it's data into struct
	FindByKey(context.Context, driver.Collection, string) error

	// Set Logger for Credentials methods
	SetLogger(*zap.Logger)
}

func Determine(auth_type string) (cred Credentials, ok bool) {
	switch auth_type {
	case "standard":
		return &StandardCredentials{}, true
	default:
		return nil, false
	}
}

func Find(ctx context.Context, db driver.Database, log *zap.Logger, auth_type string, args ...string) (cred Credentials, err error) {
	var ok bool
	switch auth_type {
	case "standard":
		cred = &StandardCredentials{Username: args[0]}
	case "ldap":
		cred = &LDAPCredentials{Username: args[0]}
	default:
		return nil, errors.New("unknown auth type")
	}

	cred.SetLogger(log)

	ok = cred.Find(ctx, db)
	if !ok {
		return nil, errors.New("couldn't find credentials")
	}

	if cred.Authorize(args...) {
		return cred, nil
	}

	return nil, errors.New("couldn't authorize")
}

func MakeCredentials(credentials *accountspb.Credentials, log *zap.Logger) (Credentials, error) {
	if credentials == nil {
		return nil, errors.New("credentials aren't given")
	}

	var cred Credentials
	var err error
	switch credentials.Type {
	case "standard":
		cred, err = NewStandardCredentials(credentials.Data[0], credentials.Data[1])
	case "ldap":
		if len(credentials.Data) != 2 {
			return nil, errors.New("missing LDAP provider key")
		}
		cred, err = NewLDAPCredentials(credentials.Data[0], credentials.Data[1])
	default:
		return nil, errors.New("auth type is wrong")
	}

	if err != nil {
		return nil, err
	}

	cred.SetLogger(log)
	return cred, nil
}

const listCredentialsAndEdgesQuery = `
RETURN FLATTEN(
FOR node, edge IN 1 OUTBOUND @account
GRAPH @credentials
    RETURN [ node._id, edge._id ]
)
`

func ListCredentialsAndEdges(ctx context.Context, log *zap.Logger, db driver.Database, account driver.DocumentID) (nodes []string, err error) {
	c, err := db.Query(ctx, listCredentialsAndEdgesQuery, map[string]interface{}{
		"account":     account,
		"credentials": schema.CREDENTIALS_COL,
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	_, err = c.ReadDocument(ctx, &nodes)
	return nodes, err
}

type ListCredentialsResponse struct {
	Type string                 `json:"type"`
	D    map[string]interface{} `json:"credentials"`
}

const listCredentialsQuery = `
FOR credentials, edge IN 1 OUTBOUND @account
GRAPH @credentials_graph
RETURN { type: edge.type, credentials }
`

// ListCredentials - Returns Credentials linked to Account
func ListCredentials(ctx context.Context, log *zap.Logger, db driver.Database, acc driver.DocumentID) (r []ListCredentialsResponse, err error) {
	c, err := db.Query(ctx, listCredentialsQuery, map[string]interface{}{
		"account":           acc.String(),
		"credentials_graph": schema.CREDENTIALS_GRAPH.Name,
	})
	if err != nil {
		log.Warn("Error executing query", zap.Error(err))
		return nil, err
	}
	defer c.Close()

	for {
		var cred ListCredentialsResponse
		_, err := c.ReadDocument(ctx, &cred)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Debug("Error unmarshalling credentials response", zap.Error(err))
			return nil, err
		}
		r = append(r, cred)
	}

	return r, nil
}

type ListableCredentials interface {
	Listable() []string
}
type ListableFabric func(map[string]interface{}) (ListableCredentials, error)

var _Listables = map[string]ListableFabric{
	"standard": StandardFromMap,
}

// MakeListable - Accepts Credentials type as string t and Credentials data as map[string]interface{} d
func MakeListable(r ListCredentialsResponse) (ListableCredentials, error) {
	f, ok := _Listables[r.Type]
	if !ok {
		return nil, fmt.Errorf("Credentials of type %s aren't Listable", r.Type)
	}

	return f(r.D)
}
