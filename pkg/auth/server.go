package auth

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/infinimesh/infinimesh/pkg/auth/authpb"
	"go.uber.org/zap"
)

type Server struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger
}

func (s *Server) Login(ctx context.Context, request *authpb.LoginRequest) (response *authpb.LoginResponse, err error) {
	return nil, nil
}

func (s *Server) CreateUser(ctx context.Context, request *authpb.CreateUserRequest) (response *authpb.CreateUserResponse, err error) {
	log := s.Log.Named("CreateUser")

	txn := s.Dgraph.NewTxn()

	q := `query userExists($name: string) {
                exists(func: eq(name, $name)) @filter(eq(type, "user")) {
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
		js, err := json.Marshal(&User{
			Node: Node{
				Type: "user",
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
		return &authpb.CreateUserResponse{Uid: userUID}, nil

	}
	return nil, errors.New("User exists already")
}

func (s *Server) SetCredentials(ctx context.Context, request *authpb.SetCredentialsRequest) (response *authpb.SetCredentialsResponse, err error) {
	return nil, nil
}

func checkExists(ctx context.Context, log *zap.Logger, txn *dgo.Txn, uid, _type string) bool {
	log = log.Named("checkExists")
	q := `query object($_uid: string, $type: string) {
                object(func: uid($_uid)) @filter(eq(type, $type)) {
                  uid
                }
              }
             `
	{

	}
	resp, err := txn.QueryWithVars(ctx, q, map[string]string{
		"$type": _type,
		"$_uid": uid,
	})
	if err != nil {
		log.Error("Query failed", zap.Error(err))
		return false
	}

	var result struct {
		Object []map[string]interface{} `json:"object"`
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Error("Failed to unmarshal response from dgraph", zap.Error(err))
		return false
	}

	return len(result.Object) > 0
}

func (s *Server) Authorize(ctx context.Context, request *authpb.AuthorizeRequest) (response *authpb.AuthorizeResponse, err error) {
	log := s.Log.Named("Authorize")

	txn := s.Dgraph.NewTxn()

	if ok := checkExists(ctx, log, txn, request.GetEntityUid(), "user"); !ok {
		return nil, errors.New("Entity does not exist")
	}

	// TODO optimize
	if ok := checkExists(ctx, log, txn, request.GetResourceUid(), "object"); !ok {
		if ok := checkExists(ctx, log, txn, request.GetResourceUid(), "device"); !ok {
			return nil, errors.New("resource does not exist")
		}
	}

	in := User{
		Node: Node{
			UID: request.GetEntityUid(),
		},
		AccessTo: &Resource{
			Node: Node{
				UID: request.GetResourceUid(),
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

	return &authpb.AuthorizeResponse{}, nil
}

func isPermissionSufficient(required, actual string) bool {
	switch required {
	case "WRITE":
		return actual == "WRITE"
	case "READ":
		return actual == "WRITE" || actual == "READ"
	default:
		return false
	}
}

func (s *Server) IsAuthorized(ctx context.Context, request *authpb.IsAuthorizedRequest) (response *authpb.IsAuthorizedResponse, err error) {

	if request.GetObject() == request.GetSubject() {
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: true}}, err
	}
	log := s.Log.Named("Authorize").With(
		zap.String("request.subject", request.GetSubject()),
		zap.String("request.action", request.GetAction().String()),
		zap.String("request.object", request.GetObject()),
	)

	params := map[string]string{
		"$device_id": request.GetObject(),
		"$user_id":   request.GetSubject(),
	}

	const qDirect = `query direct_access($device_id: string, $user_id: string){
                         direct(func: uid(0x9c70)) @normalize @cascade {
                           access.to  @filter(uid(0x9c7c) AND eq(type, "device")) @facets(permission,inherit) {
                             type: type
                           }
                         }

                         direct_via_one_object(func: uid(0x9c70)) @normalize @cascade {
                           access.to @filter(eq(type, "object")) @facets(permission,inherit) {
                             contains @filter(uid(0x9c77) AND eq(type, "device")) {
                               uid
                               type: type
                             }
                           }
                         }
                        }`

	res, err := s.Dgraph.NewTxn().QueryWithVars(ctx, qDirect, params)
	if err != nil {
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	var permissions struct {
		Direct          []Resource `json:"direct"`
		DirectViaObject []Resource `json:"direct_via_one_object"`
	}

	err = json.Unmarshal(res.Json, &permissions)
	if err != nil {
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	log.Debug("Dgraph response", zap.Any("json", permissions))

	if len(permissions.Direct) > 0 {
		if isPermissionSufficient(request.GetAction().String(), permissions.Direct[0].AccessToPermission) {
			log.Info("Granting access")
			return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: true}}, err
		}
	}

	if len(permissions.DirectViaObject) > 0 {
		if isPermissionSufficient(request.GetAction().String(), permissions.DirectViaObject[0].AccessToPermission) {
			log.Info("Granting access")
			return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: true}}, err
		}
	}

	// TODO: recursive lookup if inherit=true

	log.Info("Denying access")
	return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
}
