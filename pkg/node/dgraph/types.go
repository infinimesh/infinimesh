package dgraph

type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type Account struct {
	Node
	Name           string              `json:"name,omitempty"`
	IsRoot         bool                `json:"isRoot"`
	AccessTo       *Object             `json:"access.to,omitempty"`
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
	Name               string  `json:"name,omitempty"`
	AccessToPermission string  `json:"access.to|permission,omitempty"`
	AccessToInherit    bool    `json:"access.to|inherit"`
	Contains           *Object `json:"contains"`
	ContainedIn        *Object `json:"~contains"` // !! Should only be used for delete
}

type ObjectList struct {
	Node
	Name        string       `json:"name,omitempty"`
	Contains    []ObjectList `json:"contains"`
	ContainedIn []ObjectList `json:"~contains"`
	AccessedBy  []Account    `json:"~access.to"`
}
