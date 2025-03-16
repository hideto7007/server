// controllers/line_auth_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/config"
	"server/models"
	"server/templates"
	"server/utils"
	"server/validation"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	LineService interface {
		LineSignInCallback(c *gin.Context)
		LineSignUpCallback(c *gin.Context)
		LineDeleteCallback(c *gin.Context)
	}

	LinePrams struct {
		UserEmail string
		UserName  string
	}

	LineManager struct {
		EmailTemplateService templates.EmailTemplateService
		UtilsFetcher         utils.UtilsFetcher
	}
)

func NewLineService(
	EmailTemplateService templates.EmailTemplateService,
	utilsFetcher utils.UtilsFetcher,
) LineService {
	return &LineManager{
		EmailTemplateService: EmailTemplateService,
		UtilsFetcher:         utilsFetcher,
	}
}

func (lm *LineManager) LineSignInCallback(c *gin.Context) {
	var err error
	params := LinePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}
	validator := validation.RequestLineCallbackData{
		UserEmail: params.UserEmail,
		UserName:  params.UserName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcherSingIn, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcherSingIn.GetExternalAuth(params.UserEmail)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "ユーザー情報取得に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	// UtilsFetcher を使用してトークンを生成
	newToken, err := lm.UtilsFetcher.NewToken(result.UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := lm.UtilsFetcher.RefreshToken(result.UserId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", result.UserId), 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := lm.EmailTemplateService.PostSignInTemplate(
		result.UserEmail,
		lm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := lm.UtilsFetcher.SendMail(result.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseData[SignInResult]{
		// Token: token,
		Result: SignInResult{
			UserId:       result.UserId,
			UserEmail:    result.UserEmail,
			UserPassword: "",
		},
	}
	c.JSON(http.StatusOK, response)
}

func (lm *LineManager) LineSignUpCallback(c *gin.Context) {
	var err error
	params := LinePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}
	validator := validation.RequestLineCallbackData{
		UserEmail: params.UserEmail,
		UserName:  params.UserName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	registerData := models.RequestSignUpData{
		UserEmail:    params.UserEmail,
		UserPassword: "line",
		UserName:     params.UserName,
	}
	userId, err := dbFetcher.PostSignUp(registerData)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	// UtilsFetcher を使用してトークンを生成
	newToken, err := lm.UtilsFetcher.NewToken(userId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := lm.UtilsFetcher.RefreshToken(userId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", userId), 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := lm.EmailTemplateService.PostSignUpTemplate(
		params.UserName,
		params.UserEmail,
		lm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := lm.UtilsFetcher.SendMail(params.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.ResponseData[SignUpResult]{
		Result: SignUpResult{
			UserId:    userId,
			UserEmail: params.UserEmail,
		},
	}
	c.JSON(http.StatusOK, response)
}

func (lm *LineManager) LineDeleteCallback(c *gin.Context) {
	var err error
	params := LinePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}
	validator := validation.RequestLineCallbackData{
		UserEmail: params.UserEmail,
		UserName:  params.UserName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ここはフロント側で実施

	// // Lineトークンを無効化
	// err := lm.LineConfig.RevokeLineAccessToken(
	// 	userInfo.LineToken.AccessToken,
	// )
	// if err != nil {
	// 	response := utils.ErrorMessageResponse{
	// 		Result: err.Error(),
	// 	}
	// 	utils.RedirectHandleError(c, http.StatusInternalServerError, response, "無効なトークンのため削除できません。")
	// 	return
	// }

	// 削除する登録ユーザー取得
	getDbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := getDbFetcher.GetExternalAuth(params.UserEmail)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	deleteDbFetcher, _, _ := models.NewSignDataFetcher(
		config.GetDataBaseSource(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	data := models.RequestSignInDeleteData{
		UserEmail: params.UserEmail,
	}
	err = deleteDbFetcher.DeleteSignIn(result.UserId, data)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.OauthState, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := lm.EmailTemplateService.DeleteSignInTemplate(
		params.UserName,
		params.UserEmail,
		lm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := lm.UtilsFetcher.SendMail(params.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseData[string]{
		Result: "line外部認証の削除が成功しました。",
	}
	c.JSON(http.StatusOK, response)
}
