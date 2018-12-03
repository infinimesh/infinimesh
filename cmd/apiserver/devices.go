package main

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

type deviceAPI struct {
	client registrypb.DevicesClient
}

func (d *deviceAPI) Create(ctx context.Context, request *registrypb.CreateRequest) (response *registrypb.CreateResponse, err error) {
	return d.client.Create(ctx, request)
}
func (d *deviceAPI) GetByName(ctx context.Context, request *registrypb.GetByNameRequest) (response *registrypb.GetByNameResponse, err error) {
	return d.client.GetByName(ctx, request)

}
func (d *deviceAPI) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	return d.client.List(ctx, request)
}
func (d *deviceAPI) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	return d.client.Delete(ctx, request)
}
