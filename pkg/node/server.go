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

type Server struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
}

func (s *Server) CreateAccount(ctx context.Context, request *nodepb.CreateAccountRequest) (response *nodepb.CreateAccountResponse, err error) {
	log := s.Log.Named("CreateAccount")

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
		js, err := json.Marshal(&Account{
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
		return &nodepb.CreateAccountResponse{Uid: userUID}, nil

	}
	return nil, errors.New("User exists already")
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

func (s *Server) Authorize(ctx context.Context, request *nodepb.AuthorizeRequest) (response *nodepb.AuthorizeResponse, err error) {
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

func (s *Server) IsAuthorized(ctx context.Context, request *nodepb.IsAuthorizedRequest) (response *nodepb.IsAuthorizedResponse, err error) {
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

func (s *Server) checkPerm(ctxAccount interface{}, node string, action string) (decision bool) {
	account, ok := ctxAccount.(string)
	if !ok {
		return false
	}
	log := s.Log.Named("checkPerm").With(zap.String("node", node), zap.String("account", account), zap.String("action", action))

	decision, err := s.Repo.IsAuthorized(context.TODO(), node, account, action)
	if err != nil {
		log.Debug("Permission checked", zap.Bool("decision", false))
	}
	log.Debug("Permission checked", zap.Bool("decision", decision))
	return decision
}

func (s *Server) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {
	id, err := s.Repo.CreateObject(ctx, request.GetName(), request.GetParent())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.Object{Uid: id}, nil
}

func (s *Server) ListObjects(ctx context.Context, request *nodepb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {
	directDevices, directObjects, inheritedObjects, err := s.Repo.ListForAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	objects := make([]*nodepb.ObjectList, 0)

	for _, internalObject := range inheritedObjects {
		object := mapObject(internalObject)
		objects = append(objects, object)
	}

	var devices []*nodepb.Device
	if len(directDevices) > 0 {
		for _, directDevice := range directDevices {
			devices = append(devices, &nodepb.Device{
				Uid:  directDevice.UID,
				Name: directDevice.Name,
			})
		}
	}

	// Add direct objects and their devices to the result set, if they are not contained yet
	// Rather inefficient if there's many inherited objects/the slice is long.
	for _, directObject := range directObjects {

		var found bool
		for _, inheritedObject := range inheritedObjects {
			if inheritedObject.Name == directObject.Name {
				found = true
			}
		}

		if !found {
			objects = append(objects, mapObject(directObject))
		}

	}

	return &nodepb.ListObjectsResponse{
		Objects: objects,
		Devices: devices,
	}, nil
}

func mapObject(o ObjectList) *nodepb.ObjectList {
	objects := make([]*nodepb.ObjectList, 0)
	if len(o.Contains) > 0 {
		for _, v := range o.Contains {
			object := mapObject(v)
			objects = append(objects, object)

		}
	}

	var devices []*nodepb.Device
	for _, device := range o.ContainsDevice {
		devices = append(devices, &nodepb.Device{
			Uid:  device.UID,
			Name: device.Name,
		})
	}

	res := &nodepb.ObjectList{
		Uid:     o.UID,
		Name:    o.Name,
		Objects: objects,
		Devices: devices,
	}

	return res
}
