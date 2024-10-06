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
		GetSingInApi(c *gin.Context)
		GetRefreshTokenApi(c *gin.Context)
		GetSingUpApi(c *gin.Context)
		// GetDateRangeApi(c *gin.Context)
		// GetYearIncomeAndDeductionApi(c *gin.Context)
		// InsertIncomeDataApi(c *gin.Context)
		// UpdateIncomeDataApi(c *gin.Context)
		// DeleteIncomeDataApi(c *gin.Context)
	}

	// JSONデータを受け取るための構造体を定義
	requestSingInData struct {
		Data []models.RequestSingInData `json:"data"`
	}

	singInResult struct {
		UserId       int    `json:"user_id"`
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	singInResponse struct {
		Token    string         `json:"token,omitempty"`
		Result   []singInResult `json:"result,omitempty"`
		ErrorMsg string         `json:"error_msg,omitempty"`
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

// getSingInApi はサインイン情報を返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) GetSingInApi(c *gin.Context) {
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
		response := singInResponse{
			ErrorMsg: "サインインに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	token, err := utils.NewToken(requestData.Data[0].UserId, 5)
	if err != nil {
		response := singInResponse{
			ErrorMsg: "ークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := singInResponse{
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

	token, err := utils.RefreshToken(userIdPrams, 1)
	if err != nil {
		response := singInResponse{
			ErrorMsg: "ークンの生成に失敗しました。",
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

// GetSingUpApi はサインイン情報を新規登録API
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingDataFetcher) GetSingUpApi(c *gin.Context) {
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
		response := singInResponse{
			ErrorMsg: "サインインに失敗しました。",
		}
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	token, err := utils.NewToken(requestData.Data[0].UserId, 5)
	if err != nil {
		response := singInResponse{
			ErrorMsg: "ークンの生成に失敗しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// サインイン成功のレスポンス
	response := singInResponse{
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
