// config/google_config.go
package config

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	GoogleService interface {
		GoogleAuthURL() string
		GoogleOauth() *oauth2.Config
	}

	GoogleManager struct{}
)

func NewGoogleManager() GoogleService {
	return &GoogleManager{}
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

func (gm *GoogleManager) GoogleAuthURL() string {
	// リダイレクトURLを環境によって切り替える
	if os.Getenv("ENV") == "production" {
		RedirectURI = fmt.Sprintf("%s/auth/google/callback", os.Getenv("DOMAIN"))
	} else {
		RedirectURI = "http://localhost:8080/auth/google/callback"
	}
	// 環境変数や設定から取得する
	// これだとリダイレクトとScopesのurlパスがエンコードされないまま出力される
	// GoogleOauthConfig = &oauth2.Config{
	// 	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	// 	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	// 	RedirectURL:  RedirectURI, // リダイレクトURL
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email",
	// 		"https://www.googleapis.com/auth/userinfo.profile",
	// 	},
	// 	Endpoint: google.Endpoint,
	// }

	scopes := strings.Join(scopesList, "+")
	AuthURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		google.Endpoint.AuthURL,
		os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectURI,
		scopes,
		"randomstate",
	)

	return AuthURL
}

func (gm *GoogleManager) GoogleOauth() *oauth2.Config {
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  RedirectURI, // リダイレクトURL
		Scopes:       scopesList,
		Endpoint:     google.Endpoint,
	}
	return GoogleOauthConfig
}
