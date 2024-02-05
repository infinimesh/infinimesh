package oauth

import (
	"github.com/infinimesh/proto/node/nodeconnect"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Registrar func(*zap.Logger, *mux.Router, *Config, nodeconnect.AccountsServiceClient, string)

var Registrars = map[string]Registrar{
	"github": RegisterGithub,
}

func RegisterGithub(logger *zap.Logger, router *mux.Router, config *Config, client nodeconnect.AccountsServiceClient, infinimeshToken string) {
	log := logger.Named("GithubOauth")
	log.Info("Init Github Oauth2 handlers")

	githubConfig := &oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  config.RedirectUrl,
		Scopes:       config.Scopes,
	}

	router.HandleFunc("oauth/github/login", GithubLoginHandler(log, githubConfig, config))

	router.HandleFunc("oauth/github/checkout", GithubCheckoutHandler(log, githubConfig, config, client, infinimeshToken))
}
