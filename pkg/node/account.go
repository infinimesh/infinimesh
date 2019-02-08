package node

import (
	"context"

	"github.com/dgraph-io/dgo"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type AccountController struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
}

func (s *AccountController) CreateAccount(ctx context.Context, request *nodepb.CreateAccountRequest) (response *nodepb.CreateAccountResponse, err error) {
	log := s.Log.Named("CreateAccount")
	uid, err := s.Repo.CreateAccount(ctx, request.GetName(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	log.Info("Successfully created account", zap.String("username", request.GetName()), zap.String("password", request.GetPassword()), zap.String("uid", uid))

	return &nodepb.CreateAccountResponse{Uid: uid}, nil
}

func (s *AccountController) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {
	err = s.Repo.Authorize(ctx, request.GetAccount(), request.GetNode(), request.GetAction(), request.GetInherit())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to authorize")
	}

	return &nodepb.AuthorizeResponse{}, nil
}

func (s *AccountController) IsAuthorized(ctx context.Context, request *nodepb.IsAuthorizedRequest) (response *nodepb.IsAuthorizedResponse, err error) {
	log := s.Log.Named("Authorize").With(
		zap.String("request.account", request.GetAccount()),
		zap.String("request.action", request.GetAction().String()),
		zap.String("request.node", request.GetNode()),
	)

	decision, err := s.Repo.IsAuthorized(ctx, request.GetNode(), request.GetAccount(), request.GetAction().String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info("Return decision", zap.Bool("decision", decision))
	return &nodepb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

func (s *AccountController) GetAccount(ctx context.Context, request *nodepb.GetAccountRequest) (response *nodepb.Account, err error) {
	account, err := s.Repo.GetAccount(ctx, request.GetName())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return account, nil
}

func (s *AccountController) Authenticate(ctx context.Context, request *nodepb.AuthenticateRequest) (response *nodepb.AuthenticateResponse, err error) {

	ok, uid, err := s.Repo.Authenticate(ctx, request.GetUsername(), request.GetPassword())
	if !ok || (err != nil) {

		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	return &nodepb.AuthenticateResponse{Success: ok, Account: &nodepb.Account{Uid: uid}}, nil
}
