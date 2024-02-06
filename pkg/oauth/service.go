package oauth

import (
	"encoding/json"
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

type OAuthService interface {
	Register(map[string]Config, nodeconnect.AccountsServiceClient, nodeconnect.NamespacesServiceClient, string)
	Run(string, []string)
}

type oauthService struct {
	log *zap.Logger

	router              *mux.Router
	registeredProviders []string
}

func NewOauthService(log *zap.Logger, router *mux.Router) *oauthService {
	return &oauthService{
		log:    log,
		router: router,
	}
}

func (s *oauthService) Register(configs map[string]Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, token string) {
	log := s.log.Named("Register")
	for key, val := range configs {
		registrar, ok := Registrars[key]
		if !ok {
			log.Warn("No such auth type in config", zap.String("type", key))
			continue
		}
		registrar(s.log, s.router, &val, accClient, nsClient, token)
		s.registeredProviders = append(s.registeredProviders, key)
	}
}

func (s *oauthService) Run(port string, corsAllowed []string) {
	s.router.HandleFunc("/oauth/providers", func(w http.ResponseWriter, r *http.Request) {
		marshal, err := json.Marshal(s.registeredProviders)
		if err != nil {
			s.log.Error("Failed to marshal providers", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to get providers. %s", err.Error())))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(marshal)
	}).Methods(http.MethodGet)

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
