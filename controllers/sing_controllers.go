// controllers/sing_in_controllers.go
package controllers

import (
	"net/http"
	"server/common"
	"server/config"
	"server/models" // モデルのインポート
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
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

	singInResult struct {
		UserId       int    `json:"user_id"`
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestRefreshToken struct {
		UserId int `json:"user_id"`
	}

	refreshTokenDataResponse struct {
		Token    string `json:"token,omitempty"`
		ErrorMsg string `json:"error_msg,omitempty"`
	}

	apiSingDataFetcher struct{}
)

func NewSingDataFetcher() SingDataFetcher {
	return &apiSingDataFetcher{}
}

// PostSingInApi はサインイン情報を返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) PostSingInApi(c *gin.Context) {
	var requestData requestSingInData
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	validator := validation.RequestSingInData{
		UserId:       requestData.Data[0].UserId,
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	result, err := dbFetcher.GetSingIn(requestData.Data[0])
	if err != nil || len(result) == 0 {
		response := utils.Response{
			ErrorMsg: "サインインに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	token, err := utils.NewToken(requestData.Data[0].UserId, 12)
	if err != nil {
		response := utils.Response{
			ErrorMsg: "トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := utils.Response{
		Token: token,
		Result: []singInResult{
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
	// パラメータからユーザー情報取得
	userIdPrams, _ := common.StrToInt(c.Query("user_id"))
	validator := validation.RequestRefreshTokenData{
		UserId: userIdPrams,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	token, err := utils.RefreshToken(userIdPrams, 3)
	if err != nil {
		response := utils.Response{
			ErrorMsg: "トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// リフレッシュトークン成功のレスポンス
	response := refreshTokenDataResponse{
		Token: token,
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
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	validator := validation.RequestSingUpData{
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
		NickName:     requestData.Data[0].NickName,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	if err := dbFetcher.PostSingUp(requestData.Data[0]); err != nil {
		response := utils.Response{
			ErrorMsg: "サインアップに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.Response{
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
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	validator := validation.RequestSingInEditData{
		UserId:       requestData.Data[0].UserId,
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	if err := dbFetcher.PutSingInEdit(requestData.Data[0]); err != nil {
		response := utils.Response{
			ErrorMsg: "サインイン情報編集に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン編集の成功レスポンス
	response := utils.Response{
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
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	validator := validation.RequestSingInDeleteData{
		UserId: requestData.Data[0].UserId,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	err := dbFetcher.DeleteSingIn(requestData.Data[0])
	if err != nil {
		response := utils.Response{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.Response{
		Result: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, response)
}
