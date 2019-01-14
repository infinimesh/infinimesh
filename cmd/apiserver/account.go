package main

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/infinimesh/infinimesh/pkg/apiserver/apipb"
	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accountAPI struct {
	signingSecret []byte // TODO use asymmetric crypto.
	client        nodepb.AccountServiceClient
}

func (a *accountAPI) Token(ctx context.Context, request *apipb.TokenRequest) (response *apipb.TokenResponse, err error) {
	// TODO check password :D
	resp, err := a.client.Authenticate(ctx, &nodepb.AuthenticateRequest{Username: request.GetUsername(), Password: request.GetPassword()})
	if err != nil {
		return nil, err
	}

	if resp.GetSuccess() {
		if resp.Account == nil {
			return nil, status.Error(codes.Internal, "Failed to check credentials")
		}
		// Issue token
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			accountIDClaim: resp.Account.Uid,
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(a.signingSecret)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed to sign token")
		}

		return &apipb.TokenResponse{Token: tokenString}, nil
	}

	return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
}

func (a *accountAPI) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}
