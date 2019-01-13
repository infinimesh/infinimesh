package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const accountIDClaim = "account_id"

var (
	registryHost     string
	shadowHost       string
	nodeHost         string
	jwtSigningSecret []byte
	port             int
)

func init() {
	viper.SetDefault("REGISTRY_HOST", "device-registry:8080")
	viper.SetDefault("SHADOW_HOST", "shadow-api:8096")
	viper.SetDefault("NODE_HOST", "nodeserver:8096")
	viper.SetDefault("PORT", 8080)
	viper.AutomaticEnv()

	registryHost = viper.GetString("REGISTRY_HOST")
	shadowHost = viper.GetString("SHADOW_HOST")
	nodeHost = viper.GetString("NODE_HOST")
	port = viper.GetInt("PORT")

	jwtSigningSecret = []byte("super secret key")
}

func main() {
	log, err := log.NewProdOrDev()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	jwtAuth := func(ctx context.Context) (context.Context, error) {
		tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		log.Debug("Extracted bearer token", zap.String("token", tokenString))

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
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
				newCtx := context.WithValue(ctx, node.ContextKeyAccount, accountID)
				return newCtx, nil

			} else {
				log.Info("Token does not contain account id field", zap.Any("token", token))
			}
		}

		return ctx, errors.New("Failed to validate token")
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(jwtAuth)),
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
	nodeClient := nodepb.NewAccountServiceClient(nodeConn)

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient})
	apipb.RegisterShadowsServer(srv, &shadowAPI{client: shadowClient})
	apipb.RegisterAccountServer(srv, &accountAPI{client: nodeClient, signingSecret: jwtSigningSecret})
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	reflection.Register(srv)

	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}

}
