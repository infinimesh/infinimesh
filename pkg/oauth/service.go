package oauth

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/infinimesh/proto/node/nodeconnect"
	"github.com/rs/cors"
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

type OauthService struct {
	log *zap.Logger

	router *mux.Router
}

func NewOauthService(log *zap.Logger, router *mux.Router) *OauthService {
	return &OauthService{
		log:    log,
		router: router,
	}
}

func (s *OauthService) Register(configs map[string]Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, token string) {
	log := s.log.Named("Register")
	for key, val := range configs {
		registrar, ok := Registrars[key]
		if !ok {
			log.Warn("No such auth type in config", zap.String("type", key))
			continue
		}
		registrar(s.log, s.router, &val, accClient, nsClient, token)
	}
}

func (s *OauthService) Run(port string, corsAllowed []string) {
	handler := cors.New(cors.Options{
		AllowedOrigins:   corsAllowed,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS", "HEAD"},
		AllowCredentials: true,
	}).Handler(s.router)

	s.log.Debug("listen", zap.String("port", port))
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), handler)
	if err != nil {
		s.log.Fatal("Failed to start server", zap.Error(err))
	}
}
