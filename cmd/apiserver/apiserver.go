package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"encoding/base64"

	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/infinimesh/infinimesh/pkg/shadow/shadowpb"
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
	viper.SetDefault("NODE_HOST", "nodeserver:8082")
	viper.SetDefault("PORT", 8080)
	viper.AutomaticEnv()

	registryHost = viper.GetString("REGISTRY_HOST")
	shadowHost = viper.GetString("SHADOW_HOST")
	nodeHost = viper.GetString("NODE_HOST")
	port = viper.GetInt("PORT")

	jwtSigningSecret = []byte("super secret key")

	b64SignSecret := viper.GetString("JWT_SIGNING_KEY")
	if b64SignSecret == "" {
		panic("Invalid signing secret")
	}

	s, err := base64.StdEncoding.DecodeString(b64SignSecret)
	if err != nil {
		panic("Failed to base64 decode sign secret")
	}

	jwtSigningSecret = s
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
				newCtx := context.WithValue(ctx, accountIDClaim, accountID)
				return newCtx, nil
			}
			log.Info("Token does not contain account id field", zap.Any("token", token))
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
	accountClient := nodepb.NewAccountServiceClient(nodeConn)
	objectClient := nodepb.NewObjectServiceClient(nodeConn)

	namespaceClient := nodepb.NewNamespaceServiceClient(nodeConn)

	apipb.RegisterDevicesServer(srv, &deviceAPI{client: devicesClient, accountClient: accountClient})
	apipb.RegisterShadowsServer(srv, &shadowAPI{client: shadowClient})
	apipb.RegisterAccountServer(srv, &accountAPI{client: accountClient, signingSecret: jwtSigningSecret})
	apipb.RegisterObjectServiceServer(srv, &objectAPI{objectClient: objectClient, accountClient: accountClient})
	apipb.RegisterNamespaceServiceServer(srv, &namespaceAPI{client: namespaceClient, accountClient: accountClient})
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
