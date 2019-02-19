package dgraph

type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type Namespace struct {
	Node
	Name string `json:"name,omitempty"`

	Owns []*Object `json:"owns,omitempty"`

	AccessedBy         []Account `json:"~access.to.namespace"`
	AccessToPermission string    `json:"access.to.namespace|permission,omitempty"`
}

type Account struct {
	Node
	Name string `json:"name,omitempty"`

	IsRoot bool `json:"isRoot,omitempty"`

	AccessTo          []*Object    `json:"access.to,omitempty"`
	AccessToNamespace []*Namespace `json:"access.to.namespace,omitempty"`

	HasCredentials *UsernameCredential `json:"has.credentials,omitempty"`
}

type UsernameCredential struct {
	Node
	Username string     `json:"username"`
	Password string     `json:"password"`
	CheckPwd bool       `json:"checkpwd(password),omitempty"`
	Account  []*Account `json:"~has.credentials,omitempty"`
}

type Object struct {
	Node

	Name string `json:"name,omitempty"`
	Kind string `json:"kind,omitempty"`

	OwnedBy *Namespace `json:"~owns,omitempty"`

	Children []*Object `json:"children"`
	Parent   []Node    `json:"~children"` // Namespace or Object

	AccessedBy         []Account `json:"~access.to"`
	AccessToPermission string    `json:"access.to|permission,omitempty"`
	AccessToInherit    bool      `json:"access.to|inherit"`
}
