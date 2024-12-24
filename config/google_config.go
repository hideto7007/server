// config/google_config.go
package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	GoogleConfig interface {
		GoogleAuthURL(RedirectURI string) string
		GoogleOauth(RedirectURI string) *oauth2.Config
		Exchange(c *gin.Context, googleAuth *oauth2.Config, code string) (*oauth2.Token, error)
		Client(c *gin.Context, googleAuth *oauth2.Config, token *oauth2.Token) *http.Client
		Get(client *http.Client, url string) (*http.Response, error)
	}

	GoogleConfigManager struct{}
)

func NewGoogleManager() GoogleConfig {
	return &GoogleConfigManager{}
}

var (
	RedirectURI       string
	GoogleOauthConfig *oauth2.Config
	scopesList        = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
)

const OauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
const OauthGoogleRevokeURLAPI = "https://accounts.google.com/o/oauth2/revoke?token="

func (gm *GoogleConfigManager) GoogleAuthURL(RedirectURI string) string {
	// 環境変数や設定から取得する
	// これだとリダイレクトとScopesのurlパスがエンコードされないまま出力される
	// GoogleOauthConfig = &oauth2.Config{
	// 	ClientID:     GlobalEnv.GoogleClientID,
	// 	ClientSecret: GlobalEnv.GoogleClientSecret,
	// 	RedirectURL:  RedirectURI, // リダイレクトURL
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email",
	// 		"https://www.googleapis.com/auth/userinfo.profile",
	// 	},
	// 	Endpoint: google.Endpoint,
	// }

	scopes := strings.Join(scopesList, "+")
	state, _ := GenerateRandomState(32)
	AuthURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		google.Endpoint.AuthURL,
		GlobalEnv.GoogleClientID,
		RedirectURI,
		scopes,
		state,
	)

	return AuthURL
}

func (gm *GoogleConfigManager) GoogleOauth(RedirectURI string) *oauth2.Config {
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     GlobalEnv.GoogleClientID,
		ClientSecret: GlobalEnv.GoogleClientSecret,
		RedirectURL:  RedirectURI, // リダイレクトURL
		Scopes:       scopesList,
		Endpoint:     google.Endpoint,
	}
	return GoogleOauthConfig
}

func (gm *GoogleConfigManager) Exchange(c *gin.Context, googleAuth *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := googleAuth.Exchange(c, code)
	return token, err
}

func (gm *GoogleConfigManager) Client(c *gin.Context, googleAuth *oauth2.Config, token *oauth2.Token) *http.Client {
	client := googleAuth.Client(c, token)
	return client
}

func (gm *GoogleConfigManager) Get(client *http.Client, url string) (*http.Response, error) {
	resp, err := client.Get(url)
	return resp, err
}
