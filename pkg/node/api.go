package node

type Node struct {
	Type string `json:"type,omitempty"`
	UID  string `json:"uid,omitempty"`
}

type Account struct {
	Node
	Name           string    `json:"name,omitempty"`
	IsRoot         bool      `json:"isRoot"`
	AccessTo       *Resource `json:"access.to,omitempty"`
	AccessToDevice *Device   `json:"access.to.device,omitempty"`
}

type Resource struct {
	Node
	Name               string    `json:"name,omitempty"`
	AccessToPermission string    `json:"access.to|permission,omitempty"`
	AccessToInherit    bool      `json:"access.to|inherit"`
	Contains           *Resource `json:"contains"`
	ContainsDevice     *Device   `json:"contains.device"`
}

type Device struct {
	Node
	Name                     string `json:"name,omitempty"`
	AccessToDevicePermission string `json:"access.to.device|permission,omitempty"`
}
