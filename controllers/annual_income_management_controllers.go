// controllers/annual_income_management_controllers.go
package controllers

import (
	"net/http"
	"server/models" // モデルのインポート

	"github.com/gin-gonic/gin"
)

// GetIncomeDataInRangeApi は登録された給料及び賞与の金額を指定期間で返すAPIです。
//
// 引数:
//   - c: Ginコンテキスト
//

func GetIncomeDataInRangeApi(c *gin.Context) {
	// パラメータから日付の始まりと終わりを取得
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// データベースから指定範囲のデータを取得
	incomeData, err := models.GetIncomeDataInRange(startDate, endDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"data": incomeData})
}
