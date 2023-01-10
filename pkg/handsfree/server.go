/*
Copyright © 2021-2023 Infinite Devices GmbH Nikita Ivanovski info@slnt-opp.xyz

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
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	pb "github.com/infinimesh/proto/handsfree"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type HandsfreeServer struct {
	pb.UnimplementedHandsfreeServiceServer

	log *zap.Logger
	db  map[string](chan []string)
}

func NewHandsfreeServer(log *zap.Logger) *HandsfreeServer {
	return &HandsfreeServer{
		log: log.Named("HandsfreeServer"), db: make(map[string]chan []string),
	}
}

func GenerateCode[T any](db map[string]T) (r string) {
begin:
	for i := 0; i < 6; i++ {
		n := rand.Intn(16)
		r += strconv.FormatInt(int64(n), 16)
	}
	if _, exists := db[r]; exists {
		r = ""
		goto begin
	}
	return r
}

func (s *HandsfreeServer) Send(ctx context.Context, req *pb.ControlPacket) (*pb.ControlPacket, error) {
	log := s.log.Named("Send")
	if len(req.GetPayload()) < 2 {
		return nil, status.Error(codes.InvalidArgument, "Payload must consist of code and actual payload")
	}
	log.Debug("Request received", zap.Strings("payload", req.GetPayload()))

	code := strings.ToLower(req.GetPayload()[0])

	ch, ok := s.db[code]
	if !ok {
		return nil, status.Error(codes.NotFound, "No App's awaiting with this code")
	}

	ch <- req.GetPayload()[1:]

	return &pb.ControlPacket{
		Code: pb.Code_SUCCESS,
	}, nil
}

func (s *HandsfreeServer) Connect(req *pb.ConnectionRequest, srv pb.HandsfreeService_ConnectServer) error {
	log := s.log.Named("Connect")
	log.Debug("Request received", zap.String("app", req.GetAppId()))

	if req.GetAppId() == "" {
		return status.Error(codes.InvalidArgument, "Application ID must be present upon connection")
	}

	code := GenerateCode(s.db)
	s.db[code] = make(chan []string)

	err := srv.Send(&pb.ControlPacket{
		Code: pb.Code_AUTH, Payload: []string{code},
	})
	if err != nil {
		return nil
	}

	refresh := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-refresh.C:
			code := GenerateCode(s.db)
			s.db[code] = make(chan []string)

			err = srv.Send(&pb.ControlPacket{
				Code: pb.Code_AUTH, Payload: []string{code},
			})
			if err != nil {
				return nil
			}
		case payload := <-s.db[code]:
			srv.Send(&pb.ControlPacket{
				Code: pb.Code_DATA, Payload: payload,
			})

			refresh.Stop()
			return nil
		}
	}
}
