/*
Copyright © 2021-2022 Infinite Devices GmbH Nikita Ivanovski info@slnt-opp.xyz

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
package handsfree

import (
	"time"

	pb "github.com/infinimesh/proto/handsfree"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandsfreeServer struct {
	pb.UnimplementedHandsfreeServiceServer

	log *zap.Logger
	db  map[string](chan string)
}

func NewHandsfreeServer(log *zap.Logger) *HandsfreeServer {
	return &HandsfreeServer{
		log: log.Named("HandsfreeServer"), db: make(map[string]chan string),
	}
}

func (s *HandsfreeServer) GetToken(req *pb.ConnectionRequest, srv pb.HandsfreeService_GetTokenServer) error {
	log := s.log.Named("GetToken")
	log.Debug("Request received", zap.String("app", req.GetAppId()))

	if req.GetAppId() == "" {
		return status.Error(codes.InvalidArgument, "Application ID must be present upon connection")
	}

	hash, err := GenerateCodeLong(req.GetAppId())
	if err != nil {
		log.Warn("Error during generating hash", zap.Error(err))
		return status.Error(codes.Internal, "Internal Error while generating code. Check your Application ID validity")
	}

	code, err := ShortenToFit(hash, s.db)
	if err != nil {
		log.Warn("Error during shortening the hash to make code", zap.Error(err))
		return status.Error(codes.Internal, "Internal Error while generating code. Check your Application ID validity and have some mercy on our servers")
	}

	s.db[code] = make(chan string)

	err = srv.Send(&pb.ControlPacket{
		Code: pb.Code_AUTH_CODE, Payload: code,
	})
	if err != nil {
		return nil
	}

	refresh := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-refresh.C:
			hash, err := GenerateCodeLong(req.GetAppId())
			if err != nil {
				log.Warn("Error during generating hash", zap.Error(err))
				return status.Error(codes.Internal, "Internal Error while generating code. Check your Application ID validity")
			}

			code, err = ShortenToFit(hash, s.db)
			if err != nil {
				log.Warn("Error during shortening the hash to make code", zap.Error(err))
				return status.Error(codes.Internal, "Internal Error while generating code. Check your Application ID validity and have some mercy on our servers")
			}
			s.db[code] = make(chan string)

			err = srv.Send(&pb.ControlPacket{
				Code: pb.Code_AUTH_CODE, Payload: code,
			})
			if err != nil {
				return nil
			}
		case token := <-s.db[code]:
			srv.Send(&pb.ControlPacket{
				Code: pb.Code_TOKEN, Payload: token,
			})

			refresh.Stop()
			return nil
		}
	}
}
