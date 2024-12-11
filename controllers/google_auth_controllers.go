// controllers/google_auth_controllers.go
package controllers

import (
	"encoding/json"
	"net/http"
	"server/config"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type (
	GoogleService interface {
		HandleGoogleSignIn(c *gin.Context)
		HandleGoogleCallback(c *gin.Context)
	}

	GoogleUserInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	GoogleManager struct {
		GoogleService config.GoogleService
	}
)

func NewGoogleService(GoogleService config.GoogleService) GoogleService {
	return &GoogleManager{
		GoogleService: GoogleService,
	}
}

func (gm *GoogleManager) HandleGoogleSignIn(c *gin.Context) {
	url := gm.GoogleService.GoogleAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) HandleGoogleCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	code := c.Query("code")
	if code == "" {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: "Code not found",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// 認証コードからトークンを取得
	token, err := gm.GoogleService.GoogleOauth().Exchange(c, code)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// トークンを使ってユーザー情報を取得
	client := gm.GoogleService.GoogleOauth().Client(c, token)
	resp, err := client.Get(config.OauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer resp.Body.Close()

	// ユーザー情報の取得
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		response := utils.ResponseWithSlice[GoogleUserInfo]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
