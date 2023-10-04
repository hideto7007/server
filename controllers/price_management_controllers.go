// controllers/price_management_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/common"

	"github.com/gin-gonic/gin"
)

// 月の残額と1年間の総貯金額の構造体

type PriceInfo struct {
	LeftAmount  int `json:"left_amount"`
	TotalAmount int `json:"total_amount"`
}

// 結果の構造体

type Response struct {
	PriceInfo PriceInfo `json:"result"`
	Error     string    `json:"error,omitempty"`
}

// PriceCalc は月の収入、ボーナス、固定費、ローン、プライベートの値を使用して、
// 月と1年の貯金額を計算し、PriceInfo 構造体で結果を返します。
//
// 引数:
//   - moneyReceived: 月の収入
//   - bouns: ボーナス
//   - fixedCost: 固定費
//   - loan: ローン
//   - private: プライベート支出
//
// 戻り値:
//   - PriceInfo: 月と1年の貯金額の結果を表す構造体

func PriceCalc(moneyReceived, bouns, fixedCost, loan, private int) PriceInfo {

	var priceinfo PriceInfo
	priceinfo.LeftAmount = moneyReceived - fixedCost - loan - private
	priceinfo.TotalAmount = (priceinfo.LeftAmount * 12) + bouns

	return priceinfo
}

// GetPriceInfo は価格情報を取得するエンドポイントハンドラーです。
//
// クライアントから送信されたクエリーパラメータ money_received、bouns、fixed_cost、loan、private を
// 解析し、それらの値を使用して価格計算を行います。正常な場合、計算結果を JSON レスポンスとして
// 返し、HTTPステータスコード 200 (OK) を返します。エラーが発生した場合、エラーメッセージを JSON
// レスポンスとして返し、HTTPステータスコード 400 (Bad Request) を返します。
//
// Parameters:
//   - c: Ginコンテキスト
//
// Example:
//
//	GET /get-price-info?money_received=100&bouns=50&fixed_cost=30&loan=20&private=10
//
// Returns:
//
//	JSONレスポンス例（成功時）:
//	{
//	  "message": {
//	    "PriceInfo": {
//	      "LeftAmount": 70,
//	      "TotalAmount": 110
//	    }
//	  }
//	}
//
//	JSONレスポンス例（エラー時）:
//	{
//	  "message": "Invalid query parameters"
//	}

func GetPriceInfo(c *gin.Context) {
	// CORSヘッダーを設定
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers")

	data, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private")

	if err == nil {
		res := PriceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])

		response := Response{PriceInfo: res}

		fmt.Println("debug", response)

		c.JSON(http.StatusOK, gin.H{"message": response})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}

}
