package main

import (
	"context"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type shadowAPI struct {
	client shadowpb.ShadowsClient
}

func (s *shadowAPI) Get(ctx context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	return s.client.Get(ctx, request)
}

func (s *shadowAPI) PatchDesiredState(ctx context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {
	return s.client.PatchDesiredState(ctx, request)
}

func (s *shadowAPI) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv apipb.Shadows_StreamReportedStateChangesServer) (err error) {
	c, err := s.client.StreamReportedStateChanges(context.Background(), request)
	if err != nil {
		return err
	}

	for {
		msg, err := c.Recv()
		if err != nil {
			return err
		}

		err = srv.Send(msg)
		if err != nil {
			return err
		}
	}
}
