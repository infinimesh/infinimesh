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
	objects, err := s.Repo.ListForAccount(ctx, request.Account, request.Namespace, request.Recurse)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &nodepb.ListObjectsResponse{
		Objects: objects,
	}, nil
}
