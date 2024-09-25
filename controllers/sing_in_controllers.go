// controllers/sing_in_controllers.go
package controllers

import (
	"net/http"
	"server/config"
	"server/enum"
	"server/models" // モデルのインポート
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
)

type (
	SingInDataFetcher interface {
		GetSingInApi(c *gin.Context)
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
		Token    string
		Result   []singInResult `json:"result"`
		ErrorMsg string         `json:"error,omitempty"`
	}

	apiSingInDataFetcher struct{}
)

func NewSingInDataFetcher() SingInDataFetcher {
	return &apiSingInDataFetcher{}
}

// getExistsSingInApi はサインイン情報を返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiSingInDataFetcher) GetSingInApi(c *gin.Context) {
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

	dbFetcher, _, _ := models.NewSingInDataFetcher(config.DataSourceName)
	result, err := dbFetcher.GetSingIn(requestData.Data[0])
	if err != nil || len(result) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "サインインに失敗しました"})
		return
	}

	token, err := utils.GenerateJWT(requestData.Data[0].UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{enum.ERROR: "トークンの生成に失敗しました"})
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

// // GetDateRangeApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// // 引数:
// //   - c: Ginコンテキスト
// //

// func (af *apiSingInDataFetcher) GetDateRangeApi(c *gin.Context) {
// 	// パラメータからユーザー情報取得
// 	userId := c.Query("user_id")

// 	// データベースから指定範囲のデータを取得
// 	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
// 	paymentDate, err := dbFetcher.GetDateRange(userId)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// JSONレスポンスを返す
// 	c.JSON(http.StatusOK, gin.H{"result": paymentDate})
// }

// // GetYearIncomeAndDeductionApi は各年ごとの収入、差引額、手取を取得するAPI
// // 引数:
// //   - c: Ginコンテキスト
// //

// func (af *apiSingInDataFetcher) GetYearIncomeAndDeductionApi(c *gin.Context) {
// 	// パラメータからユーザー情報取得
// 	userId := c.Query("user_id")

// 	// データベースから指定範囲のデータを取得
// 	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
// 	yearIncomeData, err := dbFetcher.GetYearsIncomeAndDeduction(userId)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// JSONレスポンスを返す
// 	c.JSON(http.StatusOK, gin.H{"result": yearIncomeData})
// }

// // InsertIncomeDataApi は新規登録
// // 引数:
// //   - c: Ginコンテキスト
// //

// func (af *apiSingInDataFetcher) InsertIncomeDataApi(c *gin.Context) {
// 	// JSONデータを受け取るための構造体を定義
// 	var requestData struct {
// 		Data []models.InsertIncomeData `json:"data"`
// 	}

// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 収入データベースへ新しいデータ登録
// 	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
// 	if err := dbFetcher.InsertIncome(requestData.Data); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
// 		return
// 	}

// 	// JSONレスポンスを返す
// 	c.JSON(http.StatusOK, gin.H{"message": "新規給料情報を登録致しました。"})
// }

// // UpdateIncomeDataApi は更新
// // 引数:
// //   - c: Ginコンテキスト
// //

// func (af *apiSingInDataFetcher) UpdateIncomeDataApi(c *gin.Context) {
// 	var requestData struct {
// 		Data []models.UpdateIncomeData `json:"data"`
// 	}

// 	if err := c.ShouldBindJSON(&requestData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 収入データベースの更新
// 	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
// 	if err := dbFetcher.UpdateIncome(requestData.Data); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
// 		return
// 	}

// 	// JSONレスポンスを返す
// 	c.JSON(http.StatusOK, gin.H{"message": "給料情報の更新が問題なく成功しました。"})
// }

// // DeleteIncomeDataApi は削除
// // 引数:
// //   - c: Ginコンテキスト
// //

// func (af *apiSingInDataFetcher) DeleteIncomeDataApi(c *gin.Context) {

// 	incomeForecastId := c.Query("income_forecast_id")

// 	// 収入データベースの指定されたIDの削除
// 	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
// 	if err := dbFetcher.DeleteIncome([]models.DeleteIncomeData{{IncomeForecastID: incomeForecastId}}); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースからの削除中にエラーが発生しました"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "給料情報の削除が問題なく成功しました。"})
// }

// [
//     {
//         "error": "リクエストのフォーマットが正しくありません。"
//     },
//     {
//         "error": "user_id: ユーザーIDは必須です"
//     },
//     {
//         "error": "user_name: 著者名は必須項目です。."
//     }
// ]
