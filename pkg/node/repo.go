package node

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type Repo interface {
	CreateAccount(ctx context.Context, username, password string) (uid string, err error)
	IsAuthorized(ctx context.Context, target, who, action string) (decision bool, err error)
	Authorize(ctx context.Context, account, node, action string, inherit bool) (err error)
	GetAccount(ctx context.Context, name string) (account *nodepb.Account, err error)
	Authenticate(ctx context.Context, username, password string) (success bool, uid string, err error)

	CreateObject(ctx context.Context, name, parent, kind, namespace string) (id string, err error)
	DeleteObject(ctx context.Context, uid string) (err error)
	ListForAccount(ctx context.Context, account string) (directDevices []*nodepb.Device, directObjects []*nodepb.Object, inheritedObjects []*nodepb.Object, err error)
}
