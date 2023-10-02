// controllers/price_management_controllers.go
package controllers

import (
	"fmt"
	"net/http"
	"server/common"

	"github.com/gin-gonic/gin"
)

type PriceInfo struct {
	LeftAmount  int `json:"left_amount"`
	TotalAmount int `json:"total_amount"`
}

type Response struct {
	PriceInfo PriceInfo `json:"result"`
	Error     string    `json:"error,omitempty"`
}

func PriceCalc(moneyReceived, bouns, fixedCost, loan, private int) PriceInfo {

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
		res := PriceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])

		response := Response{PriceInfo: res}

		fmt.Println("debug", response)

		c.JSON(http.StatusOK, gin.H{"message": response})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
	}

}
