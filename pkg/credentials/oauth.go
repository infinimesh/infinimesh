package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/infinimesh/pkg/oauth"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type OauthVerifier func(string) (string, error)

var verifiers = map[string]OauthVerifier{
	"github": GithubVerifier,
}

func GithubVerifier(token string) (string, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var user oauth.GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", err
	}

	return user.Login, nil
}

type OauthCredentials struct {
	OauthType string `json:"oauth_type"`
	Value     string `json:"value"`

	log *zap.Logger
	driver.DocumentMeta
}

func NewOauthCredentials(oauthType, token string) (*OauthCredentials, error) {
	split := strings.Split(oauthType, "-")
	if len(split) != 2 {
		return nil, errors.New("wrong oauth type")
	}

	verifier, ok := verifiers[split[1]]
	if !ok {
		return nil, errors.New("wrong oauth type")
	}
	value, err := verifier(token)
	if err != nil {
		return nil, err
	}

	return &OauthCredentials{
		OauthType: oauthType,
		Value:     value,
	}, nil
}

func (o *OauthCredentials) Listable() []string {
	return []string{o.OauthType, o.Value}
}

func (o *OauthCredentials) Authorize(args ...string) bool {
	split := strings.Split(o.OauthType, "-")
	if len(split) != 2 {
		return false
	}

	verifier, ok := verifiers[split[1]]
	if !ok {
		return false
	}

	s, err := verifier(args[0])
	if err != nil {
		return false
	}
	return o.Value == s
}

func (o *OauthCredentials) Type() string {
	return o.OauthType
}

func (o *OauthCredentials) Key() string {
	return o.ID.Key()
}

func (o *OauthCredentials) Find(ctx context.Context, db driver.Database) bool {
	query := `FOR cred IN @@credentials FILTER cred.oauth_type == @oauth_type AND cred.value == @value RETURN cred`
	c, err := db.Query(ctx, query, map[string]interface{}{
		"oauth_type":   o.OauthType,
		"value":        o.Value,
		"@credentials": schema.CREDENTIALS_COL,
	})
	if err != nil {
		return false
	}
	defer c.Close()

	_, err = c.ReadDocument(ctx, o)
	return err == nil
}

func (o *OauthCredentials) FindByKey(ctx context.Context, collection driver.Collection, key string) error {
	_, err := collection.ReadDocument(ctx, key, o)
	return err
}

func (o *OauthCredentials) SetLogger(logger *zap.Logger) {
	o.log = logger.Named("Oauth Auth")
}
