package auth

import "github.com/golang-jwt/jwt/v4"

type JWTHandler interface {
	Parse(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error)
}

type defaultJWTHandler struct{}

func (d defaultJWTHandler) Parse(tokenString string, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
	return jwt.Parse(tokenString, keyFunc, options...)
}
