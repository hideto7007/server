// config/google_config.go
package config

import (
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type (
	GoogleService interface {
		GoogleAuthURL(RedirectURI string) string
		GoogleOauth(RedirectURI string) *oauth2.Config
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
const OauthGoogleRevokeURLAPI = "https://accounts.google.com/o/oauth2/revoke?token="

func (gm *GoogleManager) GoogleAuthURL(RedirectURI string) string {
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
	AuthURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		google.Endpoint.AuthURL,
		GlobalEnv.GoogleClientID,
		RedirectURI,
		scopes,
		"randomstate",
	)

	return AuthURL
}

func (gm *GoogleManager) GoogleOauth(RedirectURI string) *oauth2.Config {
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     GlobalEnv.GoogleClientID,
		ClientSecret: GlobalEnv.GoogleClientSecret,
		RedirectURL:  RedirectURI, // リダイレクトURL
		Scopes:       scopesList,
		Endpoint:     google.Endpoint,
	}
	return GoogleOauthConfig
}
