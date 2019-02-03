package node

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
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
	log := s.Log.Named("Authorize")

	txn := s.Dgraph.NewTxn()

	if ok := checkExists(ctx, txn, request.GetAccount(), "user"); !ok {
		return nil, errors.New("Entity does not exist")
	}

	if ok := checkExists(ctx, txn, request.GetNode(), "object"); !ok {
		if ok := checkExists(ctx, txn, request.GetNode(), "device"); !ok {
			return nil, errors.New("resource does not exist")
		}
	}

	in := Account{
		Node: Node{
			UID: request.GetAccount(),
		},
		AccessTo: &Object{
			Node: Node{
				UID: request.GetNode(),
			},
			AccessToPermission: request.GetAction(),
			AccessToInherit:    request.GetInherit(),
		},
	}

	js, err := json.Marshal(&in)
	if err != nil {
		return nil, err
	}

	log.Debug("Run mutation", zap.Any("json", &in))

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		log.Info("Mutate fail", zap.Error(err))
		return nil, errors.New("Failed to mutate")
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
	return &nodepb.Account{
		Uid:  account.UID,
		Name: account.Name,
	}, nil
}

func (s *AccountController) Authenticate(ctx context.Context, request *nodepb.AuthenticateRequest) (response *nodepb.AuthenticateResponse, err error) {

	txn := s.Dgraph.NewReadOnlyTxn()

	const q = `query authenticate($username: string, $password: string){
  login(func: eq(username, $username)) @filter(eq(type, "credentials")) {
    uid
    checkpwd(password, $password)
    ~has.credentials {
      uid
      type
    }
  }
}
`

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$username": request.GetUsername(), "$password": request.GetPassword()})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result struct {
		Login []*UsernameCredential `json:"login"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(result.Login) > 0 {
		login := result.Login[0]
		if login.CheckPwd {
			// Success
			if len(login.Account) > 0 {
				return &nodepb.AuthenticateResponse{
					Success: result.Login[0].CheckPwd,
					Account: &nodepb.Account{
						Uid: login.Account[0].UID,
					},
				}, nil
			}
		}
	}

	return nil, status.Error(codes.Unauthenticated, "Invalid credentials")

}
