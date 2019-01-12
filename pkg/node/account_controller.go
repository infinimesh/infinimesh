package node

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountController struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
}

func (s *AccountController) CreateAccount(ctx context.Context, request *nodepb.CreateAccountRequest) (response *nodepb.CreateAccountResponse, err error) {
	log := s.Log.Named("CreateAccount")

	txn := s.Dgraph.NewTxn()

	q := `query userExists($name: string) {
                exists(func: eq(name, $name)) @filter(eq(type, "account")) {
                  uid
                }
              }
             `

	var result struct {
		Exists []map[string]interface{} `json:"exists"`
	}

	resp, err := txn.QueryWithVars(ctx, q, map[string]string{"$name": request.GetName()})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Exists) == 0 {
		js, err := json.Marshal(&Account{
			Node: Node{
				Type: "account",
				UID:  "_:user",
			},
			Name: request.GetName(),
		})
		if err != nil {
			return nil, err
		}
		m := &api.Mutation{SetJson: js}
		a, err := txn.Mutate(ctx, m)
		if err != nil {
			return nil, err
		}

		err = txn.Commit(ctx)
		if err != nil {
			log.Error("Failed to commit txn", zap.Error(err))
			return nil, errors.New("Failed to commit")
		}
		userUID := a.GetUids()["user"]
		log.Info("Created user", zap.String("name", request.GetName()), zap.String("uid", userUID))
		return &nodepb.CreateAccountResponse{Uid: userUID}, nil

	}
	return nil, errors.New("User exists already")
}

func (s *AccountController) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {
	log := s.Log.Named("Authorize")

	txn := s.Dgraph.NewTxn()

	if ok := checkExists(ctx, log, txn, request.GetAccount(), "user"); !ok {
		return nil, errors.New("Entity does not exist")
	}

	if ok := checkExists(ctx, log, txn, request.GetNode(), "object"); !ok {
		if ok := checkExists(ctx, log, txn, request.GetNode(), "device"); !ok {
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
