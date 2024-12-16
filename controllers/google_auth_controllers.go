// controllers/google_auth_controllers.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"server/models"
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type (
	GoogleService interface {
		GoogleAuthCommon(c *gin.Context, code Prams) (
			int,
			GoogleUserInfo,
			utils.ErrorResponse,
		)
		GoogleSignIn(c *gin.Context)
		GoogleSignUp(c *gin.Context)
		GoogleDelete(c *gin.Context)
		GoogleSignInCallback(c *gin.Context)
		GoogleSignUpCallback(c *gin.Context)
		GoogleDeleteCallback(c *gin.Context)
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

	Prams struct {
		Code        string
		RedirectUri string
	}

	requesGoogleSignUpData struct {
		Data []models.RequestSignUpData `json:"data"`
	}

	GoogleManager struct {
		GoogleService config.GoogleService
		UtilsFetcher  utils.UtilsFetcher
	}
)

func NewGoogleService(
	GoogleService config.GoogleService,
	utilsFetcher utils.UtilsFetcher,
) GoogleService {
	return &GoogleManager{
		GoogleService: GoogleService,
		UtilsFetcher:  utilsFetcher,
	}
}

const CODE = "code"
const REDIRECT_URI = "redirect_uri"

func (gm *GoogleManager) GoogleSignIn(c *gin.Context) {
	url := gm.GoogleService.GoogleAuthURL(config.GoogleSignInEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleSignUp(c *gin.Context) {
	url := gm.GoogleService.GoogleAuthURL(config.GoogleSignUpEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleDelete(c *gin.Context) {
	url := gm.GoogleService.GoogleAuthURL(config.GoogleDeleteEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleAuthCommon(c *gin.Context, params Prams) (
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

	var googleAuth *oauth2.Config = gm.GoogleService.GoogleOauth(params.RedirectUri)

	// 認証コードからトークンを取得
	token, err := googleAuth.Exchange(c, params.Code)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		return http.StatusInternalServerError, userInfo, response
	}

	// トークンを使ってユーザー情報を取得
	client := googleAuth.Client(c, token)
	resp, err := client.Get(config.OauthGoogleURLAPI + token.AccessToken)
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

func (gm *GoogleManager) GoogleSignInCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo GoogleUserInfo
	var response utils.ErrorResponse
	var err error
	params := Prams{
		Code:        c.Query(CODE),
		RedirectUri: config.GoogleSignInEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.GoogleAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	dbFetcherSingIn, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcherSingIn.GetExternalAuth(userInfo.Email)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	// UtilsFetcher を使用してトークンを生成
	newToken, err := gm.UtilsFetcher.NewToken(result[0].UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := gm.UtilsFetcher.RefreshToken(result[0].UserId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", result[0].UserId), 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	// DB登録ユーザーIDも取得
	userInfo.UserId = result[0].UserId
	c.JSON(http.StatusOK, userInfo)
}

func (gm *GoogleManager) GoogleSignUpCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo GoogleUserInfo
	var response utils.ErrorResponse
	params := Prams{
		Code:        c.Query(CODE),
		RedirectUri: config.GoogleSignUpEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.GoogleAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	registerData := requesGoogleSignUpData{
		Data: []models.RequestSignUpData{
			{
				UserName:     userInfo.Email,
				UserPassword: "google",
				NickName:     userInfo.Name,
			},
		},
	}
	if err := dbFetcher.PostSignUp(registerData.Data[0]); err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

func (gm *GoogleManager) GoogleDeleteCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo GoogleUserInfo
	var response utils.ErrorResponse
	var err error
	params := Prams{
		Code:        c.Query(CODE),
		RedirectUri: config.GoogleDeleteEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.GoogleAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	// Googleトークンを無効化
	revokeURL := fmt.Sprintf("%s%s", config.OauthGoogleRevokeURLAPI, userInfo.Token.AccessToken)
	resp, err := http.Get(revokeURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		response := utils.ErrorResponse{
			ErrorMsg: "無効なトークンのため削除できません",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 削除する登録ユーザー取得
	getDbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := getDbFetcher.GetExternalAuth(userInfo.Email)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	deleteDbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	data := models.RequestSignInDeleteData{
		UserId:   result[0].UserId,
		UserName: userInfo.Email,
	}
	err = deleteDbFetcher.DeleteSignIn(data)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	Okresponse := utils.ResponseWithSingle[string]{
		Result: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, Okresponse)
}
