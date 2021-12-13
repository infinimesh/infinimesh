//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package node

import (
	"context"

	"github.com/dgraph-io/dgo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slntopp/infinimesh/pkg/node/nodepb"
)

//ObjectController is a Data type for Object Controller file
type ObjectController struct {
	nodepb.UnimplementedObjectServiceServer

	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
}

//CreateObject is a method for creating objects in heirarchy
func (s *ObjectController) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {

	log := s.Log.Named("Create Object Controller")
	//Added logging
	log.Info("Function Invoked",
		zap.String("Name", request.Name),
		zap.String("Kind", request.Kind),
		zap.String("Namespace", request.Namespaceid),
		zap.String("Parent", request.Parent))

	//Initialize the Account Controller with Namespace controller data
	a.Repo = s.Repo
	a.Log = s.Log

	//Get metadata from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	if request.Name == "" {
		return nil, status.Error(codes.FailedPrecondition, "Please provide Object name")
	}

	// If a parent is given, we need permission on the parent. otherwise, we need permission on the namespace as it's created without a parent
	var authorized bool
	if request.Parent != "" {
		resp, err := a.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
			Node:    request.Parent,
			Account: requestorID,
			Action:  nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, err
		}
		authorized = resp.Decision.GetValue()
	} else {
		resp, err := a.IsAuthorizedNamespace(ctx, &nodepb.IsAuthorizedNamespaceRequest{
			Namespaceid: request.Namespaceid,
			Account:     requestorID,
			Action:      nodepb.Action_WRITE,
		})
		if err != nil {
			return nil, err
		}

		authorized = resp.Decision.GetValue()
	}

	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "The account does not have permission to create Objects")
	}

	id, err := s.Repo.CreateObject(ctx, request.GetName(), request.GetParent(), request.GetKind(), request.GetNamespaceid())
	if err != nil {
		//Added logging
		log.Error("Failed to create Object", zap.String("Name", request.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Object Created", zap.String("Name", request.Name), zap.String("Object Id", id))
	return &nodepb.Object{Uid: id, Name: request.GetName()}, nil
}

//DeleteObject is a method for deleting objects in heirarchy
func (s *ObjectController) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {

	log := s.Log.Named("Delete Object Controller")
	//Added logging
	log.Info("Function Invoked", zap.String("Account", request.Uid))

	//Initialize the Account Controller with Namespace controller data
	a.Repo = s.Repo
	a.Log = s.Log

	//Get metadata from context and perform validation
	_, requestorID, err := Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	resp, err := a.IsAuthorized(ctx, &nodepb.IsAuthorizedRequest{
		Node:    request.GetUid(),
		Account: requestorID,
		Action:  nodepb.Action_WRITE,
	})
	if err != nil {
		return nil, err
	}

	if !resp.Decision.GetValue() {
		return nil, status.Error(codes.PermissionDenied, "The Account does not have permission to access the resource")
	}

	err = s.Repo.DeleteObject(ctx, request.GetUid())
	if err != nil {
		//Added logging
		log.Error("Failed to delete object", zap.Error(err))
		return nil, err
	}

	//Added Logging
	log.Info("Delete Object successful")
	return &nodepb.DeleteObjectResponse{}, nil
}

//ListObjects is a method for listing objects in heirarchy
func (s *ObjectController) ListObjects(ctx context.Context, request *nodepb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {

	log := s.Log.Named("List Objects Controller")
	//Added logging
	log.Debug("Function Invoked",
		zap.String("Account", request.Account),
		zap.String("Namespace", request.Namespace),
		zap.Bool("Recurse", request.Recurse))

	//Initialize the Account Controller with Namespace controller data
	a.Repo = s.Repo
	a.Log = s.Log

	//Get metadata from context and perform validation
	_, _, err = Validation(ctx, log)
	if err != nil {
		return nil, err
	}

	objects, err := s.Repo.ListForAccount(ctx, request.Account, request.Namespace, request.Recurse)
	if err != nil {
		//Added logging
		log.Error("Failed to list Objects for the Account", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Debug("List Objects successful")
	return &nodepb.ListObjectsResponse{
		Objects: objects,
	}, nil
}
