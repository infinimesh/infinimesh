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
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/grpchealth"
	"github.com/bufbuild/connect-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/infinimesh/infinimesh/pkg/handsfree"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	cc "github.com/infinimesh/proto/handsfree/handsfreeconnect"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log *zap.Logger

	port string

	redisHost string

	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")
	viper.SetDefault("REDIS_HOST", "redis:6379")

	port = viper.GetString("PORT")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))

	redisHost = viper.GetString("REDIS_HOST")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Connecting to Redis", zap.String("URL", redisHost))
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0, // use default DB
	})
	log.Info("Redis connection established")

	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
			h.ServeHTTP(w, r)
		})
	})

	auth.SetContext(log, rdb, SIGNING_KEY)
	authInterceptor := auth.NewAuthInterceptor(log, rdb, SIGNING_KEY)

	interceptors := connect.WithInterceptors(authInterceptor)

	handsfreeServer := handsfree.NewHandsfreeServer(log)
	path, handler := cc.NewHandsfreeServiceHandler(handsfreeServer, interceptors)
	router.PathPrefix(path).Handler(handler)

	checker := grpchealth.NewStaticChecker()
	path, handler = grpchealth.NewHandler(checker)
	router.PathPrefix(path).Handler(handler)

	host := fmt.Sprintf("0.0.0.0:%s", port)

	handler = cors.New(cors.Options{
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"*", "Connect-Protocol-Version"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
	}).Handler(h2c.NewHandler(router, &http2.Server{}))

	log.Info("Serving", zap.String("host", host))
	err := http.ListenAndServe(host, handler)
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
