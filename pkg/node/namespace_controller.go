package node

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type NamespaceController struct {
	Repo Repo
}

func (n *NamespaceController) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.CreateNamespaceResponse, err error) {
	return nil, nil
}
func (n *NamespaceController) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.GetNamespaceResponse, err error) {
	return nil, nil
}
func (n *NamespaceController) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {
	return nil, nil
}
