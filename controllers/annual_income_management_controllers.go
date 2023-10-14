// controllers/annual_income_management_controllers.go
package controllers

import (
	"net/http"
	"server/config"
	"server/models" // モデルのインポート

	"github.com/gin-gonic/gin"
)

type (
	IncomeDataFetcher interface {
		GetIncomeDataInRangeApi(c *gin.Context)
		GetStartDataAndEndDateApi(c *gin.Context)
		GetYearIncomeAndDeductionApi(c *gin.Context)
		InsertIncomeDataApi(c *gin.Context)
		UpdateIncomeDataApi(c *gin.Context)
		DeleteIncomeDataApi(c *gin.Context)
	}

	apiGetIncomeDataFetcher struct{}
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
	// パラメータから日付の始まりと終わりを取得
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// データベースから指定範囲のデータを取得
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
	incomeData, err := dbFetcher.GetIncomeDataInRange(startDate, endDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"result": incomeData})
}

// GetStartDataAndEndDateApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetStartDataAndEndDateApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userId := c.Query("user_id")

	// データベースから指定範囲のデータを取得
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
	paymentDate, err := dbFetcher.GetStartDataAndEndDate(userId)

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
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
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
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.InsertIncome(requestData.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "登録成功"})
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
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.UpdateIncome(requestData.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの挿入中にエラーが発生しました"})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteIncomeDataApi は削除
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) DeleteIncomeDataApi(c *gin.Context) {
	var requestData struct {
		Data []models.DeleteIncomeData `json:"data"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 収入データベースの指定されたIDの削除
	var dbFetcher models.AnuualIncomeFetcher = models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.DeleteIncome(requestData.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースからの削除中にエラーが発生しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "削除データ成功"})
}
