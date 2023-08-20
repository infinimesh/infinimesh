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

	"connectrpc.com/connect"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/infinimesh/infinimesh/pkg/sessions"
	infinimesh "github.com/infinimesh/infinimesh/pkg/shared"
)

type interceptor struct {
	log         *zap.Logger
	rdb         *redis.Client
	signing_key []byte
}

func NewAuthInterceptor(log *zap.Logger, _rdb *redis.Client, signing_key []byte) *interceptor {
	return &interceptor{
		log:         log.Named("AuthInterceptor"),
		rdb:         _rdb,
		signing_key: signing_key,
	}
}

func (i *interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		path := req.Header().Get(":path")
		header := req.Header().Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
		}

		var middleware func(context.Context, []byte, string) (context.Context, error)

		switch {
		case path == "/infinimesh.node.DevicesService/GetByToken":
			middleware = i.ConnectDeviceAuthMiddleware
		case strings.HasPrefix(path, "/infinimesh.node.ShadowService/"):
			middleware = i.ConnectDeviceAuthMiddleware
		default:
			middleware = i.ConnectStandardAuthMiddleware
		}
		i.log.Debug("Authorization Header", zap.String("header", header))

		ctx, err := middleware(ctx, i.signing_key, segments[1])
		if path != "/infinimesh.node.AccountsService/Token" && err != nil {
			return nil, err
		}

		go i.connectHandleLogActivity(ctx)

		return next(ctx, req)
	})
}

func (i *interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	i.log.Debug("WrapStreamingClient")
	return next
}
func (i *interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	i.log.Debug("Setup Wrap Streaming Handler")
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		path := shc.RequestHeader().Get(":path")
		header := shc.RequestHeader().Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
		}

		var middleware func(context.Context, []byte, string) (context.Context, error)

		switch {
		case path == "/infinimesh.node.DevicesService/GetByToken":
			middleware = i.ConnectDeviceAuthMiddleware
		case strings.HasPrefix(path, "/infinimesh.node.ShadowService/"):
			middleware = i.ConnectDeviceAuthMiddleware
		default:
			middleware = i.ConnectStandardAuthMiddleware
		}
		i.log.Debug("Authorization Header", zap.String("header", header))

		ctx, err := middleware(ctx, i.signing_key, segments[1])
		if path != "/infinimesh.node.AccountsService/Token" && err != nil {
			return err
		}

		go i.connectHandleLogActivity(ctx)

		return next(ctx, shc)
	}
}

func (i *interceptor) ConnectStandardAuthMiddleware(ctx context.Context, signingKey []byte, tokenString string) (context.Context, error) {
	token, err := connectValidateToken(signingKey, tokenString)
	if err != nil {
		return ctx, err
	}
	i.log.Debug("Validated token", zap.Any("claims", token))

	account := token[infinimesh.INFINIMESH_ACCOUNT_CLAIM]
	if account == nil {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: no requestor ID")
	}
	uuid, ok := account.(string)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "Invalid token format: requestor ID isn't string")
	}

	if token[infinimesh.INFINIMESH_NOSESSION_CLAIM] == nil {
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
			i.log.Debug("Session check failed", zap.Any("error", err))
			return ctx, status.Error(codes.Unauthenticated, "Session is expired, revoked or invalid")
		}

		ctx = context.WithValue(ctx, infinimesh.InfinimeshSessionCtxKey, sid)
	}

	var exp int64
	if token["exp"] != nil {
		exp = int64(token["exp"].(float64))
	}

	ctx = context.WithValue(ctx, infinimesh.InfinimeshAccountCtxKey, uuid)
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

func (i *interceptor) ConnectDeviceAuthMiddleware(ctx context.Context, signingKey []byte, tokenString string) (context.Context, error) {
	token, err := connectValidateToken(signingKey, tokenString)
	if err != nil {
		return nil, err
	}
	i.log.Debug("Validated token", zap.Any("claims", token))

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

func connectValidateToken(signing_key []byte, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Unexpected signing method: %v", t.Header["alg"])
		}
		return signing_key, nil
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

func (i *interceptor) connectHandleLogActivity(ctx context.Context) {
	sid_ctx := ctx.Value(infinimesh.InfinimeshSessionCtxKey)
	if sid_ctx == nil {
		return
	}

	sid := sid_ctx.(string)
	req := ctx.Value(infinimesh.InfinimeshAccountCtxKey).(string)
	exp := ctx.Value(infinimesh.ContextKey("exp")).(int64)

	if err := sessions.LogActivity(rdb, req, sid, exp); err != nil {
		i.log.Warn("Error logging activity", zap.Any("error", err))
	}
}
