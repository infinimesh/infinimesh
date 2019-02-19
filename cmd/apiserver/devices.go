package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

type deviceAPI struct {
	client        registrypb.DevicesClient
	accountClient nodepb.AccountServiceClient
}

func (d *deviceAPI) Create(ctx context.Context, request *registrypb.CreateRequest) (response *registrypb.CreateResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := d.accountClient.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
		Namespace: request.GetNamespace(),
		Account:   account,
		Action:    nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	return d.client.Create(ctx, request)
}

func (d *deviceAPI) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {
	return d.client.Update(ctx, request)
}

func (d *deviceAPI) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {
	return d.client.Get(ctx, request)

}
func (d *deviceAPI) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	return d.client.List(ctx, request)
}
func (d *deviceAPI) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	return d.client.Delete(ctx, request)
}
