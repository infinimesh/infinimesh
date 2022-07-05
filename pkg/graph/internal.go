/*
Copyright © 2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graph

import (
	"context"

	pb "github.com/infinimesh/proto/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cred "github.com/infinimesh/infinimesh/pkg/credentials"
)

type InternalService struct {
	pb.UnimplementedInternalServiceServer
}

func (*InternalService) GetLDAPProviders(ctx context.Context, _ *pb.EmptyMessage) (*pb.LDAPProviders, error) {
	if !cred.LDAP_CONFIGURED {
		return nil, status.Error(codes.OK, "LDAP Auth is not configured")
	}

	res := make(map[string]string)
	for key := range cred.LDAP.Providers {
		res[key] = "" // TODO: add title
	}

	return &pb.LDAPProviders{
		Providers: res,
	}, nil
}
