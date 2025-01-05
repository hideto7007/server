// controllers/line_auth_controllers.go
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
	LineService interface {
		LineSignIn(c *gin.Context)
		LineSignUp(c *gin.Context)
		LineDelete(c *gin.Context)
		LineSignInCallback(c *gin.Context)
		LineSignUpCallback(c *gin.Context)
		LineDeleteCallback(c *gin.Context)
	}

	requesLineSignUpData struct {
		Data []models.RequestSignUpData `json:"data"`
	}

	LineManager struct {
		LineConfig               config.LineConfig
		ControllersCommonService common.ControllersCommonService
		EmailTemplateService     templates.EmailTemplateService
		UtilsFetcher             utils.UtilsFetcher
	}
)

func NewLineService(
	LineConfig config.LineConfig,
	ControllersCommonService common.ControllersCommonService,
	EmailTemplateService templates.EmailTemplateService,
	utilsFetcher utils.UtilsFetcher,
) LineService {
	return &LineManager{
		LineConfig:               LineConfig,
		ControllersCommonService: ControllersCommonService,
		EmailTemplateService:     EmailTemplateService,
		UtilsFetcher:             utilsFetcher,
	}
}

func (gm *LineManager) LineSignIn(c *gin.Context) {
	url, state := gm.LineConfig.LineAuthURL(config.LineSignInEnv.RedirectURI)
	// stateを保存（CSRF対策用）
	c.SetCookie(utils.OauthState, state, utils.SecondsInHour, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *LineManager) LineSignUp(c *gin.Context) {
	url, state := gm.LineConfig.LineAuthURL(config.LineSignUpEnv.RedirectURI)
	// stateを保存（CSRF対策用）
	c.SetCookie(utils.OauthState, state, utils.SecondsInHour, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *LineManager) LineDelete(c *gin.Context) {
	url, state := gm.LineConfig.LineAuthURL(config.LineDeleteEnv.RedirectURI)
	// stateを保存（CSRF対策用）
	c.SetCookie(utils.OauthState, state, utils.SecondsInHour, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *LineManager) LineSignInCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo *config.LineUserInfo
	var response utils.ErrorResponse
	params := common.LinePrams{
		Code:        c.Query(common.CODE),
		State:       c.Query(common.STATE),
		RedirectUri: config.LineSignInEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.LineAuthCommon(c, params)

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

	// リダイレクト
	url := gm.ControllersCommonService.RedirectSignIn(result[0].UserId, result[0].UserName, true)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *LineManager) LineSignUpCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo *config.LineUserInfo
	var response utils.ErrorResponse
	params := common.LinePrams{
		Code:        c.Query(common.CODE),
		State:       c.Query(common.STATE),
		RedirectUri: config.LineSignUpEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.LineAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	registerData := requesLineSignUpData{
		Data: []models.RequestSignUpData{
			{
				UserName:     userInfo.UserName,
				UserPassword: "line",
				NickName:     userInfo.DisplayName,
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
		userInfo.DisplayName,
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

	// リダイレクト
	url := gm.ControllersCommonService.RedirectSignIn(0, "", false)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (gm *LineManager) LineDeleteCallback(c *gin.Context) {
	// コールバックから認証コードを取得
	var httpStatus int
	var userInfo *config.LineUserInfo
	var response utils.ErrorResponse
	params := common.LinePrams{
		Code:        c.Query(common.CODE),
		State:       c.Query(common.STATE),
		RedirectUri: config.LineDeleteEnv.RedirectURI,
	}

	httpStatus, userInfo, response = gm.ControllersCommonService.LineAuthCommon(c, params)

	if httpStatus != http.StatusOK {
		c.JSON(httpStatus, response)
		return
	}

	// Lineトークンを無効化
	err := gm.LineConfig.RevokeLineAccessToken(
		userInfo.LineToken.AccessToken,
	)
	if err != nil {
		response := utils.ErrorResponse{
			ErrorMsg: err.Error(),
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
	c.SetCookie(utils.OauthState, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := gm.EmailTemplateService.DeleteSignInTemplate(
		userInfo.DisplayName,
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
	// リダイレクト
	url := gm.ControllersCommonService.RedirectSignIn(0, "", false)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
