package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type namespaceAPI struct {
	client        nodepb.NamespaceServiceClient
	accountClient nodepb.AccountServiceClient
}

func (n *namespaceAPI) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check permissions")
	}

	if resp.GetIsRoot() {
		return n.client.ListNamespaces(ctx, &nodepb.ListNamespacesRequest{})
	} else {
		return n.client.ListNamespacesForAccount(ctx, &nodepb.ListNamespacesForAccountRequest{
			Account: account,
		})
	}
}

func (n *namespaceAPI) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.CreateNamespaceResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := n.accountClient.IsRoot(ctx, &nodepb.IsRootRequest{Account: account})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check permissions")
	}

	if resp.GetIsRoot() {
		return n.client.CreateNamespace(ctx, request)
	}
	return nil, status.Error(codes.PermissionDenied, "Account is not root")
}
