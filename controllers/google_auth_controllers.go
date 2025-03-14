// controllers/google_auth_controllers.go
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
	GoogleService interface {
		GoogleSignInCallback(c *gin.Context)
		GoogleSignUpCallback(c *gin.Context)
		GoogleDeleteCallback(c *gin.Context)
	}

	GooglePrams struct {
		UserEmail string
		UserName  string
	}

	GoogleManager struct {
		GoogleConfig         config.GoogleConfig
		EmailTemplateService templates.EmailTemplateService
		UtilsFetcher         utils.UtilsFetcher
	}
)

func NewGoogleService(
	GoogleConfig config.GoogleConfig,
	EmailTemplateService templates.EmailTemplateService,
	utilsFetcher utils.UtilsFetcher,
) GoogleService {
	return &GoogleManager{
		GoogleConfig:         GoogleConfig,
		EmailTemplateService: EmailTemplateService,
		UtilsFetcher:         utilsFetcher,
	}
}

func (gm *GoogleManager) GoogleSignInCallback(c *gin.Context) {
	var err error
	params := GooglePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}

	validator := validation.RequestGoogleCallbackData{
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
	newToken, err := gm.UtilsFetcher.NewToken(result.UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := gm.UtilsFetcher.RefreshToken(result.UserId, utils.RefreshAuthTokenHour)
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

	subject, body, err := gm.EmailTemplateService.PostSignInTemplate(
		result.UserEmail,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(result.UserEmail, subject, body, true); err != nil {
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

func (gm *GoogleManager) GoogleSignUpCallback(c *gin.Context) {
	var err error
	params := GooglePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}

	validator := validation.RequestGoogleCallbackData{
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
		UserPassword: "google",
		UserName:     params.UserName,
	}
	if err := dbFetcher.PostSignUp(registerData); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	subject, body, err := gm.EmailTemplateService.PostSignUpTemplate(
		params.UserName,
		params.UserEmail,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(params.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseData[string]{
		Result: "google外部認証の登録が成功しました。",
	}
	c.JSON(http.StatusOK, response)
}

func (gm *GoogleManager) GoogleDeleteCallback(c *gin.Context) {
	var err error
	params := GooglePrams{
		UserEmail: c.Query("user_email"),
		UserName:  c.Query("user_name"),
	}

	validator := validation.RequestGoogleCallbackData{
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

	// client := http.DefaultClient

	// // Googleトークンを無効化
	// resp, err := gm.ControllersCommonService.GetRevoke(
	// 	client,
	// 	config.OauthGoogleRevokeURLAPI,
	// 	userInfo.Token.AccessToken,
	// )
	// if err != nil || resp.StatusCode != http.StatusOK {
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

	subject, body, err := gm.EmailTemplateService.DeleteSignInTemplate(
		params.UserName,
		params.UserEmail,
		gm.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := gm.UtilsFetcher.SendMail(params.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.ResponseData[string]{
		Result: "google外部認証の削除が成功しました。",
	}
	c.JSON(http.StatusOK, response)
}
