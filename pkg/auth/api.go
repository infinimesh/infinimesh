package auth

type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type User struct {
	Node
	Name       string      `json:"name,omitempty"`
	Credential *Credential `json:"credential,omitempty"`
}

type Credential struct {
	Node
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"client_secret"`
}

type Clearance struct {
	Node
	GrantedTo User     `json:"granted.to,omitempty"`
	AccessTo  Resource `json:"access.to,omitempty"`
}

type Resource struct {
	Node
}
