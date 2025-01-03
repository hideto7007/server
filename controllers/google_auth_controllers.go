// controllers/google_auth_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/config"
	"server/controllers/common"
	"server/models"
	"server/templates"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	GoogleService interface {
		GoogleSignIn(c *gin.Context)
		GoogleSignUp(c *gin.Context)
		GoogleDelete(c *gin.Context)
		GoogleSignInCallback(c *gin.Context)
		GoogleSignUpCallback(c *gin.Context)
		GoogleDeleteCallback(c *gin.Context)
	}

	requesGoogleSignUpData struct {
		Data []models.RequestSignUpData `json:"data"`
	}

	GoogleManager struct {
		GoogleConfig             config.GoogleConfig
		ControllersCommonService common.ControllersCommonService
		EmailTemplateService     templates.EmailTemplateService
		UtilsFetcher             utils.UtilsFetcher
	}
)

func NewGoogleService(
	GoogleConfig config.GoogleConfig,
	ControllersCommonService common.ControllersCommonService,
	EmailTemplateService templates.EmailTemplateService,
	utilsFetcher utils.UtilsFetcher,
) GoogleService {
	return &GoogleManager{
		GoogleConfig:             GoogleConfig,
		ControllersCommonService: ControllersCommonService,
		EmailTemplateService:     EmailTemplateService,
		UtilsFetcher:             utilsFetcher,
	}
}

func (gm *GoogleManager) GoogleSignIn(c *gin.Context) {
	url := gm.GoogleConfig.GoogleAuthURL(config.GoogleSignInEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleSignUp(c *gin.Context) {
	url := gm.GoogleConfig.GoogleAuthURL(config.GoogleSignUpEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleDelete(c *gin.Context) {
	url := gm.GoogleConfig.GoogleAuthURL(config.GoogleDeleteEnv.RedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *GoogleManager) GoogleSignInCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo common.GoogleUserInfo
	var response utils.ErrorResponse
	var err error
	params := common.GooglePrams{
		Code:        c.Query(common.CODE),
		RedirectUri: config.GoogleSignInEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.GoogleAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	dbFetcherSingIn, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcherSingIn.GetExternalAuth(userInfo.UserName)
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

	subject, body, err := gm.EmailTemplateService.PostSignInTemplate(
		result[0].UserName,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(result[0].UserName, subject, body, true); err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseOk := utils.ResponseWithSlice[common.GoogleUserInfo]{
		Result: []common.GoogleUserInfo{
			{
				UserId:   userInfo.UserId,
				UserName: userInfo.UserName,
			},
		},
	}
	c.JSON(http.StatusOK, responseOk)
}

func (gm *GoogleManager) GoogleSignUpCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo common.GoogleUserInfo
	var response utils.ErrorResponse
	params := common.GooglePrams{
		Code:        c.Query(common.CODE),
		RedirectUri: config.GoogleSignUpEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.GoogleAuthCommon(c, params)

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
				UserName:     userInfo.UserName,
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

	subject, body, err := gm.EmailTemplateService.PostSignUpTemplate(
		userInfo.Name,
		userInfo.UserName,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(userInfo.UserName, subject, body, true); err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	responseOk := utils.ResponseWithSingle[string]{
		Result: "google外部認証の登録成功しました。",
	}

	c.JSON(http.StatusOK, responseOk)
}

func (gm *GoogleManager) GoogleDeleteCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo common.GoogleUserInfo
	var response utils.ErrorResponse
	var err error
	params := common.GooglePrams{
		Code:        c.Query(common.CODE),
		RedirectUri: config.GoogleDeleteEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.GoogleAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	client := http.DefaultClient

	// Googleトークンを無効化
	resp, err := gm.ControllersCommonService.GetRevoke(
		client,
		config.OauthGoogleRevokeURLAPI,
		userInfo.Token.AccessToken,
	)
	if err != nil || resp.StatusCode != http.StatusOK {
		response := utils.ErrorResponse{
			ErrorMsg: "無効なトークンのため削除できません。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 削除する登録ユーザー取得
	getDbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := getDbFetcher.GetExternalAuth(userInfo.UserName)
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
		UserName: userInfo.UserName,
	}
	err = deleteDbFetcher.DeleteSignIn(data)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := gm.EmailTemplateService.DeleteSignInTemplate(
		userInfo.Name,
		userInfo.UserName,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(userInfo.UserName, subject, body, true); err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	Okresponse := utils.ResponseWithSingle[string]{
		Result: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, Okresponse)
}
