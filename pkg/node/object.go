package node

import (
	"context"

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

func (s *ObjectController) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {
	id, err := s.Repo.CreateObject(ctx, request.GetName(), request.GetParent(), request.GetKind(), request.GetNamespace())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.Object{Uid: id, Name: request.GetName()}, nil
}

func (s *ObjectController) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {
	err = s.Repo.DeleteObject(ctx, request.GetUid())
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
		// object := mapObject(internalObject)
		objects = append(objects, internalObject)
	}

	var devices []*nodepb.Device
	if len(directDevices) > 0 {
		for _, directDevice := range directDevices {
			devices = append(devices, directDevice)
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
			// objects = append(objects, mapObject(directObject))
			objects = append(objects, directObject)
		}

	}

	return &nodepb.ListObjectsResponse{
		Objects: objects,
		Devices: devices,
	}, nil
}
