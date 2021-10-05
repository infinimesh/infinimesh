//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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

package dgraph

//Node Data strcuture for Dgraph database
type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

//Namespace Data structure for Dgraph database
type Namespace struct {
	Node
	Name                 string `json:"name,omitempty"`
	MarkForDeletion      bool   `json:"markfordeletion,omitempty"`
	DeleteInitiationTime string `json:"deleteinitiationtime,omitempty"`
	RetentionPeriod      int32  `json:"retentionperiod,omitempty"`

	Owns []*Object `json:"owns,omitempty"`

	AccessedBy         []Account `json:"~access.to.namespace"`
	AccessToPermission string    `json:"access.to.namespace|permission,omitempty"`
}

//Account Data structure for Dgraph database
type Account struct {
	Node
	Name string `json:"name,omitempty"`

	IsRoot  bool `json:"isRoot"`
	IsAdmin bool `json:"isAdmin"`
	Enabled bool `json:"enabled"`

	AccessTo          []*Object    `json:"access.to,omitempty"`
	AccessToNamespace []*Namespace `json:"access.to.namespace,omitempty"`

	DefaultNamespace []*Namespace `json:"default.namespace"`

	HasCredentials []*UsernameCredential `json:"has.credentials,omitempty"`

	Owns    []*Account `json:"owns,omitempty"`
	OwnedBy []*Account `json:"~owns,omitempty"`
}

//UsernameCredential Data structure for Dgraph database
type UsernameCredential struct {
	Node
	Username string     `json:"username"`
	Password string     `json:"password"`
	CheckPwd bool       `json:"checkpwd(password),omitempty"`
	Account  []*Account `json:"~has.credentials,omitempty"`
}

//Object Data structure for Dgraph database
type Object struct {
	Node

	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`

	OwnedBy []*Namespace `json:"~owns,omitempty"`

	Children []*Object `json:"children"`
	Parent   []Node    `json:"~children"` // Namespace or Object

	AccessedBy         []Account `json:"~access.to"`
	AccessToPermission string    `json:"access.to|permission,omitempty"`
	AccessToInherit    bool      `json:"access.to|inherit"`
}
