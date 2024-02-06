package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/infinimesh/infinimesh/pkg/oauth/config"
	"github.com/infinimesh/infinimesh/pkg/oauth/handlers"
	"github.com/infinimesh/proto/node/nodeconnect"
	"go.uber.org/zap"
	"net/http"
)

type Registrar interface {
	Register(*zap.Logger, *mux.Router, *config.Config, nodeconnect.AccountsServiceClient, nodeconnect.NamespacesServiceClient, string)
}

var Registrars = map[string]Registrar{
	"github": &handlers.GithubRegistrar{},
}

type OAuthService interface {
	Register(map[string]config.Config, nodeconnect.AccountsServiceClient, nodeconnect.NamespacesServiceClient, string)
	Run(string, []string)
}

type oauthService struct {
	log *zap.Logger

	registeredProviders []string
}

func NewOauthService(log *zap.Logger) *oauthService {
	return &oauthService{
		log: log,
	}
}

func (s *oauthService) Register(router *mux.Router, configs map[string]config.Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, token string) {
	log := s.log.Named("Register")
	for key, val := range configs {
		registrar, ok := Registrars[key]
		if !ok {
			log.Warn("No such auth type in config", zap.String("type", key))
			continue
		}
		registrar.Register(s.log, router, &val, accClient, nsClient, token)
		s.registeredProviders = append(s.registeredProviders, key)
	}

	router.HandleFunc("/oauth/providers", func(w http.ResponseWriter, r *http.Request) {
		marshal, err := json.Marshal(s.registeredProviders)
		if err != nil {
			s.log.Error("Failed to marshal providers", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get providers. %s", err.Error())))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(marshal)
	}).Methods(http.MethodGet)
}
