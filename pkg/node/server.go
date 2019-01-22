package node

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

type ObjectController struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
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

func (s *ObjectController) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {
	id, err := s.Repo.CreateObject(ctx, request.GetName(), request.GetParent())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.Object{Uid: id}, nil
}

func (s *ObjectController) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {
	err = s.Repo.DeleteObject(ctx, request.GetUid(), request.GetParentUid())
	if err != nil {
		return nil, err
	}
	return &nodepb.DeleteObjectResponse{}, nil
}

func (s *ObjectController) ListObjects(ctx context.Context, request *nodepb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {
	directDevices, directObjects, inheritedObjects, err := s.Repo.ListForAccount(ctx, request.GetAccount())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	objects := make([]*nodepb.Object, 0)

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

func mapObject(o ObjectList) *nodepb.Object {
	objects := make([]*nodepb.Object, 0)
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

	res := &nodepb.Object{
		Uid:     o.UID,
		Name:    o.Name,
		Objects: objects,
		Devices: devices,
	}

	return res
}
