//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
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

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

//ObjectController is a Data type for Object Controller file
type ObjectController struct {
	Dgraph *dgo.Dgraph
	Log    *zap.Logger

	Repo Repo
}

//CreateObject is a method for creating objects in heirarchy
func (s *ObjectController) CreateObject(ctx context.Context, request *nodepb.CreateObjectRequest) (response *nodepb.Object, err error) {

	log := s.Log.Named("Create Object Controller")
	//Added logging
	log.Info("Create Object Controller", zap.Bool("Function Invoked", true),
		zap.String("Name", request.Name),
		zap.String("Kind", request.Kind),
		zap.String("Namespace", request.Namespace),
		zap.String("Parent", request.Parent))

	id, err := s.Repo.CreateObject(ctx, request.GetName(), request.GetParent(), request.GetKind(), request.GetNamespace())
	if err != nil {
		//Added logging
		log.Error("Create Object Controller", zap.Bool("Failed to create Object", true), zap.String("Name", request.Name), zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("Create Object Controller", zap.Bool("Object Created", true), zap.String("Name", request.Name), zap.String("Object Id", id))
	return &nodepb.Object{Uid: id, Name: request.GetName()}, nil
}

//DeleteObject is a method for deleting objects in heirarchy
func (s *ObjectController) DeleteObject(ctx context.Context, request *nodepb.DeleteObjectRequest) (response *nodepb.DeleteObjectResponse, err error) {

	log := s.Log.Named("Delete Object Controller")
	//Added logging
	log.Info("Delete Object Controller", zap.Bool("Function Invoked", true), zap.String("Account", request.Uid))

	err = s.Repo.DeleteObject(ctx, request.GetUid())
	if err != nil {
		//Added logging
		log.Error("Delete Object Controller", zap.Bool("Failed to delete object", true), zap.Error(err))
		return nil, err
	}

	//Added Logging
	log.Info("Delete Object Controller", zap.Bool("Delete Object successful", true))
	return &nodepb.DeleteObjectResponse{}, nil
}

//ListObjects is a method for listing objects in heirarchy
func (s *ObjectController) ListObjects(ctx context.Context, request *nodepb.ListObjectsRequest) (response *nodepb.ListObjectsResponse, err error) {

	log := s.Log.Named("List Objects Controller")
	//Added logging
	log.Info("List Objects Controller", zap.Bool("Function Invoked", true),
		zap.String("Account", request.Account),
		zap.String("Namespace", request.Namespace),
		zap.Bool("Recurse", request.Recurse))

	objects, err := s.Repo.ListForAccount(ctx, request.Account, request.Namespace, request.Recurse)
	if err != nil {
		//Added logging
		log.Error("List Objects Controller", zap.Bool("Failed to list accounts", true), zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	//Added logging
	log.Info("List Objects Controller", zap.Bool("List Objects successful", true))
	return &nodepb.ListObjectsResponse{
		Objects: objects,
	}, nil
}
