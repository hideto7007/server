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
