package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cache "github.com/patrickmn/go-cache"

	"github.com/infinimesh/infinimesh/pkg/node/nodepb"
)

var (
	signingSecret  []byte
	accountIDClaim = "account_id"

	nodeserver string

	log *zap.Logger

	accountService nodepb.AccountServiceClient

	userCache *cache.Cache
)

func init() {
	log, _ = zap.NewDevelopment()

	viper.SetDefault("NODE_HOST", "localhost:8082")
	viper.AutomaticEnv()

	nodeserver = viper.GetString("NODE_HOST")

	{
		b64SignSecret := viper.GetString("JWT_SIGNING_KEY")
		if b64SignSecret == "" {
			panic("Invalid signing secret")
		}
		s, err := base64.StdEncoding.DecodeString(b64SignSecret)
		if err != nil {
			panic("Failed to base64 decode sign secret")
		}
		signingSecret = s
	}

	userCache = cache.New(1*time.Second, 5*time.Second)
}

func main() {
	u, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)

	c, err := grpc.Dial(nodeserver, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	accountService = nodepb.NewAccountServiceClient(c)

	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := ParseJWT(cookie.Value)
		if err != nil {
			log.Error("Failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var name string
		if cachedUser, ok := userCache.Get(userID); ok {
			if userString, ok := cachedUser.(string); ok {
				log.Debug("Using cached value", zap.String("userID", userID), zap.String("username", userString))
				name = userString
			}
		} else {
			acc, err := accountService.GetAccount(context.Background(), &nodepb.GetAccountRequest{
				Id: userID,
			})
			if err != nil {
				log.Error("Failed to fetch account", zap.String("userID", userID), zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			name = acc.Name

			_ = userCache.Add(userID, acc.Name, time.Minute*1)
		}

		r.Header.Set("X-WEBAUTH-USER", name)
		p.ServeHTTP(w, r)
	}
}

func ParseJWT(tokenString string) (user string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return signingSecret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Info("Validated token", zap.Any("claims", claims))

		if accountID, ok := claims[accountIDClaim]; ok {
			if accountIDStr, ok := accountID.(string); ok {
				return accountIDStr, nil
			}

		}
		log.Info("Token does not contain account id field", zap.Any("token", token))
	}

	return "", errors.New("Invalid token")

}
