// controllers/annual_income_management_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/models" // モデルのインポート

	"github.com/gin-gonic/gin"
)

// GetIncomeDataInRangeApi は登録された給料及び賞与の金額を指定期間で返すAPI
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
	c.JSON(http.StatusOK, gin.H{"result": incomeData})
}

// GetStartDataAndEndDateApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func GetStartDataAndEndDateApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userId := c.Query("user_id")

	// データベースから指定範囲のデータを取得
	paymentDate, err := models.GetStartDataAndEndDate(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, result := range paymentDate {
		fmt.Println("debug", result)
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"result": paymentDate})
}
