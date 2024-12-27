// config/line_config.go
package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"server/common"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type (
	LineConfig interface {
		LineAuthURL(RedirectURI string) (string, string)
		GetLineAccessToken(code, redirectURI string) (*LineTokenResponse, error)
		GetLineUserInfo(accessToken string) (*LineUserInfo, error)
		GetEmail(IdToken string) (string, error)
		RevokeLineAccessToken(accessToken string) error
	}

	LineConfigManager struct {
		HTTPClient *common.HTTPClient
	}

	LineTokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		IdToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
	}

	LineUserInfo struct {
		Id          string             `json:"userId,omitempty"`
		UserId      int                `json:"user_id"`
		UserName    string             `json:"email"`
		DisplayName string             `json:"displayName,omitempty"`
		LineToken   *LineTokenResponse `json:"line_token,omitempty"`
	}
)

var localHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

// var localHeaders = map[string]string{
// 	"Content-Type": "application/json",
// }

func NewLineManager(HTTPClient *common.HTTPClient) LineConfig {
	return &LineConfigManager{
		HTTPClient: HTTPClient,
	}
}

var (
	LineRedirectURI string
	LineOauthConfig *oauth2.Config
	scopes          string
)

const OauthLineURLAPI = "https://access.line.me/oauth2/v2.1/authorize"
const OauthLineRevokeURLAPI = "https://api.line.me/oauth2/v2.1/revoke"
const OauthLineAccessTokenURLAPI = "https://api.line.me/oauth2/v2.1/token"
const ProfileLineURLAPI = "https://api.line.me/v2/profile"

func (gm *LineConfigManager) LineAuthURL(LineRedirectURI string) (string, string) {
	scopes = "profile openid email"
	state, _ := GenerateRandomState(32)
	AuthURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=%s",
		OauthLineURLAPI,
		GlobalEnv.LineClientID,
		url.QueryEscape(LineRedirectURI),
		state,
		scopes,
	)

	return AuthURL, state
}

func (gm *LineConfigManager) GetLineAccessToken(code, redirectURI string) (*LineTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", GlobalEnv.LineClientID)
	data.Set("client_secret", GlobalEnv.LineClientSecret)

	// URL エンコードされた文字列を取得
	dataBytes := []byte(data.Encode())

	resp, err := gm.HTTPClient.Post(
		OauthLineAccessTokenURLAPI,
		localHeaders,
		dataBytes,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp LineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (gm *LineConfigManager) GetLineUserInfo(accessToken string) (*LineUserInfo, error) {
	localHeaders["Authorization"] = "Bearer " + accessToken

	resp, err := gm.HTTPClient.Get(
		ProfileLineURLAPI,
		localHeaders,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo LineUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (gm *LineConfigManager) GetEmail(IdToken string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(IdToken, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["email"].(string), err
	} else {
		return "", fmt.Errorf("解析不可能でした。。。")
	}
}

func (gm *LineConfigManager) RevokeLineAccessToken(accessToken string) error {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", GlobalEnv.LineClientID)
	data.Set("client_secret", GlobalEnv.LineClientSecret)
	data.Set("access_token", accessToken)

	// URL エンコードされた文字列を取得
	dataBytes := []byte(data.Encode())

	resp, err := gm.HTTPClient.Post(
		OauthLineRevokeURLAPI,
		localHeaders,
		dataBytes,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"トークン無効化に失敗しました。ステータスコード:%d エラー内容：%s",
			resp.StatusCode,
			resp.Body,
		)
	}

	return nil
}
