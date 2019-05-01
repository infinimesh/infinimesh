package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

type shadowAPI struct {
	accountClient nodepb.AccountServiceClient
	client        shadowpb.ShadowsClient
}

func (s *shadowAPI) Get(ctx context.Context, request *shadowpb.GetRequest) (response *shadowpb.GetResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := s.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return s.client.Get(ctx, request)
}

func (s *shadowAPI) PatchDesiredState(ctx context.Context, request *shadowpb.PatchDesiredStateRequest) (response *shadowpb.PatchDesiredStateResponse, err error) {
	account, ok := ctx.Value("account_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := s.accountClient.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}

	return s.client.PatchDesiredState(ctx, request)
}

func (s *shadowAPI) StreamReportedStateChanges(request *shadowpb.StreamReportedStateChangesRequest, srv apipb.States_StreamReportedStateChangesServer) (err error) {
	account, ok := srv.Context().Value("account_id").(string)
	if !ok {
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	resp, err := s.accountClient.IsAuthorized(srv.Context(), &nodepb.IsAuthorizedRequest{
		Node:    request.Id,
		Account: account,
		Action:  nodepb.Action_READ,
	})
	if err != nil {
		return status.Error(codes.PermissionDenied, "Permission denied")
	}
	if !resp.GetDecision().GetValue() {
		return status.Error(codes.PermissionDenied, "Permission denied")
	}

	c, err := s.client.StreamReportedStateChanges(srv.Context(), request)
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
