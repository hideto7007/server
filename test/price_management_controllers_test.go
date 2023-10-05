// test/price_management_controllers_test.go
package common

import (
	"net/http"
	"net/http/httptest"

	"server/common"
	"server/controllers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPriceCalc(t *testing.T) {
	t.Run("success PriceCalc()", func(t *testing.T) {
		// テストケース1: 正常な整数の計算

		result := controllers.PriceCalc(300, 100, 50, 50, 50)

		assert.Equal(t, 150, result.LeftAmount)
		assert.Equal(t, 1900, result.TotalAmount)

		t.Logf("result.LeftAmount: %d", result.LeftAmount)
		t.Logf("result.TotalAmount: %d", result.TotalAmount)
	})
}

func TestGetPriceInfo(t *testing.T) {
	t.Run("success GetPriceInfoApi()", func(t *testing.T) {
		// テスト用のGinコンテキストを作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?money_received=300&bouns=100&fixed_cost=50&loan=50&private=50", nil)

		paramMap, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private")

		// PriceCalc 関数をモック化
		res := controllers.PriceCalc(paramMap["money_received"], paramMap["bouns"], paramMap["fixed_cost"], paramMap["loan"], paramMap["private"])

		// GetPriceInfoApi 関数を呼び出し
		controllers.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusOK, c.Writer.Status())

		// レスポンスの JSON データを取得
		response := controllers.Response{
			PriceInfo: controllers.PriceInfo{
				LeftAmount:  res.LeftAmount,
				TotalAmount: res.TotalAmount,
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, 150, response.PriceInfo.LeftAmount)
		assert.Equal(t, 1900, response.PriceInfo.TotalAmount)

		t.Logf("response.PriceInfo.LeftAmount: %d", response.PriceInfo.LeftAmount)
		t.Logf("response.PriceInfo.TotalAmount: %d", response.PriceInfo.TotalAmount)
		t.Logf("err: %v", err)
	})

	t.Run("error case GetPriceInfoApi()", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?money_received=300&bouns=100&fixed_cost=notanumber&loan=50&private=50", nil)

		paramMap, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private")

		// GetPriceInfoApi 関数を呼び出し
		controllers.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, c.Writer.Status())

		assert.Empty(t, paramMap)
		expectedErrorMessage := "strconv.Atoi: parsing \"notanumber\": invalid syntax"
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("paramMap[]: %v", paramMap)
		t.Logf("err: %s", err)
	})
}
