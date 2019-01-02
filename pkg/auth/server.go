package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
		fmt.Println("zz")
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
	// Upsert add node user -> CRED -> resource

	log := s.Log.Named("Authorize")

	txn := s.Dgraph.NewTxn()

	if ok := checkExists(ctx, log, txn, request.GetEntityUid(), "user"); !ok {
		return nil, errors.New("Entity does not exist")
	}

	if ok := checkExists(ctx, log, txn, request.GetResourceUid(), "resource"); !ok {
		return nil, errors.New("Resource does not exist")
	}

	in := User{
		Node: Node{
			UID: request.GetEntityUid(),
		},
		AccessTo: &Resource{
			Node: Node{
				UID: request.GetResourceUid(),
			},
			AccessToPermission: "WRITE",
		},
	}

	js, err := json.Marshal(&in)
	if err != nil {
		return nil, err
	}

	a, err := txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		log.Info("Mutate fail", zap.Error(err))
		return nil, errors.New("Failed to mutate")
	}

	spew.Dump(a)

	return &authpb.AuthorizeResponse{}, nil
}

func (s *Server) IsAuthorized(ctx context.Context, request *authpb.IsAuthorizedRequest) (response *authpb.IsAuthorizedResponse, err error) {

	if request.GetObject() == request.GetSubject() {
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: true}}, err
	}
	log := s.Log.Named("Authorize").With(
		zap.String("subject", request.GetSubject()),
		zap.String("action", request.GetAction()),
		zap.String("object", request.GetObject()),
	)

	// Compute all clearances which are sufficient to perform the requested action
	var sufficientClearances string
	switch request.GetAction() {
	case "write":
		sufficientClearances = "write"
	case "read":
		sufficientClearances = "write read"
	default:
		return nil, errors.New("Invalid action")
	}

	params := map[string]string{
		"$device_id":    request.GetObject(),
		"$subject_uuid": request.GetSubject(),
		"$action":       sufficientClearances,
	}
	const q = `query permissions($action: string, $device_id: string, $subject_uuid: string){
                     var(func: eq(device_id,$device_id)) @recurse @normalize @cascade {
                       parentObjectUIDs as uid
                       contained_in  {
                       }
                     }

                     var(func: uid(parentObjectUIDs)) @normalize  @cascade {
                       accessed_through @filter(anyofterms(action, $action)) {
                         clearanceIDs as uid
                       }
                     }

                     clearance(func: uid(clearanceIDs), first: 1) @cascade {
                       uid
                       action
                       granted_to @filter(eq(uuid, $subject_uuid)) {}
                     }
                   }`

	res, err := s.Dgraph.NewTxn().QueryWithVars(ctx, q, params)
	if err != nil {
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	type Clearance struct {
		Action string `json:"action"`
	}

	type Permissions struct {
		Permissions []Clearance `json:"clearance"`
	}

	var p Permissions
	err = json.Unmarshal(res.Json, &p)
	if err != nil {
		// s.Log.Info("Failed to unmarshal result from dgraph", fields ...zapcore.Field)
		return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
	}

	if len(p.Permissions) > 0 {
		permission := p.Permissions[0]
		if permission.Action == request.Action {
			log.Info("Granting access")
			return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: true}}, err
		}
	}

	log.Info("Denying access")
	return &authpb.IsAuthorizedResponse{Decision: &wrappers.BoolValue{Value: false}}, err
}
