package node

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type NamespaceController struct {
	Repo Repo
}

func (n *NamespaceController) CreateNamespace(ctx context.Context, request *nodepb.CreateNamespaceRequest) (response *nodepb.Namespace, err error) {
	id, err := n.Repo.CreateNamespace(ctx, request.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create namespace")
	}

	return &nodepb.Namespace{
		Id:   id,
		Name: request.GetName(),
	}, nil
}

func (n *NamespaceController) ListNamespaces(ctx context.Context, request *nodepb.ListNamespacesRequest) (response *nodepb.ListNamespacesResponse, err error) {
	namespaces, err := n.Repo.ListNamespaces(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed list namespaces")
	}

	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

func (n *NamespaceController) ListNamespacesForAccount(ctx context.Context, request *nodepb.ListNamespacesForAccountRequest) (response *nodepb.ListNamespacesResponse, err error) {
	namespaces, err := n.Repo.ListNamespacesForAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed list namespaces")
	}

	return &nodepb.ListNamespacesResponse{
		Namespaces: namespaces,
	}, nil
}

func (n *NamespaceController) GetNamespace(ctx context.Context, request *nodepb.GetNamespaceRequest) (response *nodepb.Namespace, err error) {
	namespace, err := n.Repo.GetNamespace(ctx, request.GetNamespace())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed get namespace")
	}

	return namespace, nil
}

func (n *NamespaceController) ListPermissions(ctx context.Context, request *nodepb.ListPermissionsRequest) (response *nodepb.ListPermissionsResponse, err error) {
	permissions, err := n.Repo.ListPermissionsInNamespace(ctx, request.Namespace)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed get permissions")
	}

	return &nodepb.ListPermissionsResponse{Permissions: permissions}, nil
}

func (n *NamespaceController) DeletePermission(ctx context.Context, request *nodepb.DeletePermissionRequest) (response *nodepb.DeletePermissionResponse, err error) {
	err = n.Repo.DeletePermissionInNamespace(ctx, request.Namespace, request.AccountId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete permission")
	}
	return &nodepb.DeletePermissionResponse{}, nil
}
