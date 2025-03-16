// controllers/sing_in_controllers.go
package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"server/common"
	"server/config"
	"server/models" // モデルのインポート
	"server/templates"
	"server/utils"
	"server/validation"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	// email_utils "github.com/hack-31/point-app-backend/utils/email"
)

type (
	SignDataFetcher interface {
		PostSignInApi(c *gin.Context)
		GetRefreshTokenApi(c *gin.Context)
		TemporaryPostSignUpApi(c *gin.Context)
		RetryAuthEmail(c *gin.Context)
		PostSignUpApi(c *gin.Context)
		PutSignInEditApi(c *gin.Context)
		DeleteSignInApi(c *gin.Context)
		SignOutApi(c *gin.Context)
		RegisterEmailCheckNotice(c *gin.Context)
		NewPasswordUpdate(c *gin.Context)
	}

	// JSONデータを受け取るための構造体を定義

	RequestRedisKeyData struct {
		RedisKey      string `json:"redis_key"`
		AuthEmailCode string `json:"auth_email_code"`
	}

	SignInResult struct {
		UserId       int    `json:"user_id"`
		UserEmail    string `json:"user_email"`
		UserPassword string `json:"user_password"`
	}

	SignUpResult struct {
		UserId       int    `json:"user_id"`
		UserEmail    string `json:"user_email"`
	}

	TemporayPostSignUpResult struct {
		RedisKey  string `json:"redis_key"`
		UserEmail string `json:"user_email"`
		UserName  string `json:"user_name"`
	}

	RetryAuthEmailResult struct {
		RedisKey  string `json:"redis_key"`
		UserEmail string `json:"user_email"`
		UserName  string `json:"user_name"`
	}

	RequestRefreshToken struct {
		UserId int `json:"user_id"`
	}

	apiSignDataFetcher struct {
		UtilsFetcher         utils.UtilsFetcher
		CommonFetcher        common.CommonFetcher
		EmailTemplateService templates.EmailTemplateService
		RedisService         config.RedisService
	}
)

func NewSignDataFetcher(
	tokenFetcher utils.UtilsFetcher,
	CommonFetcher common.CommonFetcher,
	EmailTemplateService templates.EmailTemplateService,
	RedisService config.RedisService,
) SignDataFetcher {
	return &apiSignDataFetcher{
		UtilsFetcher:         tokenFetcher,
		CommonFetcher:        CommonFetcher,
		EmailTemplateService: EmailTemplateService,
		RedisService:         RedisService,
	}
}

// PostSignInApi はサインイン情報を返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//   - tokenFetcher utils.UtilsFetcher: tokenフィーチャー構造体
//

