// controllers/sing_in_controllers.go
package controllers

import (
	"net/http"
	"os"
	"server/common"
	"server/config"
	"server/models" // モデルのインポート
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type (
	SingDataFetcher interface {
		PostSingInApi(c *gin.Context)
		GetRefreshTokenApi(c *gin.Context)
		PostSingUpApi(c *gin.Context)
		PutSingInEditApi(c *gin.Context)
		DeleteSingInApi(c *gin.Context)
	}

	// JSONデータを受け取るための構造体を定義
	requestSingInData struct {
		Data []models.RequestSingInData `json:"data"`
	}

	requestSingUpData struct {
		Data []models.RequestSingUpData `json:"data"`
	}

	requestSingInEditData struct {
		Data []models.RequestSingInEditData `json:"data"`
	}

	requestSingInDeleteData struct {
		Data []models.RequestSingInDeleteData `json:"data"`
	}

	SingInResult struct {
		UserId       int    `json:"user_id"`
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestRefreshToken struct {
		UserId int `json:"user_id"`
	}

	singUpResult struct{}

	singInEditResult struct{}

	singInDeleteResult struct{}

	refreshTokenDataResponse struct {
		Token    string `json:"token,omitempty"`
		ErrorMsg string `json:"error_msg,omitempty"`
	}

	apiSingDataFetcher struct {
		UtilsFetcher  utils.UtilsFetcher
		CommonFetcher common.CommonFetcher
	}
)

func NewSingDataFetcher(
	tokenFetcher utils.UtilsFetcher,
	CommonFetcher common.CommonFetcher,
) SingDataFetcher {
	return &apiSingDataFetcher{
		UtilsFetcher:  tokenFetcher,
		CommonFetcher: CommonFetcher,
	}
}

// PostSingInApi はサインイン情報を返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//   - tokenFetcher utils.UtilsFetcher: tokenフィーチャー構造体
//

func (af *apiSingDataFetcher) PostSingInApi(c *gin.Context) {
	var requestData requestSingInData
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	if err := c.ShouldBindJSON(&requestData); err != nil {
		response := utils.ResponseWithSlice[requestSingInData]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.RequestSingInData{
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

	dbFetcher, _, _ := models.NewSingDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	result, err := dbFetcher.GetSingIn(requestData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[requestSingInData]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// UtilsFetcher を使用してトークンを生成
	newToken, err := af.UtilsFetcher.NewToken(result[0].UserId, utils.AuthTokenHour)
	if err != nil {
		response := utils.ResponseWithSlice[requestSingInData]{
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

	c.SetCookie(utils.AuthToken, newToken, utils.AuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)
	c.SetCookie(utils.RefreshAuthToken, refreshToken, utils.RefreshAuthTokenHour*utils.SecondsInHour, "/", domain, secure, httpOnly)

	// サインイン成功のレスポンス
	response := utils.ResponseWithSlice[SingInResult]{
		// Token: token,
		Result: []SingInResult{
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

func (af *apiSingDataFetcher) GetRefreshTokenApi(c *gin.Context) {
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

// PostSingUpApi はサインイン情報を新規登録API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) PostSingUpApi(c *gin.Context) {
	var requestData requestSingUpData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[singUpResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	validator := validation.RequestSingUpData{
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

	dbFetcher, _, _ := models.NewSingDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	if err := dbFetcher.PostSingUp(requestData.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "サインアップに成功",
	}
	c.JSON(http.StatusOK, response)
}

// PutSingInEditApi はサインイン情報を編集API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) PutSingInEditApi(c *gin.Context) {
	var requestData requestSingInEditData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[singInEditResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(requestData.Data[0].UserId)

	validator := validation.RequestSingInEditData{
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

	dbFetcher, _, _ := models.NewSingDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	if err := dbFetcher.PutSingInEdit(requestData.Data[0]); err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインイン情報編集に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン編集の成功レスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "サインイン編集に成功",
	}
	c.JSON(http.StatusOK, response)
}

// DeleteSingInApi はサインイン情報を削除API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) DeleteSingInApi(c *gin.Context) {
	var requestData requestSingInDeleteData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ResponseWithSlice[singInDeleteResult]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(requestData.Data[0].UserId)

	validator := validation.RequestSingInDeleteData{
		UserId: userIdCheck,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(
		config.DataSourceName,
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	err := dbFetcher.DeleteSingIn(requestData.Data[0])
	if err != nil {
		response := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.ResponseWithSingle[string]{
		Result: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, response)
}
