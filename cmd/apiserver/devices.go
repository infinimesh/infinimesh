package main

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
)

type deviceAPI struct {
}

func (d *deviceAPI) Create(context.Context, *registrypb.CreateRequest) (*registrypb.CreateResponse, error) {
	return &registrypb.CreateResponse{}, nil
}
func (d *deviceAPI) GetByName(context.Context, *registrypb.GetByNameRequest) (*registrypb.GetByNameResponse, error) {
	return &registrypb.GetByNameResponse{}, nil

}
func (d *deviceAPI) List(context.Context, *registrypb.ListDevicesRequest) (*registrypb.ListResponse, error) {
	return &registrypb.ListResponse{}, nil

}
func (d *deviceAPI) Delete(context.Context, *registrypb.DeleteRequest) (*registrypb.DeleteResponse, error) {
	return &registrypb.DeleteResponse{}, nil
}
