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
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"strconv"

	jwt "github.com/golang-jwt/jwt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"encoding/base64"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	inflog "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
	"robpike.io/filter"
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
	if info.FullMethod == "/infinimesh.api.Accounts/Token" {
		return handler(ctx, req)
	}

	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return jwtSigningSecret, nil
	})
	if err != nil {
		return ctx, err
	}

	if !token.Valid {
		return ctx, errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Debug("GRPC API Server: Validated token Function Invoked", zap.Any("Claims", claims))

		if accountID, ok := claims[accountIDClaim]; ok {

			if accountIDStr, ok := accountID.(string); ok {
				//Added the requestor account id to context metadata so that it can be passed on to the server
				ctx = metadata.AppendToOutgoingContext(ctx, "requestorid", accountIDStr)

				resp, err := accountClient.GetAccount(context.Background(), &nodepb.GetAccountRequest{Id: accountIDStr})
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
				}

				if !resp.Enabled {
					return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Account is disabled"))
				}

				ctx = context.WithValue(ctx, accountIDClaim, accountID)

				if restricted, ok := claims[tokenRestrictedClaim]; ok && restricted.(bool) {
					log.Info("Token is restricted", zap.Any("restricted", restricted))

					fullMethod := strings.Split(info.FullMethod, "/")
					reqNS, reqMethod := fullMethod[1], fullMethod[2]
					for ns, ids := range claims {
						if reqNS == ns {
							idSet := make(map[string]bool)
							if ids != nil {
								for _, id := range ids.([]interface{}) {
									idSet[id.(string)] = true
								}
							}
							if reqMethod == "List" {
								r, err := handler(ctx, req)
								if err != nil {
									return r, err
								}
								if ids != nil {
									switch reqNS {
									case "infinimesh.api.Devices":
										res := r.(*registrypb.ListResponse)
										res.Devices = filter.Choose(res.Devices, func(el *registrypb.Device) bool { return idSet[el.Id] }).([]*registrypb.Device)
										r = res
									case "infinimesh.api.Accounts":
										res := r.(*nodepb.ListAccountsResponse)
										res.Accounts = filter.Choose(res.Accounts, func(el *nodepb.Account) bool { return idSet[el.Uid] }).([]*nodepb.Account)
										r = res
									case "infinimesh.api.Namespaces":
										res := r.(*nodepb.ListNamespacesResponse)
										res.Namespaces = filter.Choose(res.Namespaces, func(el *nodepb.Namespace) bool { return idSet[el.Id] }).([]*nodepb.Namespace)
										r = res
									case "infinimesh.api.Objects":
										res := r.(*nodepb.ListObjectsResponse)
										res.Objects = filter.Choose(res.Objects, func(el *nodepb.Object) bool { return idSet[el.Uid] }).([]*nodepb.Object)
										r = res
									}
								}
								return r, err
							} else if reqMethod == "Get" {
								if ids == nil || idSet[req.(map[string]interface{})["id"].(string)] {
									return handler(ctx, req)
								}
							}
							return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Method is restricted"))
						}
					}
					return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Method is restricted"))
				}

				return handler(ctx, req)
			}

		}
		log.Info("Token does not contain account id field", zap.Any("token", token))
	}

	return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
}

var jwtAuth = func(ctx context.Context) (context.Context, error) {
	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Debug("Extracted bearer token", zap.String("token", tokenString))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return jwtSigningSecret, nil
	})
	if err != nil {
		return ctx, err
	}

	if !token.Valid {
		return ctx, errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Info("Validated token", zap.Any("claims", claims))

		if accountID, ok := claims[accountIDClaim]; ok {

			if accountIDStr, ok := accountID.(string); ok {
				resp, err := accountClient.GetAccount(context.Background(), &nodepb.GetAccountRequest{Id: accountIDStr})
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
				}

				if !resp.Enabled {
					return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Account is disabled"))
				}

				newCtx := context.WithValue(ctx, accountIDClaim, accountID)
				return newCtx, nil
			}

		}
		log.Info("Token does not contain account id field", zap.Any("token", token))
	}

	return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Failed to validate token"))
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
