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
		GetRevoke(client *http.Client, url string, AccessToken string) (*http.Response, error)
	}

	ControllersCommonManager struct {
		GoogleConfig config.GoogleConfig
	}

	GoogleUserInfo struct {
		ID            string        `json:"id"`
		UserId        int           `json:"user_id"`
		Email         string        `json:"email"`
		VerifiedEmail bool          `json:"verified_email"`
		Name          string        `json:"name"`
		GivenName     string        `json:"given_name"`
		FamilyName    string        `json:"family_name"`
		Picture       string        `json:"picture"`
		Locale        string        `json:"locale"`
		Token         *oauth2.Token `json:"token"`
	}

	GooglePrams struct {
		Code        string
		RedirectUri string
	}
)

func NewControllersCommonManager(
	GoogleConfig config.GoogleConfig,
) ControllersCommonService {
	return &ControllersCommonManager{
		GoogleConfig: GoogleConfig,
	}
}

var (
	RedirectURI       string
	GoogleOauthConfig *oauth2.Config
)

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

func (gm *ControllersCommonManager) GetRevoke(client *http.Client, url string, AccessToken string) (*http.Response, error) {
	revokeURL := fmt.Sprintf("%s%s", url, AccessToken)
	resp, err := client.Get(revokeURL)
	return resp, err
}
