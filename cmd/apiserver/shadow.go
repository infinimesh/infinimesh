package main

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type shadowAPI struct {
	client shadowpb.ShadowClient
}

func (s *shadowAPI) GetReported(ctx context.Context, request *shadowpb.GetReportedRequest) (response *shadowpb.GetReportedResponse, err error) {
	return s.client.GetReported(ctx, request)
}
