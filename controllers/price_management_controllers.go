// controllers/price_management_controllers.go
package controllers

import (
	"net/http"
	"server/common"
	"server/validation"

	"github.com/gin-gonic/gin"
)

type (
	PriceManagementFetcher interface {
		PriceCalc(moneyReceived, bouns, fixedCost, loan, private, insurance int) PriceInfo
		GetPriceInfoApi(c *gin.Context)
	}

	PriceInfo struct {
		LeftAmount  int `json:"left_amount"`
		TotalAmount int `json:"total_amount"`
	}

	Response struct {
		Result   []PriceInfo `json:"result,omitempty"`
		ErrorMsg string      `json:"error_msg,omitempty"`
	}

	apiPriceManagementFetcher struct{}
)

func NewPriceManagementFetcher() PriceManagementFetcher {
	return &apiPriceManagementFetcher{}
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

func (af *apiPriceManagementFetcher) PriceCalc(moneyReceived, bouns, fixedCost, loan, private, insurance int) PriceInfo {

	var priceinfo PriceInfo
	priceinfo.LeftAmount = moneyReceived - fixedCost - loan - private
	priceinfo.TotalAmount = ((priceinfo.LeftAmount * 12) + bouns) - insurance

	return priceinfo
}

// GetPriceInfoApi は価格情報を取得するエンドポイントハンドラーです。
//
// クライアントから送信されたクエリーパラメータ money_received、bouns、fixed_cost、loan、private を
// 解析し、それらの値を使用して価格計算を行います。正常な場合、計算結果を JSON レスポンスとして
// 返し、HTTPステータスコード 200 (OK) を返します。エラーが発生した場合、エラーメッセージを JSON
// レスポンスとして返し、HTTPステータスコード 400 (Bad Request) を返します。
//
// 引数:
//   - c: Ginコンテキスト
//
// 期待するURL:
//
//	GET /get-price-info?money_received=100&bouns=50&fixed_cost=30&loan=20&private=10
//
// 戻り値:
//
//	JSONレスポンス例（成功時）:
//	{
//	  "message": {
//	    "PriceInfo": {
//	      "left_amount": 70,
//	      "total_amount": 110
//	    }
//	  }
//	}
//
//	JSONレスポンス例（エラー時）:
//	{
//	  "message": "Invalid query parameters"
//	}

func (af *apiPriceManagementFetcher) GetPriceInfoApi(c *gin.Context) {

	validator := validation.RequestPriceManagementData{
		MoneyReceived: c.Query("money_received"),
		Bouns:         c.Query("bouns"),
		FixedCost:     c.Query("fixed_cost"),
		Loan:          c.Query("loan"),
		Private:       c.Query("private"),
		Insurance:     c.Query("insurance"),
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	var common common.CommonFetcher = common.NewCommonFetcher()
	data, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private", "insurance")

	if err == nil {
		var price PriceManagementFetcher = NewPriceManagementFetcher()
		res := price.PriceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"], data["insurance"])

		response := Response{
			Result: []PriceInfo{
				{
					LeftAmount:  res.LeftAmount,
					TotalAmount: res.TotalAmount,
				},
			},
		}
		c.JSON(http.StatusOK, response)
	} else {
		response := Response{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
	}

}
