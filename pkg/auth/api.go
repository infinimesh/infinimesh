package auth

type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type User struct {
	Node
	Name               string      `json:"name,omitempty"`
	Credential         *Credential `json:"credential,omitempty"`
	AccessTo           *Resource   `json:"access.to,omitempty"`
	AccessToPermission string      `json:"access.to|permission,omitempty"`
}

type Credential struct {
	Node
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"client_secret"`
}

type Resource struct {
	Node
	AccessToPermission string `json:"access.to|permission,omitempty"`
}
