// controllers/annual_income_management_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/common"
	"server/config"
	"server/models" // モデルのインポート
	"server/validation"

	"github.com/gin-gonic/gin"
)

type (
	IncomeDataFetcher interface {
		GetIncomeDataInRangeApi(c *gin.Context)
		GetDateRangeApi(c *gin.Context)
		GetYearIncomeAndDeductionApi(c *gin.Context)
		InsertIncomeDataApi(c *gin.Context)
		UpdateIncomeDataApi(c *gin.Context)
		DeleteIncomeDataApi(c *gin.Context)
	}

	apiGetIncomeDataFetcher struct{}

	incomeDataInRangeResponse struct {
		Result   []models.IncomeData `json:"result,omitempty"`
		ErrorMsg string              `json:"error_msg,omitempty"`
	}
)

func NewIncomeDataFetcher() IncomeDataFetcher {
	return &apiGetIncomeDataFetcher{}
}

// GetIncomeDataInRangeApi は登録された給料及び賞与の金額を指定期間で返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetIncomeDataInRangeApi(c *gin.Context) {
	var common common.CommonFetcher = common.NewCommonFetcher()
	// パラメータから日付の始まりと終わりを取得
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userIdPrams, _ := common.StrToInt(c.Query("user_id"))

	validator := validation.RequestYearIncomeAndDeductiontData{
		UserId:    userIdPrams,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		fmt.Println(errMsgList)
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	incomeData, err := dbFetcher.GetIncomeDataInRange(startDate, endDate, userIdPrams)

	if err != nil {
		response := incomeDataInRangeResponse{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := incomeDataInRangeResponse{
		Result: incomeData,
	}
	c.JSON(http.StatusOK, response)
}

// GetDateRangeApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetDateRangeApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userId := c.Query("user_id")

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	paymentDate, err := dbFetcher.GetDateRange(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"result": paymentDate})
}

// GetYearIncomeAndDeductionApi は各年ごとの収入、差引額、手取を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetYearIncomeAndDeductionApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userId := c.Query("user_id")

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	yearIncomeData, err := dbFetcher.GetYearsIncomeAndDeduction(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"result": yearIncomeData})
}

// InsertIncomeDataApi は新規登録
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) InsertIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData struct {
		Data []models.InsertIncomeData `json:"data"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 収入データベースへ新しいデータ登録
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.InsertIncome(requestData.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "新規給料情報を登録致しました。"})
}

// UpdateIncomeDataApi は更新
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) UpdateIncomeDataApi(c *gin.Context) {
	var requestData struct {
		Data []models.UpdateIncomeData `json:"data"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 収入データベースの更新
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.UpdateIncome(requestData.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "給料情報の更新が問題なく成功しました。"})
}

// DeleteIncomeDataApi は削除
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) DeleteIncomeDataApi(c *gin.Context) {

	incomeForecastId := c.Query("income_forecast_id")

	// 収入データベースの指定されたIDの削除
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.DeleteIncome([]models.DeleteIncomeData{{IncomeForecastID: incomeForecastId}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースからの削除中にエラーが発生しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "給料情報の削除が問題なく成功しました。"})
}
