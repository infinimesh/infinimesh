package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/infinimesh/proto/node/nodeconnect"
	"go.uber.org/zap"
	"net/http"
)

type OrgAccess struct {
	Namespace string `yaml:"ns"`
	Level     int32  `yaml:"level"`
}

type Config struct {
	ClientId            string               `yaml:"client_id"`
	ClientSecret        string               `yaml:"client_secret"`
	RedirectUrl         string               `yaml:"redirect_url"`
	Scopes              []string             `yaml:"scopes"`
	State               string               `yaml:"state"`
	ApiUrl              string               `yaml:"api_url"`
	AuthUrl             string               `yaml:"auth_url"`
	TokenUrl            string               `yaml:"token_url"`
	OrganizationMapping map[string]OrgAccess `yaml:"organization_mapping"`
}

type OAuthService interface {
	Register(map[string]Config, nodeconnect.AccountsServiceClient, nodeconnect.NamespacesServiceClient, string)
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

func (s *oauthService) Register(router *mux.Router, configs map[string]Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, token string) {
	log := s.log.Named("Register")
	for key, val := range configs {
		registrar, ok := Registrars[key]
		if !ok {
			log.Warn("No such auth type in config", zap.String("type", key))
			continue
		}
		registrar(s.log, router, &val, accClient, nsClient, token)
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