func (af *apiSignDataFetcher) PostSignInApi(c *gin.Context) {
	var requestData models.RequestSignInData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.RequestSignInData{
		UserEmail:    requestData.UserEmail,
		UserPassword: requestData.UserPassword,
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
	result, err := dbFetcher.GetSignIn(requestData)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// UtilsFetcher を使用してトークンを生成
	newToken, err := af.UtilsFetcher.NewToken(result.UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := af.UtilsFetcher.RefreshToken(result.UserId, utils.RefreshAuthTokenHour)
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

	subject, body, err := af.EmailTemplateService.PostSignInTemplate(
		result.UserEmail,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(result.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := utils.ResponseData[SignInResult]{
		// Token: token,
		Result: SignInResult{
			UserId:       result.UserId,
			UserEmail:    result.UserEmail,
			UserPassword: result.UserPassword,
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetRefreshTokenApi はリフレッシュトークンを返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) GetRefreshTokenApi(c *gin.Context) {
	var common common.CommonFetcher = common.NewCommonFetcher()

	// パラメータからユーザー情報取得
	userIdCheck := c.Query("user_id")
	validator := validation.RequestRefreshTokenData{
		UserId: userIdCheck,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	signInUserId, err := c.Cookie(utils.UserId)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新しいアクセストークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else if signInUserId != userIdCheck {
		response := utils.ErrorMessageResponse{
			Result: "サインインユーザーが異なっています。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	refreshToken, err := c.Cookie(utils.RefreshAuthToken)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "リフレッシュトークンがありません。再ログインしてください。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// リフレッシュトークンの検証
	token, err := af.UtilsFetcher.ParseWithClaims(refreshToken)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "リフレッシュトークンが無効です。再ログインしてください。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// クレームからユーザー情報を取得
	_, ok := af.UtilsFetcher.MapClaims(token.(*jwt.Token))
	if !ok {
		response := utils.ErrorMessageResponse{
			Result: "無効なリフレッシュトークン。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userId, _ := common.StrToInt(userIdCheck)

	newToken, err := af.UtilsFetcher.NewToken(userId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新しいアクセストークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 新しいアクセストークンをクッキーとしてセット（またはJSONとして返す）
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	// // リフレッシュトークンも更新しておく
	// c.SetCookie(utils.RefreshAuthToken, newToken, 2*60*60, "/", domain, secure, true)

	// log.Println("INFO: ", newToken)

	// リフレッシュトークン成功のレスポンス
	response := utils.ResponseData[string]{
		Result: "新しいアクセストークンが発行されました。",
	}
	c.JSON(http.StatusOK, response)
}

// TemporaryPostSignUpApi はサインイン情報を仮登録API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) TemporaryPostSignUpApi(c *gin.Context) {
	var requestData models.RequestSignUpData
	var err error
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.TemporayRequestSignUpData{
		UserEmail:    requestData.UserEmail,
		UserPassword: requestData.UserPassword,
		UserName:     requestData.UserName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// パスワードハッシュ化
	hashPassword, _ := af.UtilsFetcher.EncryptPassword(requestData.UserPassword)
	uid := uuid.New().String()
	confirmCode, _ := rand.Int(rand.Reader, big.NewInt(10000))
	// redisに登録する際のkey
	confirmCodeStr := fmt.Sprintf("%04d", confirmCode.Int64())
	key := fmt.Sprintf("%s:%s", confirmCodeStr, uid)
	// redisに登録する際のvalue
	userInfo := [...]string{
		requestData.UserEmail,
		hashPassword,
		requestData.UserName,
	}
	value := strings.Join(userInfo[:], ",") // 配列をカンマ区切りの文字列に変換

	// 保存
	if err = af.RedisService.RedisSet(key, value, time.Hour); err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	subject, body, err := af.EmailTemplateService.TemporayPostSignUpTemplate(requestData.UserName, confirmCodeStr)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(仮登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.UserEmail, subject, body, false); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール仮登録送信エラー(仮登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインアップ仮登録成功のレスポンス
	response := utils.ResponseData[TemporayPostSignUpResult]{
		Result: TemporayPostSignUpResult{
			RedisKey:  key,
			UserEmail: requestData.UserEmail,
			UserName:  requestData.UserName,
		},
	}
	c.JSON(http.StatusOK, response)
}

// RetryAuthEmail はAPIはメール認証を再通知するために使用
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) RetryAuthEmail(c *gin.Context) {
	UserEmail := c.Query("user_email")
	UserName := c.Query("user_name")
	RedisKey := c.Query("redis_key")

	var err error

	validator := validation.RequestRetryAuthEmail{
		UserEmail: UserEmail,
		UserName:  UserName,
		RedisKey:  RedisKey,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// サインアップ仮登録した情報を取得
	redisGet, err := af.RedisService.RedisGet(RedisKey)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	uid := uuid.New().String()
	confirmCode, _ := rand.Int(rand.Reader, big.NewInt(10000))

	// redisに再登録する際のキー
	confirmCodeStr := fmt.Sprintf("%04d", confirmCode.Int64())
	newKey := fmt.Sprintf("%s:%s", confirmCodeStr, uid)

	// 更新して保存
	if err = af.RedisService.RedisSet(newKey, redisGet, time.Hour); err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 前の情報は削除する
	if err = af.RedisService.RedisDel(RedisKey); err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	subject, body, err := af.EmailTemplateService.TemporayPostSignUpTemplate(UserName, confirmCodeStr)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(UserEmail, subject, body, false); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール再通知成功のレスポンス
	response := utils.ResponseData[RetryAuthEmailResult]{
		Result: RetryAuthEmailResult{
			RedisKey:  newKey,
			UserEmail: UserEmail,
			UserName:  UserName,
		},
	}
	c.JSON(http.StatusOK, response)
}

// PostSignUpApi はサインイン情報を新規登録API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) PostSignUpApi(c *gin.Context) {
	var requestData RequestRedisKeyData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 認証コード取得
	auth := strings.Split(requestData.RedisKey, ":")
	if auth[0] != requestData.AuthEmailCode {
		response := utils.ErrorMessageResponse{
			Result: "メール認証コードが間違っています。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインアップ仮登録した情報を取得
	redisGet, err := af.RedisService.RedisGet(requestData.RedisKey)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// username,password,usernameの順で文字列が連結されている
	info := strings.Split(redisGet, ",")

	userEmail := info[0]
	userPassword := info[1]
	userName := info[2]

	validator := validation.RequestSignUpData{
		UserEmail:    userEmail,
		UserPassword: userPassword,
		UserName:     userName,
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
	// requesTemporaySignUpDataの構造体を流用してデータ構造作成
	data := models.RequestSignUpData{
		UserEmail:    userEmail,
		UserPassword: userPassword,
		UserName:     userName,
	}
	userId, err := dbFetcher.PostSignUp(data)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	// 情報は削除する
	if err = af.RedisService.RedisDel(requestData.RedisKey); err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// UtilsFetcher を使用してトークンを生成
	newToken, err := af.UtilsFetcher.NewToken(userId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := af.UtilsFetcher.RefreshToken(userId, utils.RefreshAuthTokenHour)
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

	subject, body, err := af.EmailTemplateService.PostSignUpTemplate(
		data.UserName,
		data.UserEmail,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(data.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.ResponseData[SignUpResult]{
		Result: SignUpResult{
			UserId:       userId,
			UserEmail:    userEmail,
		},
	}
	c.JSON(http.StatusOK, response)
}

// PutSignInEditApi はサインイン情報を編集API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) PutSignInEditApi(c *gin.Context) {
	var requestData models.RequestSignInEditData
	var updateValue string
	var result string
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(c.Param("user_id"))

	validator := validation.RequestSignInEditData{
		UserId:       userIdCheck,
		UserEmail:    requestData.UserEmail,
		UserPassword: requestData.UserPassword,
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
	result, err := dbFetcher.PutCheck(requestData)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "更新チェックエラー",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	UserId, _ := af.CommonFetcher.StrToInt(userIdCheck)
	if err := dbFetcher.PutSignInEdit(UserId, requestData); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "サインイン情報編集に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	if result == "メールアドレス更新" {
		updateValue = requestData.UserEmail
	} else {
		updateValue = requestData.UserPassword
	}

	subject, body, err := af.EmailTemplateService.PostSignInEditTemplate(
		result,
		updateValue,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(更新): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(更新): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン編集の成功レスポンス
	response := utils.ResponseData[string]{
		Result: fmt.Sprintf("%s:%s", "サインイン編集に成功", updateValue),
	}
	c.JSON(http.StatusOK, response)
}

// DeleteSignInApi はサインイン情報を削除API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) DeleteSignInApi(c *gin.Context) {
	var requestData models.RequestSignInDeleteData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(c.Param("user_id"))

	validator := validation.RequestSignInDeleteData{
		DeleteName: requestData.DeleteName,
		UserId:     userIdCheck,
		UserEmail:  requestData.UserEmail,
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

	UserId, _ := af.CommonFetcher.StrToInt(userIdCheck)
	err := dbFetcher.DeleteSignIn(UserId, requestData)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := af.EmailTemplateService.DeleteSignInTemplate(
		requestData.DeleteName,
		requestData.UserEmail,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.ResponseData[string]{
		Result: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, response)
}

// SignOutApi はサインアウトAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) SignOutApi(c *gin.Context) {

	// パラメータからユーザー情報取得
	userEmail := c.Query("user_email")
	validator := validation.RequestSignOutData{
		UserEmail: userEmail,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", config.GlobalEnv.Domain, config.GlobalEnv.Secure, config.GlobalEnv.HttpOnly)

	subject, body, err := af.EmailTemplateService.SignOutTemplate(
		userEmail,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(サインアウト): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(userEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(サインアウト): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.ResponseData[string]{
		Result: "サインアウトに成功",
	}
	c.JSON(http.StatusOK, response)
}

// RegisterEmailCheckNotice はパスワード再発行時にすでに登録済みかを確認するために使用する
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) RegisterEmailCheckNotice(c *gin.Context) {
	UserEmail := c.Query("user_email")

	validator := validation.EmailCheckRequestData{
		UserEmail: UserEmail,
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
	// user_id取得
	userId, err := dbFetcher.GetUserId(UserEmail)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	tokenId := uuid.New().String() + common.AnyToStr(userId)

	var link string = fmt.Sprintf("%ssign_password_reset?token_id=%s", utils.GetBaseURL(), tokenId)

	subject, body, err := af.EmailTemplateService.RegisterEmailCheckNoticeTemplate(
		link,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(パスワード再発行メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(UserEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(パスワード再発行メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール再通知成功のレスポンス
	response := utils.ResponseData[string]{
		Result: "パスワード再設定通知成功",
	}
	c.JSON(http.StatusOK, response)
}

// NewPasswordUpdate はパスワード再発行時の更新API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) NewPasswordUpdate(c *gin.Context) {
	var requestData models.RequestNewPasswordUpdateData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.RequestNewPasswordUpdateData{
		TokenId:         requestData.TokenId,
		NewUserPassword: requestData.NewUserPassword,
		ConfirmPassword: requestData.ConfirmPassword,
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
	userEmail, err := dbFetcher.NewPasswordUpdate(requestData)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	subject, body, err := af.EmailTemplateService.NewPasswordUpdateTemplate(
		requestData.NewUserPassword,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(パスワード再発行メール): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(userEmail, subject, body, true); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "メール送信エラー(パスワード再発行メール): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール再通知成功のレスポンス
	response := utils.ResponseData[string]{
		Result: "パスワード再発行成功",
	}
	c.JSON(http.StatusOK, response)
}
