//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package registry

import (
	"context"
	"encoding/json"

	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dgraph-io/dgo"
)

//DGraphRepo is a Data type for executing Dgraph Query
type DGraphRepo struct {
	Dg *dgo.Dgraph
}

//List is a method to execute Dgraph Query to List details of all Devices
func (dr *DGraphRepo) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	txn := dr.Dg.NewReadOnlyTxn()

	const q = `query list($namespaceid: string){
		var(func: uid($namespaceid)) @filter(eq(type, "namespace")) {
		  owns {
			OBJs as uid
		  } @filter(eq(kind, "device"))
		}

		nodes(func: uid(OBJs)) @recurse {
		  children{}
		  uid
		  name
		  kind
		  enabled
		  tags
		}
	  }`

	vars := map[string]string{
		"$namespaceid": request.Namespace,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Nodes []Device `json:"nodes"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	var devices []*registrypb.Device
	for _, device := range res.Nodes {
		devices = append(devices, toProto(&device))
	}

	return &registrypb.ListResponse{
		Devices: devices,
	}, nil
}
