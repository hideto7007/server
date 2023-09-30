// controllers/price_management_controllers.go
package controllers

import (
	"net/http"
	"server/common"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

type PriceInfo struct {
	LeftAmount  int `json:"left_amount"`
	TotalAmount int `json:"total_amount"`
}

type Response struct {
	PriceInfo PriceInfo `json:"result"`
	Error     string    `json:"error,omitempty"`
}

func priceCalc(moneyReceived, bouns, fixedCost, loan, private int) PriceInfo {

	var priceinfo PriceInfo
	priceinfo.LeftAmount = moneyReceived - fixedCost - loan - private
	priceinfo.TotalAmount = (priceinfo.LeftAmount * 12) + bouns

	return priceinfo
}

func GetPriceInfo(c *gin.Context) {
	// CORSヘッダーを設定
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers")

	data, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private")

	if err == nil {
		res := priceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])

		response := Response{PriceInfo: res}

		c.JSON(http.StatusOK, gin.H{"message": response})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": err})
	}

}

// func GetPriceInfo(Response_Writer http.ResponseWriter, req *http.Request) {

// 	// レスポンスのContent-Typeを設定
// 	Response_Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

// 	// HTTPリクエストからクエリパラメータを取得してパラメーラ値を整数値に変換
// 	data, err := common.IntgetPrameter(req, "money_received", "bouns", "fixed_cost", "loan", "private")

// 	var response Response
// 	if err == nil {
// 		res := priceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])
// 		response.PriceInfo = res
// 	} else {
// 		// エラーメッセージをErrorフィールドに設定
// 		response.Error = err.Error()
// 	}

// 	// Response型のデータをJSONに変換
// 	jsonResponse, _ := json.Marshal(response)

// 	// JSONレスポンスを書き込む
// 	Response_Writer.WriteHeader(http.StatusOK)
// 	Response_Writer.Write(jsonResponse)
// }
