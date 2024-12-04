// controllers/sing_in_controllers.go
package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
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
		TemporayPostSignUpApi(c *gin.Context)
		RetryAuthEmail(c *gin.Context)
		PostSignUpApi(c *gin.Context)
		PutSignInEditApi(c *gin.Context)
		DeleteSignInApi(c *gin.Context)
		SignOutApi(c *gin.Context)
	}

	// JSONデータを受け取るための構造体を定義
	requestSignInData struct {
		Data []models.RequestSignInData `json:"data"`
	}

	RequestRedisKeyData struct {
		RedisKey      string `json:"redis_key"`
		AuthEmailCode string `json:"auth_email_code"`
	}

	requestRegisterSignUpData struct {
		Data []RequestRedisKeyData `json:"data"`
	}

	requesTemporaySignUpData struct {
		Data []models.RequestSignUpData `json:"data"`
	}

	requestSignInEditData struct {
		Data []models.RequestSignInEditData `json:"data"`
	}

	requestSignInDeleteData struct {
		Data []models.RequestSignInDeleteData `json:"data"`
	}

	SignInResult struct {
		UserId       int    `json:"user_id"`
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	TemporayPostSignUpResult struct {
		RedisKey string `json:"redis_key"`
		UserName string `json:"user_name"`
		NickName string `json:"nick_name"`
	}

	RetryAuthEmailResult struct {
		RedisKey string `json:"redis_key"`
		UserName string `json:"user_name"`
		NickName string `json:"nick_name"`
	}

	RequestRefreshToken struct {
		UserId int `json:"user_id"`
	}

	signUpResult struct{}

	signInEditResult struct{}

	signInDeleteResult struct{}

	refreshTokenDataResponse struct {
		Token    string `json:"token,omitempty"`
		ErrorMsg string `json:"error_msg,omitempty"`
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
	var requestData requestSignInData
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	if err := c.ShouldBindJSON(&requestData); err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.RequestSignInData{
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcher.GetSignIn(requestData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// UtilsFetcher を使用してトークンを生成
	newToken, err := af.UtilsFetcher.NewToken(result[0].UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	refreshToken, err := af.UtilsFetcher.RefreshToken(result[0].UserId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	c.SetCookie(utils.UserId, fmt.Sprintf("%d", result[0].UserId), 0, "/", domain, secure, httpOnly)
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)

	subject, body, err := af.EmailTemplateService.PostSignInTemplate(
		result[0].UserName,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: "メールテンプレート生成エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(result[0].UserName, subject, body, true); err != nil {
		response := utils.ResponseWithSlice[requestSignInData]{
			ErrorMsg: "メール送信エラー(サインイン): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := utils.ResponseWithSlice[SignInResult]{
		// Token: token,
		Result: []SignInResult{
			{
				UserId:       result[0].UserId,
				UserName:     result[0].UserName,
				UserPassword: result[0].UserPassword,
			},
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
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	// パラメータからユーザー情報取得
	userIdCheck := c.Query("user_id")
	validator := validation.RequestRefreshTokenData{
		UserId: userIdCheck,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	signInUserId, err := c.Cookie(utils.UserId)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "新しいアクセストークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	} else if signInUserId != userIdCheck {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "サインインユーザーが異なっています。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	refreshToken, err := c.Cookie(utils.RefreshAuthToken)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンがありません。再ログインしてください。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// リフレッシュトークンの検証
	token, err := af.UtilsFetcher.ParseWithClaims(refreshToken)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンが無効です。再ログインしてください。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// クレームからユーザー情報を取得
	_, ok := af.UtilsFetcher.MapClaims(token.(*jwt.Token))
	if !ok {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "無効なリフレッシュトークン。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userId, _ := common.StrToInt(userIdCheck)

	newToken, err := af.UtilsFetcher.NewToken(userId, utils.RefreshAuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "新しいアクセストークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	// 新しいアクセストークンをクッキーとしてセット（またはJSONとして返す）
	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	// // リフレッシュトークンも更新しておく
	// c.SetCookie(utils.RefreshAuthToken, newToken, 2*60*60, "/", domain, secure, true)

	// log.Println("INFO: ", newToken)

	// リフレッシュトークン成功のレスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "新しいアクセストークンが発行されました。",
	}
	c.JSON(http.StatusOK, response)
}

// TemporayPostSignUpApi はサインイン情報を仮登録API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) TemporayPostSignUpApi(c *gin.Context) {
	var requestData requesTemporaySignUpData
	var err error
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[TemporayPostSignUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.TemporayRequestSignUpData{
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
		NickName:     requestData.Data[0].NickName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// パスワードハッシュ化
	hashPassword, _ := af.UtilsFetcher.EncryptPassword(requestData.Data[0].UserPassword)
	uid := uuid.New().String()
	confirmCode, _ := rand.Int(rand.Reader, big.NewInt(10000))
	// redisに登録する際のkey
	confirmCodeStr := fmt.Sprintf("%04d", confirmCode.Int64())
	key := fmt.Sprintf("%s:%s", confirmCodeStr, uid)
	// redisに登録する際のvalue
	userInfo := [...]string{
		requestData.Data[0].UserName,
		hashPassword,
		requestData.Data[0].NickName,
	}
	value := strings.Join(userInfo[:], ",") // 配列をカンマ区切りの文字列に変換

	// 保存
	if err = af.RedisService.RedisSet(key, value, time.Hour); err != nil {
		response := utils.ResponseWithSlice[TemporayPostSignUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	subject, body, err := af.EmailTemplateService.TemporayPostSignUpTemplate(requestData.Data[0].NickName, confirmCodeStr)
	if err != nil {
		response := utils.ResponseWithSlice[TemporayPostSignUpResult]{
			ErrorMsg: "メールテンプレート生成エラー(仮登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.Data[0].UserName, subject, body, false); err != nil {
		response := utils.ResponseWithSlice[TemporayPostSignUpResult]{
			ErrorMsg: "メール仮登録送信エラー(仮登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインアップ仮登録成功のレスポンス
	response := utils.ResponseWithSingle[TemporayPostSignUpResult]{
		Result: TemporayPostSignUpResult{
			RedisKey: key,
			UserName: requestData.Data[0].UserName,
			NickName: requestData.Data[0].NickName,
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
	UserName := c.Query("user_name")
	NickName := c.Query("nick_name")
	RedisKey := c.Query("redis_key")

	var err error

	validator := validation.RequestRetryAuthEmail{
		UserName: UserName,
		NickName: NickName,
		RedisKey: RedisKey,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// サインアップ仮登録した情報を取得
	redisGet, err := af.RedisService.RedisGet(RedisKey)
	if err != nil {
		response := utils.ResponseWithSlice[RetryAuthEmailResult]{
			ErrorMsg: err.Error(),
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
		response := utils.ResponseWithSlice[RetryAuthEmailResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// 前の情報は削除する
	if err = af.RedisService.RedisDel(RedisKey); err != nil {
		response := utils.ResponseWithSlice[RetryAuthEmailResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	subject, body, err := af.EmailTemplateService.TemporayPostSignUpTemplate(NickName, confirmCodeStr)
	if err != nil {
		response := utils.ResponseWithSlice[RetryAuthEmailResult]{
			ErrorMsg: "メールテンプレート生成エラー(メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(UserName, subject, body, false); err != nil {
		response := utils.ResponseWithSlice[RetryAuthEmailResult]{
			ErrorMsg: "メール送信エラー(メール再通知): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール再通知成功のレスポンス
	response := utils.ResponseWithSingle[RetryAuthEmailResult]{
		Result: RetryAuthEmailResult{
			RedisKey: newKey,
			UserName: UserName,
			NickName: NickName,
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
	var requestData requestRegisterSignUpData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 認証コード取得
	auth := strings.Split(requestData.Data[0].RedisKey, ":")
	if auth[0] != requestData.Data[0].AuthEmailCode {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メール認証コードが間違っています。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインアップ仮登録した情報を取得
	redisGet, err := af.RedisService.RedisGet(requestData.Data[0].RedisKey)
	if err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// username,password,nicknameの順で文字列が連結されている
	info := strings.Split(redisGet, ",")

	userName := info[0]
	userPassword := info[1]
	nickName := info[2]

	validator := validation.RequestSignUpData{
		UserName:     userName,
		UserPassword: userPassword,
		NickName:     nickName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	// requesTemporaySignUpDataの構造体を流用してデータ構造作成
	data := requesTemporaySignUpData{
		Data: []models.RequestSignUpData{
			{
				UserName:     userName,
				UserPassword: userPassword,
				NickName:     nickName,
			},
		},
	}
	if err := dbFetcher.PostSignUp(data.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	// 情報は削除する
	if err = af.RedisService.RedisDel(requestData.Data[0].RedisKey); err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	subject, body, err := af.EmailTemplateService.PostSignUpTemplate(
		data.Data[0].NickName,
		data.Data[0].UserName,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メールテンプレート生成エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(data.Data[0].UserName, subject, body, true); err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メール送信エラー(登録): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "サインアップに成功",
	}
	c.JSON(http.StatusOK, response)
}

// PutSignInEditApi はサインイン情報を編集API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSignDataFetcher) PutSignInEditApi(c *gin.Context) {
	var requestData requestSignInEditData
	var updateValue string
	var result string
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[signInEditResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(requestData.Data[0].UserId)

	validator := validation.RequestSignInEditData{
		UserId:       userIdCheck,
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcher.PutCheck(requestData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "更新チェックエラー",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}
	if err := dbFetcher.PutSignInEdit(requestData.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインイン情報編集に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	if result == "ユーザー名更新" {
		updateValue = requestData.Data[0].UserName
	} else {
		updateValue = requestData.Data[0].UserPassword
	}

	subject, body, err := af.EmailTemplateService.PostSignInEditTemplate(
		result,
		updateValue,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メールテンプレート生成エラー(更新): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.Data[0].UserName, subject, body, true); err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メール送信エラー(更新): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン編集の成功レスポンス
	response := utils.ResponseWithSingle[string]{
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
	var requestData requestSignInDeleteData
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[signInDeleteResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(requestData.Data[0].UserId)

	validator := validation.RequestSignInDeleteData{
		UserId:   userIdCheck,
		UserName: requestData.Data[0].UserName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSignDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	err := dbFetcher.DeleteSignIn(requestData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", domain, secure, httpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", domain, secure, httpOnly)

	subject, body, err := af.EmailTemplateService.DeleteSignInTemplate(
		requestData.Data[0].UserName,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メールテンプレート生成エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(requestData.Data[0].UserName, subject, body, true); err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メール送信エラー(削除): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.ResponseWithSingle[string]{
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
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false

	// パラメータからユーザー情報取得
	userName := c.Query("user_name")
	validator := validation.RequestSignOutData{
		UserName: userName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ローカルの場合
	if os.Getenv("ENV") != "local" {
		domain = os.Getenv("DOMAIN")
		secure = true
		httpOnly = true
	}

	// Cookie無効化
	c.SetCookie(utils.UserId, "", 0, "/", domain, secure, httpOnly)
	c.SetCookie(utils.AuthToken, "", -1, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, "", -1, "/", domain, secure, httpOnly)

	subject, body, err := af.EmailTemplateService.SignOutTemplate(
		userName,
		af.UtilsFetcher.DateTimeStr(time.Now(), "2006年01月02日 15:04"),
	)
	if err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メールテンプレート生成エラー(サインアウト): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// メール送信ユーティリティを呼び出し
	if err := af.UtilsFetcher.SendMail(userName, subject, body, true); err != nil {
		response := utils.ResponseWithSlice[signUpResult]{
			ErrorMsg: "メール送信エラー(サインアウト): " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "サインアウトに成功",
	}
	c.JSON(http.StatusOK, response)
}
