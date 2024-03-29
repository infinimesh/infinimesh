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

	"github.com/infinimesh/proto/handsfree/handsfreeconnect"
	"github.com/infinimesh/proto/node/access"
	"github.com/infinimesh/proto/node/nodeconnect"

	"github.com/infinimesh/infinimesh/pkg/sessions"
	infinimesh "github.com/infinimesh/infinimesh/pkg/shared"
)

type AuthInterceptor interface {
	connect.Interceptor

	MakeToken(string) (string, error)
	ConnectStandardAuthMiddleware(context.Context, []byte, string) (context.Context, bool, error)
	ConnectDeviceAuthMiddleware(context.Context, []byte, string) (context.Context, bool, error)

	SetSessionsHandler(sessions.SessionsHandler)
}

type interceptor struct {
	log         *zap.Logger
	rdb         *redis.Client
	jwt         JWTHandler
	sessions    sessions.SessionsHandler
	signing_key []byte
}

type middleware func(context.Context, []byte, string) (context.Context, bool, error)

func NewAuthInterceptor(log *zap.Logger, _rdb *redis.Client, _jwth JWTHandler, signing_key []byte) *interceptor {
	jwth := _jwth
	if jwth == nil {
		jwth = defaultJWTHandler{}
	}

	sessions := sessions.NewSessionsHandlerModule(_rdb).Handler()

	return &interceptor{
		log:         log.Named("AuthInterceptor"),
		rdb:         _rdb,
		jwt:         jwth,
		sessions:    sessions,
		signing_key: signing_key,
	}
}

func (i *interceptor) MakeToken(account string) (string, error) {
	claims := jwt.MapClaims{}
	claims[infinimesh.INFINIMESH_ACCOUNT_CLAIM] = account
	claims[infinimesh.INFINIMESH_ROOT_CLAIM] = 4
	claims[infinimesh.INFINIMESH_NOSESSION_CLAIM] = true

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(i.signing_key)
}

func SelectMiddleware(i *interceptor, procedure string) middleware {
	var middleware middleware

	switch {
	case procedure == nodeconnect.DevicesServiceGetByTokenProcedure:
		middleware = i.ConnectDeviceAuthMiddleware
	case strings.HasPrefix(procedure, "/"+nodeconnect.ShadowServiceName):
		middleware = i.ConnectDeviceAuthMiddleware
	default:
		middleware = i.ConnectStandardAuthMiddleware
	}

	return middleware
}

func (i *interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		procedure := req.Spec().Procedure
		header := req.Header().Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			segments = []string{"", ""}
		}

		i.log.Debug("Authorization Header", zap.String("header", header))

		middleware := SelectMiddleware(i, procedure)

		ctx, log_activity, err := middleware(ctx, i.signing_key, segments[1])
		if procedure != nodeconnect.AccountsServiceTokenProcedure && err != nil {
			return nil, err
		}

		if log_activity {
			go i.connectHandleLogActivity(ctx)
		}

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
		procedure := shc.Spec().Procedure
		header := shc.RequestHeader().Get("Authorization")

		segments := strings.Split(header, " ")
		if len(segments) != 2 {
			segments = []string{"", ""}
		}

		var middleware middleware

		switch {
		case strings.HasPrefix(procedure, "/"+nodeconnect.ShadowServiceName):
			middleware = i.ConnectDeviceAuthMiddleware
		case procedure == handsfreeconnect.HandsfreeServiceConnectProcedure:
			middleware = i.ConnectBlankMiddleware
		default:
			middleware = i.ConnectStandardAuthMiddleware
		}
		i.log.Debug("Authorization Header", zap.String("header", header))

		ctx, log_activity, err := middleware(ctx, i.signing_key, segments[1])
		if err != nil {
			return err
		}

		if log_activity {
			go i.connectHandleLogActivity(ctx)
		}

		return next(ctx, shc)
	}
}

func (i *interceptor) ConnectStandardAuthMiddleware(_ctx context.Context, signingKey []byte, tokenString string) (ctx context.Context, log_activity bool, err error) {
	ctx = _ctx

	log := i.log.Named("StandardAuthMiddleware")

	token, err := connectValidateToken(i.jwt, signingKey, tokenString)
	if err != nil {
		err = status.Error(codes.Unauthenticated, "Invalid token format")
		return
	}
	log.Debug("Validated token", zap.Any("claims", token))

	account := token[infinimesh.INFINIMESH_ACCOUNT_CLAIM]
	if account == nil {
		err = status.Error(codes.Unauthenticated, "Invalid token format: no requestor ID")
		return
	}
	uuid, ok := account.(string)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Invalid token format: requestor ID isn't string")
		return
	}

	if token[infinimesh.INFINIMESH_NOSESSION_CLAIM] == nil {
		log_activity = true
		session := token[infinimesh.INFINIMESH_SESSION_CLAIM]
		if session == nil {
			err = status.Error(codes.Unauthenticated, "Invalid token format: no session ID")
			return
		}
		sid, ok := session.(string)
		if !ok {
			err = status.Error(codes.Unauthenticated, "Invalid token format: session ID isn't string")
			return
		}

		// Check if session is valid
		if err = i.sessions.Check(uuid, sid); err != nil {
			i.log.Debug("Session check failed", zap.Any("error", err))
			err = status.Error(codes.Unauthenticated, "Session is expired, revoked or invalid")
			return
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

	return
}

func (i *interceptor) ConnectBlankMiddleware(_ctx context.Context, signingKey []byte, tokenString string) (ctx context.Context, log_activity bool, err error) {
	return _ctx, false, nil
}

func (i *interceptor) ConnectDeviceAuthMiddleware(_ctx context.Context, signingKey []byte, tokenString string) (ctx context.Context, log_activity bool, err error) {
	log_activity = false
	token, err := connectValidateToken(i.jwt, signingKey, tokenString)
	if err != nil {
		return
	}
	log := i.log.Named("DeviceAuthMiddleware")
	log.Debug("Validated token", zap.Any("claims", token))

	devices := token[infinimesh.INFINIMESH_DEVICES_CLAIM]
	if devices == nil {
		err = status.Error(codes.Unauthenticated, "Invalid token format: no devices scope")
		return
	}

	ipool, ok := devices.(map[string]any)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Invalid token format: devices scope isn't a map")
		return
	}

	pool := make(map[string]access.Level, len(ipool))
	for key, value := range ipool {
		val, ok := value.(float64)
		if !ok {
			err = status.Errorf(codes.Unauthenticated, "Invalid token format: element %v is not a number", value)
			return
		}
		pool[key] = access.Level(val)
	}
	ctx = context.WithValue(_ctx, infinimesh.InfinimeshDevicesCtxKey, pool)

	return
}

func connectValidateToken(jwth JWTHandler, signing_key []byte, tokenString string) (jwt.MapClaims, error) {
	token, err := jwth.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
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

	if err := i.sessions.LogActivity(req, sid, exp); err != nil {
		i.log.Warn("Error logging activity", zap.Any("error", err))
	}
}

func (i *interceptor) SetSessionsHandler(sessions sessions.SessionsHandler) {
	i.sessions = sessions
}
