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
package schema

var (
	DB_NAME = "infinimesh"
)

type infinimeshClaim string
const (
	InfinimeshAccount infinimeshClaim = "requestorID"
)

const (
	ROOT_ACCOUNT_KEY = "infinimesh"
	ROOT_NAMESPACE_KEY = "infinimesh"
)

const (
	ACCOUNTS_COL = "Accounts"
	ACC2NS = ACCOUNTS_COL + "2" + NAMESPACES_COL
	ACC2CRED = ACCOUNTS_COL + "2" + CREDENTIALS_COL
)

const (
	NAMESPACES_COL = "Namespaces"
	NS2ACC = NAMESPACES_COL + "2" + ACCOUNTS_COL
)

const (
	CREDENTIALS_COL = "Credentials"
	CREDENTIALS_EDGE_COL = ACCOUNTS_COL + "2" + CREDENTIALS_COL
)

const (
	DEVICES_COL = "Devices"
	DEVICES_EDGE_COL = ACCOUNTS_COL + "2" + DEVICES_COL
)

type InfinimeshGraphSchema struct {
	Name string
	Edges [][]string
}

var COLLECTIONS = []string{
	ACCOUNTS_COL, NAMESPACES_COL,
	CREDENTIALS_COL, DEVICES_COL,
}

var PERMISSIONS_GRAPH = InfinimeshGraphSchema{
	Name: "Permissions",
	Edges: [][]string{
		{ACCOUNTS_COL, NAMESPACES_COL},
		{NAMESPACES_COL, ACCOUNTS_COL},
		{NAMESPACES_COL, DEVICES_COL},
	},
}
var CREDENTIALS_GRAPH = InfinimeshGraphSchema{
	Name: "Credentials",
	Edges: [][]string{
		{ACCOUNTS_COL, CREDENTIALS_COL},
	},
}

var GRAPHS_SCHEMAS = []InfinimeshGraphSchema{
	PERMISSIONS_GRAPH, CREDENTIALS_GRAPH,
}

type InfinimeshAccessLevel int32;
const NONE InfinimeshAccessLevel = 0;
const READ InfinimeshAccessLevel = 1;
const MGMT InfinimeshAccessLevel = 2;
const ADMIN InfinimeshAccessLevel = 3;
const ROOT InfinimeshAccessLevel = 4;
