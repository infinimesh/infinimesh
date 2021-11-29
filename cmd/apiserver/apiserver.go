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
	"context"
	"errors"
	"fmt"
	"net"

	"strconv"

	jwt "github.com/golang-jwt/jwt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"encoding/base64"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	inflog "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
)

const (
	accountIDClaim       = "account_id"
	tokenRestrictedClaim = "restricted"
	expiresAt            = "exp"
)

var (
	registryHost     string
	shadowHost       string
	nodeHost         string
	jwtSigningSecret []byte
	port             int

	accountClient nodepb.AccountServiceClient

	log *zap.Logger
)


var jwtAuthInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Info("JWT: Request received")

	if info.FullMethod == "/infinimesh.api.Accounts/Token" {
		return handler(ctx, req)
	}

	ctx, err := jwtAuth(ctx)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

var jwtAuth = func(ctx context.Context) (context.Context, error) {
	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Debug("Extracted bearer token", zap.String("token", tokenString))

	var account_id_value interface{}

	account_id_value, err = validateSelfSignedToken(tokenString)
	if err != nil {
		return ctx, err
	}

	if account_id, ok := account_id_value.(string); ok {
		return authorizeAccount(ctx, account_id)
	}

	return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
}

func validateSelfSignedToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return jwtSigningSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Info("Validated token", zap.Any("claims", claims))
		if accountID, ok := claims[accountIDClaim]; ok {
			return accountID, nil
		}
	}

	return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
}

func authorizeAccount(ctx context.Context, account_id string) (context.Context, error) {
	resp, err := accountClient.GetAccount(context.Background(), &nodepb.GetAccountRequest{Id: account_id})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
	}

	if !resp.Enabled {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Account is disabled"))
	}

	newCtx := context.WithValue(ctx, accountIDClaim, account_id)
	return newCtx, nil
}

func init() {
	viper.SetDefault("REGISTRY_HOST", "device-registry:8080")
	viper.SetDefault("SHADOW_HOST", "shadow-api:8096")
	viper.SetDefault("NODE_HOST", "nodeserver:8082")
	viper.SetDefault("PORT", 8080)
	viper.AutomaticEnv()

	registryHost = viper.GetString("REGISTRY_HOST")
	shadowHost = viper.GetString("SHADOW_HOST")
	nodeHost = viper.GetString("NODE_HOST")
	port = viper.GetInt("PORT")

	b64SignSecret := viper.GetString("JWT_SIGNING_KEY")
	if b64SignSecret == "" {
		panic("Invalid signing secret")
	}

	s, err := base64.StdEncoding.DecodeString(b64SignSecret)
	if err != nil {
		panic("Failed to base64 decode sign secret")
	}

	jwtSigningSecret = s

	logger, err := inflog.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	log = logger

}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(jwtAuth)),
		grpc.UnaryInterceptor(jwtAuthInterceptor),
	)

	registryConn, err := grpc.Dial(registryHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	devicesClient := registrypb.NewDevicesClient(registryConn)

	shadowConn, err := grpc.Dial(shadowHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	shadowClient := shadowpb.NewShadowsClient(shadowConn)

	nodeConn, err := grpc.Dial(nodeHost, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	accountClient = nodepb.NewAccountServiceClient(nodeConn)
	objectClient := nodepb.NewObjectServiceClient(nodeConn)

	namespaceClient := nodepb.NewNamespacesClient(nodeConn)

	//Added logging
	log.Info("GRPC API Server: Starting GRPC Service")

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient, accountClient: accountClient})
	apipb.RegisterStatesServer(srv, &shadowAPI{client: shadowClient, accountClient: accountClient})
	apipb.RegisterAccountsServer(srv, &accountAPI{client: accountClient, signingSecret: jwtSigningSecret})
	apipb.RegisterObjectsServer(srv, &objectAPI{objectClient: objectClient, accountClient: accountClient})
	apipb.RegisterNamespacesServer(srv, &namespaceAPI{client: namespaceClient, accountClient: accountClient})
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		//Added logging
		log.Error("GRPC API Server: Unable to start GRPC Service", zap.Error(err))
		panic(err)
	}

	//Added logging
	log.Info("GRPC API Server: GRPC Service Started")
	reflection.Register(srv)

	err = srv.Serve(listener)
	if err != nil {
		//Added logging
		log.Error("GRPC API Server: Unable to start GRPC Service", zap.Error(err))
		panic(err)
	}

}
