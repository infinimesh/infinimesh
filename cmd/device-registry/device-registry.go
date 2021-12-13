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

package main

import (
	"fmt"
	"net"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/slntopp/infinimesh/pkg/registry"
	"github.com/slntopp/infinimesh/pkg/registry/registrypb"
	"github.com/slntopp/infinimesh/pkg/repo"

	logger "github.com/slntopp/infinimesh/pkg/log"
)

var (
	port      string
	dgraphURL string
	dbAddr    string
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DGRAPH_HOST", "localhost:9080")
	viper.SetDefault("DB_ADDR2", ":6379")

	dgraphURL = viper.GetString("DGRAPH_HOST")
	port = viper.GetString("PORT")
	dbAddr = viper.GetString("DB_ADDR2")
}

func main() {
	log, err := logger.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	log.Info("Connecting to dgraph", zap.String("URL", dgraphURL))
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to dgraph", zap.Error(err))
	}
	log.Info("Connected to dgraph")
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	rep, err := repo.NewRedisRepo(dbAddr)
	if err != nil {
		log.Fatal("Failed to connect to redis2 with db addr", zap.Error(err))
	}
	repServ := repo.Server{
		Repo: rep,
		Log:  log.Named("RepoController"),
	}
	server := registry.NewServer(dg, repServ)

	server.Log = log.Named("deviceController")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}
	s := grpc.NewServer()
	registrypb.RegisterDevicesServer(s, server)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}
