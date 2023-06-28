/*
Copyright © 2021-2023 Infinite Devices GmbH

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
package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/infinimesh/infinimesh/pkg/sessions"
	infinimesh "github.com/infinimesh/infinimesh/pkg/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	log *zap.Logger
	rdb *redis.Client

	SIGNING_KEY []byte
)

func SetContext(logger *zap.Logger, _rdb *redis.Client, key []byte) {
	log = logger.Named("JWT")
	rdb = _rdb

	SIGNING_KEY = key
	log.Debug("Context set", zap.ByteString("signing_key", key))
}

func MakeToken(account string) (string, error) {
	claims := jwt.MapClaims{}
	claims[infinimesh.INFINIMESH_ACCOUNT_CLAIM] = account
	claims[infinimesh.INFINIMESH_ROOT_CLAIM] = 4
	claims[infinimesh.INFINIMESH_NOSESSION_CLAIM] = true

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SIGNING_KEY)
}

func JWT_AUTH_INTERCEPTOR(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	l := log.Named("Interceptor")
	l.Debug("Invoked", zap.String("method", info.FullMethod))

	if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health/") {
		return handler(ctx, req)
	}

	// Middleware selector
	var middleware func(context.Context) (context.Context, error)
	switch {
	case info.FullMethod == "/infinimesh.node.DevicesService/GetByToken":
		middleware = JwtDeviceAuthMiddleware
	case strings.HasPrefix(info.FullMethod, "/infinimesh.node.ShadowService/"):
		middleware = JwtDeviceAuthMiddleware
	default:
		middleware = JwtStandardAuthMiddleware
	}

	ctx, err := middleware(ctx)
	if info.FullMethod != "/infinimesh.node.AccountsService/Token" && err != nil {
		return nil, err
	}

	go handleLogActivity(ctx)

	return handler(ctx, req)
}

func JwtStandardAuthMiddleware(ctx context.Context) (context.Context, error) {
	l := log.Named("StandardAuthMiddleware")
	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		l.Debug("Error extracting token", zap.Any("error", err))
		return ctx, err
	}

	token, err := validateToken(tokenString)
	if err != nil {
		return ctx, err
	}
	log.Debug("Validated token", zap.Any("claims", token))

	account := token[infinimesh.INFINIMESH_ACCOUNT_CLAIM]
	if account == nil {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: no requestor ID")
	}
	uuid, ok := account.(string)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: requestor ID isn't string")
	}

	session := token[infinimesh.INFINIMESH_SESSION_CLAIM]
	if session == nil {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: no session ID")
	}
	sid, ok := session.(string)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: session ID isn't string")
	}

	// Check if session is valid
	if err := sessions.Check(rdb, uuid, sid); err != nil {
		log.Debug("Session check failed", zap.Any("error", err))
		return ctx, status.Error(codes.Unauthenticated, "Session is expired, revoked or invalid")
	}

	var exp int64
	if token["exp"] != nil {
		exp = int64(token["exp"].(float64))
	}

	ctx = context.WithValue(ctx, infinimesh.InfinimeshAccountCtxKey, uuid)
	ctx = context.WithValue(ctx, infinimesh.InfinimeshSessionCtxKey, sid)
	ctx = context.WithValue(ctx, infinimesh.ContextKey("exp"), exp)

	ctx = metadata.AppendToOutgoingContext(ctx, infinimesh.INFINIMESH_ACCOUNT_CLAIM, uuid)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+tokenString)

	if root := token[infinimesh.INFINIMESH_ROOT_CLAIM]; root != nil {
		r, ok := root.(bool)
		if ok {
			ctx = context.WithValue(ctx, infinimesh.InfinimeshRootCtxKey, r)
		}
	}

	return ctx, nil
}

func JwtDeviceAuthMiddleware(ctx context.Context) (context.Context, error) {
	l := log.Named("DeviceAuthMiddleware")
	tokenString, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		l.Error("Error extracting token", zap.Any("error", err))
		return nil, err
	}

	token, err := validateToken(tokenString)
	if err != nil {
		return nil, err
	}
	log.Debug("Validated token", zap.Any("claims", token))

	devices := token[infinimesh.INFINIMESH_DEVICES_CLAIM]
	if devices == nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid token format: no devices scope")
	}

	ipool, ok := devices.([]interface{})
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Invalid token format: devices scope isn't a slice")
	}

	pool := make([]string, len(ipool))
	for i, el := range ipool {
		pool[i], ok = el.(string)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token format: element %d is not a string", i)
		}
	}
	ctx = context.WithValue(ctx, infinimesh.InfinimeshDevicesCtxKey, pool)

	post := false
	ipost := token[infinimesh.INFINIMESH_POST_STATE_ALLOWED_CLAIM]
	if ipost != nil {
		post, ok = ipost.(bool)
		if !ok {
			post = false
		}
	}
	ctx = context.WithValue(ctx, infinimesh.InfinimeshPostAllowedCtxKey, post)

	return ctx, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Unexpected signing method: %v", t.Header["alg"])
		}
		return SIGNING_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Cannot Validate Token")
}

func handleLogActivity(ctx context.Context) {
	sid_ctx := ctx.Value(infinimesh.InfinimeshSessionCtxKey)
	if sid_ctx == nil {
		return
	}

	sid := sid_ctx.(string)
	req := ctx.Value(infinimesh.InfinimeshAccountCtxKey).(string)
	exp := ctx.Value(infinimesh.ContextKey("exp")).(int64)

	if err := sessions.LogActivity(rdb, req, sid, exp); err != nil {
		log.Warn("Error logging activity", zap.Any("error", err))
	}
}
