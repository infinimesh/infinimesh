package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type objectAPI struct {
	objectClient  nodepb.ObjectServiceClient
	accountClient nodepb.AccountServiceClient
}

func (o *objectAPI) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := o.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.GetParent(),
		Account: account,
		Action:  nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Decision.GetValue() {
		return nil, status.Error(codes.PermissionDenied, "No permission to access resource")
	}

	return o.objectClient.CreateObject(ctx, request)
}

func (o *objectAPI) ListObjects(ctx context.Context, request *apipb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return o.objectClient.ListObjects(ctx, &nodepb.ListObjectsRequest{Account: account})
}

func (o *objectAPI) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	_ = account
	//TODO

	return nil, nil
}
