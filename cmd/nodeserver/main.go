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

package main

import (
	"net"
	"syscall"

	"os"
	"os/signal"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

var (
	dgraphURL string
	port      = ":8082"
)

func init() {
	viper.SetDefault("DGRAPH_HOST", "localhost:9080")
	viper.AutomaticEnv()

	dgraphURL = viper.GetString("DGRAPH_HOST")
}

func main() {
	log, err := log.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Connecting to dgraph", zap.String("URL", dgraphURL))
	conn, err := grpc.Dial(dgraphURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to dgraph", zap.Error(err))
	}
	log.Info("Connected to dgraph")
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}

	srv := grpc.NewServer()

	repo := dgraph.NewDGraphRepo(dg)

	objectController := &node.ObjectController{
		Repo:   repo,
		Dgraph: dg,
		Log:    log.Named("server"),
	}

	accountController := &node.AccountController{
		Repo:   repo,
		Dgraph: dg,
		Log:    log.Named("accountController"),
	}

	namespaceController := &node.NamespaceController{
		Repo: repo,
	}

	nodepb.RegisterObjectServiceServer(srv, objectController)
	nodepb.RegisterAccountServiceServer(srv, accountController)
	nodepb.RegisterNamespacesServer(srv, namespaceController)
	reflection.Register(srv)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		log.Info("Starting gRPC server", zap.String("address", port))
		if err := srv.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	<-signals
	log.Info("Stopping gRPC server")
	srv.GracefulStop()
	log.Info("Stopped gRPC server")
}
