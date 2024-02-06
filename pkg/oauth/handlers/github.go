package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/infinimesh/infinimesh/pkg/oauth/config"

	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	"github.com/infinimesh/proto/node/accounts"
	"github.com/infinimesh/proto/node/nodeconnect"

	"connectrpc.com/connect"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Organization struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Team struct {
	Name string       `json:"name"`
	Org  Organization `json:"organization"`
}

func getGithubUser(token string) (*GithubUser, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var user GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getGithubUserTeams(token string) ([]*Team, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/user/teams", nil)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var teams []*Team
	err = json.Unmarshal(body, &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

type StateInfo struct {
	Token  string
	Method string
}

type GithubRegistrar struct {
	states map[string]*StateInfo
	mu     *sync.Mutex
}

func (g *GithubRegistrar) Register(logger *zap.Logger, router *mux.Router, config *config.Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, infinimeshToken string) {
	log := logger.Named("GithubOauth")
	log.Info("Init Github Oauth2 handlers")

	g.states = map[string]*StateInfo{}
	g.mu = &sync.Mutex{}

	githubConfig := &oauth2.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  config.RedirectUrl,
		Scopes:       config.Scopes,
	}

	router.HandleFunc("oauth/github/login", g.GithubLoginHandler(log, githubConfig, config))

	router.HandleFunc("oauth/github/checkout", g.GithubCheckoutHandler(log, githubConfig, config, accClient, nsClient, infinimeshToken))
}

func (g *GithubRegistrar) GithubLoginHandler(log *zap.Logger, githubConfig *oauth2.Config, config *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Get login request")
		if r.FormValue("state") == "" {
			log.Error("No state in request")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		if r.FormValue("method") == "" {
			log.Error("No method in request")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		url := githubConfig.AuthCodeURL(r.FormValue("state"))

		token := ""
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			split := strings.Split(authHeader, " ")
			if len(split) == 2 {
				token = split[1]
			}
		}

		g.mu.Lock()
		g.states[r.FormValue("state")] = &StateInfo{
			Token:  token,
			Method: r.FormValue("method"),
		}
		g.mu.Unlock()

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (g *GithubRegistrar) GithubCheckoutHandler(log *zap.Logger, githubConfig *oauth2.Config, cfg *config.Config, accClient nodeconnect.AccountsServiceClient, nsClient nodeconnect.NamespacesServiceClient, infinimeshToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Get checkout request")

		g.mu.Lock()
		defer g.mu.Unlock()
		state, ok := g.states[r.FormValue("state")]
		if !ok {
			log.Error("Wrong request state", zap.String("state", r.FormValue("state")))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		delete(g.states, r.FormValue("state"))

		token, err := githubConfig.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			log.Error("Failed to get token", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		user, err := getGithubUser(token.AccessToken)
		if err != nil {
			log.Error("Failed to get user", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		teams, err := getGithubUserTeams(token.AccessToken)
		if err != nil {
			log.Error("Failed to get user teams", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		var orgAccess *config.OrgAccess

		for _, team := range teams {
			var ok = false
			for key, org := range cfg.OrganizationMapping {
				if team.Org.Name == key {
					ok = true
					orgAccess = &org
					break
				}
			}
			if ok {
				break
			}
		}
		if orgAccess == nil {
			log.Error("Wrong access check")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		if state.Method == "sign_in" {
			tokenReq := connect.NewRequest(&node.TokenRequest{
				Auth: &accounts.Credentials{
					Type: "oauth2-github",
					Data: []string{token.AccessToken},
				},
				Exp: time.Now().Unix() + 24*60*60*30,
			})
			tokenResponse, err := accClient.Token(context.Background(), tokenReq)
			if err != nil {
				createAccRequest := connect.NewRequest(&accounts.CreateRequest{
					Account: &accounts.Account{
						Title: user.Name,
					},
					Credentials: &accounts.Credentials{
						Type: "oauth2-github",
						Data: []string{token.AccessToken},
					},
				})
				createAccRequest.Header().Set("Authorization", "bearer "+infinimeshToken)
				createAccResponse, err := accClient.Create(context.Background(), createAccRequest)
				if err != nil {
					log.Error("Failed to create account", zap.Error(err))
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
					return
				}

				joinRequest := connect.NewRequest(&node.JoinRequest{
					Namespace: orgAccess.Namespace,
					Account:   createAccResponse.Msg.GetAccount().GetUuid(),
					Access:    access.Level(orgAccess.Level),
				})
				joinRequest.Header().Set("Authorization", "bearer "+infinimeshToken)
				_, err = nsClient.Join(context.Background(), joinRequest)
				if err != nil {
					log.Error("Failed to join account", zap.Error(err))
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
					return
				}

				tokenResponse, err = accClient.Token(context.Background(), tokenReq)
				if err != nil {
					log.Error("Failed to get token of created account", zap.Error(err))
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
					return
				}
			}
			w.Write([]byte(tokenResponse.Msg.GetToken()))
		} else if state.Method == "link" {
			accToken := state.Token

			accRequest := connect.NewRequest(&accounts.Account{
				Uuid: "me",
			})
			accRequest.Header().Set("Authorization", "bearer "+accToken)

			get, err := accClient.Get(context.Background(), accRequest)
			if err != nil {
				log.Error("Failed to get user", zap.Error(err))
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			setReq := connect.NewRequest(&node.SetCredentialsRequest{
				Uuid: get.Msg.GetUuid(),
				Credentials: &accounts.Credentials{
					Type: "oauth2-github",
					Data: []string{token.AccessToken},
				},
			})
			setReq.Header().Set("Authorization", "bearer"+infinimeshToken)
			_, err = accClient.SetCredentials(context.Background(), setReq)
			if err != nil {
				log.Error("Failed to set credentials", zap.Error(err))
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}
}
