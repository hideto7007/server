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
	if err := c.ShouldBindJSON(&requestData); err != nil {
		response := utils.Response[requestSingInData]{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userIdCheck := common.AnyToStr(requestData.Data[0].UserId)

	validator := validation.RequestSingInData{
		UserId:       userIdCheck,
		UserName:     requestData.Data[0].UserName,
		UserPassword: requestData.Data[0].UserPassword,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.Response[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	result, err := dbFetcher.GetSingIn(requestData.Data[0])
	if err != nil || len(result) == 0 {
		response := utils.Response[requestSingInData]{
			ErrorMsg: "サインインに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userId, _ := af.CommonFetcher.StrToInt(userIdCheck)

	// UtilsFetcher を使用してトークンを生成
	token, err := af.UtilsFetcher.NewToken(userId, 12)
	if err != nil {
		response := utils.Response[requestSingInData]{
			ErrorMsg: "トークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := utils.Response[SingInResult]{
		Token: token,
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
	// パラメータからユーザー情報取得
	userIdCheck := c.Query("user_id")
	validator := validation.RequestRefreshTokenData{
		UserId: userIdCheck,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.Response[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId, _ := common.StrToInt(userIdCheck)

	token, err := af.UtilsFetcher.RefreshToken(userId, 3)
	if err != nil {
		response := utils.Response[RequestRefreshToken]{
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
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.Response[singUpResult]{
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
		response := utils.Response[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	if err := dbFetcher.PostSingUp(requestData.Data[0]); err != nil {
		response := utils.Response[string]{
			ErrorMsg: "サインアップに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインアップ成功のレスポンス
	response := utils.Response[string]{
		ResultMsg: "サインアップに成功",
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
		response := utils.Response[singInEditResult]{
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
		response := utils.Response[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	if err := dbFetcher.PutSingInEdit(requestData.Data[0]); err != nil {
		response := utils.Response[string]{
			ErrorMsg: "サインイン情報編集に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン編集の成功レスポンス
	response := utils.Response[string]{
		ResultMsg: "サインイン編集に成功",
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
		response := utils.Response[singInDeleteResult]{
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
		response := utils.Response[utils.ErrorMessages]{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dbFetcher, _, _ := models.NewSingDataFetcher(config.DataSourceName)
	err := dbFetcher.DeleteSingIn(requestData.Data[0])
	if err != nil {
		response := utils.Response[string]{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// サインイン削除の成功レスポンス
	response := utils.Response[string]{
		ResultMsg: "サインイン削除に成功",
	}
	c.JSON(http.StatusOK, response)
}
