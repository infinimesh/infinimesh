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

func (s *AccountController) IsRoot(ctx context.Context, request *nodepb.IsRootRequest) (response *nodepb.IsRootResponse, err error) {
	account, err := s.Repo.GetAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find account")
	}

	return &nodepb.IsRootResponse{IsRoot: account.IsRoot}, nil
}

func (s *AccountController) CreateUserAccount(ctx context.Context, request *nodepb.CreateUserAccountRequest) (response *nodepb.CreateUserAccountResponse, err error) {
	log := s.Log.Named("CreateUserAccount")
	uid, err := s.Repo.CreateUserAccount(ctx, request.Account.Name, request.Password, request.Account.IsRoot, request.Account.Enabled)
	if err != nil {
		log.Error("Failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	log.Info("Successfully created account", zap.String("username", request.Account.Name), zap.String("password", request.Password), zap.String("uid", uid))

	return &nodepb.CreateUserAccountResponse{Uid: uid}, nil
}

func (s *AccountController) AuthorizeNamespace(ctx context.Context, request *nodepb.AuthorizeNamespaceRequest) (response *nodepb.AuthorizeNamespaceResponse, err error) {
	err = s.Repo.AuthorizeNamespace(ctx, request.GetAccount(), request.GetNamespace(), request.GetAction())

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to authorize")
	}

	return &nodepb.AuthorizeNamespaceResponse{}, nil
}

func (s *AccountController) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {
	err = s.Repo.Authorize(ctx, request.GetAccount(), request.GetNode(), request.GetAction(), request.GetInherit())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to authorize")
	}

	return &nodepb.AuthorizeResponse{}, nil
}

func (s *AccountController) IsAuthorizedNamespace(ctx context.Context, request *nodepb.IsAuthorizedNamespaceRequest) (response *nodepb.IsAuthorizedNamespaceResponse, err error) {
	root, err := s.IsRoot(ctx, &nodepb.IsRootRequest{
		Account: request.GetAccount(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Authorization check failed")
	}

	if root.GetIsRoot() {
		return &nodepb.IsAuthorizedNamespaceResponse{
			Decision: &wrappers.BoolValue{Value: true},
		}, nil
	}

	decision, err := s.Repo.IsAuthorizedNamespace(ctx, request.GetNamespace(), request.GetAccount(), request.GetAction())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.IsAuthorizedNamespaceResponse{Decision: &wrappers.BoolValue{Value: decision}}, nil
}

func (s *AccountController) SetPassword(ctx context.Context, request *nodepb.SetPasswordRequest) (response *nodepb.SetPasswordResponse, err error) {
	err = s.Repo.SetPassword(ctx, request.Username, request.Password)
	if err != nil {
		return &nodepb.SetPasswordResponse{}, err
	}

	return &nodepb.SetPasswordResponse{}, nil
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
	account, err := s.Repo.GetAccount(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return account, nil
}

func (s *AccountController) Authenticate(ctx context.Context, request *nodepb.AuthenticateRequest) (response *nodepb.AuthenticateResponse, err error) {
	ok, uid, defaultNs, err := s.Repo.Authenticate(ctx, request.GetUsername(), request.GetPassword())
	if !ok || (err != nil) {

		return &nodepb.AuthenticateResponse{}, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	return &nodepb.AuthenticateResponse{Success: ok, Account: &nodepb.Account{Uid: uid}, DefaultNamespace: defaultNs}, nil
}

func (s *AccountController) ListAccounts(ctx context.Context, request *nodepb.ListAccountsRequest) (response *nodepb.ListAccountsResponse, err error) {
	accounts, err := s.Repo.ListAccounts(ctx)
	if err != nil {
		return &nodepb.ListAccountsResponse{}, status.Error(codes.Internal, "Failed to list accounts")
	}

	return &nodepb.ListAccountsResponse{
		Accounts: accounts,
	}, nil
}

func (s *AccountController) UpdateAccount(ctx context.Context, request *nodepb.UpdateAccountRequest) (response *nodepb.Account, err error) {
	err = s.Repo.UpdateAccount(ctx, request)
	if err != nil {
		return &nodepb.Account{}, err
	}
	return request.Account, nil
}
