package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type (
	ControllersCommonService interface {
		GoogleAuthCommon(c *gin.Context, params GooglePrams) (
			int,
			GoogleUserInfo,
			utils.ErrorResponse,
		)
		LineAuthCommon(c *gin.Context, params LinePrams) (
			int,
			*config.LineUserInfo,
			utils.ErrorResponse,
		)
		GetRevoke(client *http.Client, url string, AccessToken string) (*http.Response, error)
	}

	ControllersCommonManager struct {
		GoogleConfig config.GoogleConfig
		LineConfig   config.LineConfig
	}

	GoogleUserInfo struct {
		ID            string        `json:"id,omitempty"`
		UserId        int           `json:"user_id"`
		UserName      string        `json:"user_name"`
		VerifiedEmail bool          `json:"verified_email,omitempty"`
		Name          string        `json:"name,omitempty"`
		GivenName     string        `json:"given_name,omitempty"`
		FamilyName    string        `json:"family_name,omitempty"`
		Picture       string        `json:"picture,omitempty"`
		Locale        string        `json:"locale,omitempty"`
		Token         *oauth2.Token `json:"token,omitempty"`
	}

	GooglePrams struct {
		Code        string
		RedirectUri string
	}

	LinePrams struct {
		Code        string
		State       string
		RedirectUri string
	}
)

func NewControllersCommonManager(
	GoogleConfig config.GoogleConfig,
	LineConfig config.LineConfig,

) ControllersCommonService {
	return &ControllersCommonManager{
		GoogleConfig: GoogleConfig,
		LineConfig:   LineConfig,
	}
}

var (
	RedirectURI       string
	GoogleOauthConfig *oauth2.Config
)

const CODE = "code"
const REDIRECT_URI = "redirect_uri"
const STATE = "state"

func (gm *ControllersCommonManager) GoogleAuthCommon(c *gin.Context, params GooglePrams) (
	int,
	GoogleUserInfo,
	utils.ErrorResponse,
) {
	var userInfo GoogleUserInfo
	validator := validation.RequestGoogleCallbackData{
		Code:        params.Code,
		RedirectUri: params.RedirectUri,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorResponse{
			Result: errMsgList,
		}
		return http.StatusBadRequest, userInfo, response
	}

	var googleAuth *oauth2.Config = gm.GoogleConfig.GoogleOauth(params.RedirectUri)

	// 認証コードからトークンを取得
	token, err := gm.GoogleConfig.Exchange(c, googleAuth, params.Code)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		return http.StatusInternalServerError, userInfo, response
	}

	// トークンを使ってユーザー情報を取得
	client := gm.GoogleConfig.Client(c, googleAuth, token)
	resp, err := gm.GoogleConfig.Get(client, config.OauthGoogleURLAPI+token.AccessToken)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		return http.StatusInternalServerError, userInfo, response
	}
	defer resp.Body.Close()

	userInfo.Token = token
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		return http.StatusInternalServerError, userInfo, response
	}

	return http.StatusOK, userInfo, utils.ErrorResponse{}
}

func (gm *ControllersCommonManager) LineAuthCommon(c *gin.Context, params LinePrams) (
	int,
	*config.LineUserInfo,
	utils.ErrorResponse,
) {
	var userInfo *config.LineUserInfo
	validator := validation.RequestLineCallbackData{
		Code:        params.Code,
		State:       params.State,
		RedirectUri: params.RedirectUri,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorResponse{
			Result: errMsgList,
		}
		return http.StatusBadRequest, nil, response
	}

	// stateを検証（CSRF対策）
	savedState, err := c.Cookie(utils.OauthState)
	if err != nil || params.State != savedState {
		response := utils.ErrorResponse{
			ErrorMsg: "無効なステートです。",
		}
		return http.StatusBadRequest, nil, response
	}

	// アクセストークン取得
	tokenResp, err := gm.LineConfig.GetLineAccessToken(
		params.Code,
		params.RedirectUri,
	)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "無効なLineアクセストークンです。",
		}
		return http.StatusInternalServerError, nil, response
	}

	// ユーザー情報取得
	userInfo, err = gm.LineConfig.GetLineUserInfo(tokenResp.AccessToken)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "Lineのユーザー情報取得に失敗しました。",
		}
		return http.StatusInternalServerError, nil, response
	}

	// メールアドレス取得
	email, err := gm.LineConfig.GetEmail(tokenResp.IdToken)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		return http.StatusInternalServerError, nil, response
	}

	// lineトークン情報セット
	userInfo.LineToken = tokenResp
	userInfo.UserName = email

	return http.StatusOK, userInfo, utils.ErrorResponse{}
}

func (gm *ControllersCommonManager) GetRevoke(client *http.Client, url string, AccessToken string) (*http.Response, error) {
	revokeURL := fmt.Sprintf("%s%s", url, AccessToken)
	resp, err := client.Get(revokeURL)
	return resp, err
}
