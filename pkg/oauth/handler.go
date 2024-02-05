package oauth

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"github.com/infinimesh/proto/node"
	"io"
	"net/http"
	"time"

	"github.com/infinimesh/proto/node/accounts"
	"github.com/infinimesh/proto/node/nodeconnect"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
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

func GithubLoginHandler(log *zap.Logger, githubConfig *oauth2.Config, config *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Get login request")
		url := githubConfig.AuthCodeURL(config.State)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func GithubCheckoutHandler(log *zap.Logger, githubConfig *oauth2.Config, config *Config, client nodeconnect.AccountsServiceClient, infinimeshToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Get checkout request")
		if r.FormValue("state") != config.State {
			log.Error("Wrong request state", zap.String("state", r.FormValue("state")))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

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

		var access *OrgAccess

		for _, team := range teams {
			var ok = false
			for key, org := range config.OrganizationMapping {
				if team.Org.Name == key {
					ok = true
					access = &org
					break
				}
			}
			if ok {
				break
			}
		}
		if access == nil {
			log.Error("Wrong access check")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		tokenReq := connect.NewRequest(&node.TokenRequest{
			Auth: &accounts.Credentials{
				Type: "oauth2-github",
				Data: []string{token.AccessToken},
			},
			Exp: time.Now().Unix() + 24*60*60*30,
		})
		tokenResponse, err := client.Token(context.Background(), tokenReq)
		if err != nil {
			createAccRequest := connect.NewRequest(&accounts.CreateRequest{
				Account: &accounts.Account{
					Title: user.Name,
				},
				Credentials: &accounts.Credentials{
					Type: "oauth2-github",
					Data: []string{token.AccessToken},
				},
				Namespace: access.Namespace,
				Access:    &access.Level,
			})
			createAccRequest.Header().Set("Authorization", "bearer "+infinimeshToken)
			_, err := client.Create(context.Background(), createAccRequest)
			if err != nil {
				log.Error("Failed to create account", zap.Error(err))
				return
			}
			tokenResponse, err = client.Token(context.Background(), tokenReq)
			if err != nil {
				log.Error("Failed to get token of created account", zap.Error(err))
				return
			}
		}
		w.Write([]byte(tokenResponse.Msg.GetToken()))
	}
}
