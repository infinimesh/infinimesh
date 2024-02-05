package main

import (
	"github.com/gorilla/mux"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	logger "github.com/infinimesh/infinimesh/pkg/log"
	"github.com/infinimesh/infinimesh/pkg/oauth"
	"github.com/infinimesh/infinimesh/pkg/shared/auth"
	"github.com/infinimesh/proto/node/nodeconnect"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strings"
)

var (
	log *zap.Logger

	port                 string
	accountsConnection   string
	namespacesConnection string
	signingKey           []byte
	corsAllowedIn        string
	configs              map[string]oauth.Config
)

func init() {
	viper.AutomaticEnv()

	log = logger.NewLogger()
	viper.SetDefault("PORT", "80")
	viper.SetDefault("REGISTRY", "repo:8000")
	viper.SetDefault("NAMESPACES", "repo:8000")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")
	viper.SetDefault("CORS_ALLOWED", []string{"*"})

	port = viper.GetString("PORT")
	accountsConnection = viper.GetString("REGISTRY")
	namespacesConnection = viper.GetString("NAMESPACES")
	signingKey = []byte(viper.GetString("SIGNING_KEY"))
	corsAllowedIn = viper.GetString("CORS_ALLOWED")

	file, err := os.ReadFile("oauth2_config.yaml")
	if err != nil {
		log.Fatal("Failed to read config", zap.Error(err))
	}
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		log.Fatal("Failed to parse config", zap.Error(err))
	}
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	router := mux.NewRouter()
	accClient := nodeconnect.NewAccountsServiceClient(http.DefaultClient, accountsConnection)
	nsClient := nodeconnect.NewNamespacesServiceClient(http.DefaultClient, namespacesConnection)

	interceptor := auth.NewAuthInterceptor(log, nil, nil, signingKey)
	token, err := interceptor.MakeToken(schema.ROOT_ACCOUNT_KEY)
	if err != nil {
		log.Fatal("Failed to create token", zap.Error(err))
	}

	cors := strings.Split(corsAllowedIn, ",")

	service := oauth.NewOauthService(log, router)
	service.Register(configs, accClient, nsClient, token)
	service.Run(port, cors)
}
