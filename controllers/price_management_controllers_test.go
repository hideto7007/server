package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"server/common"
	"server/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestRes struct {
}

func TestPriceCalc(t *testing.T) {
	t.Run("success PriceCalc()", func(t *testing.T) {
		// テストケース1: 正常な整数の計算

		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		result := pm.PriceCalc(300, 100, 50, 50, 50, 30)

		assert.Equal(t, 150, result.LeftAmount)
		assert.Equal(t, 1870, result.TotalAmount)

		t.Logf("result.LeftAmount: %d", result.LeftAmount)
		t.Logf("result.TotalAmount: %d", result.TotalAmount)
	})
}

func TestGetPriceInfo(t *testing.T) {
	t.Run("success GetPriceInfoApi()", func(t *testing.T) {
		// テスト用のGinコンテキストを作成
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?money_received=300&bouns=100&fixed_cost=50&loan=50&private=50&insurance=30", nil)
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}

		paramMap, err := pm.CommonFetcher.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private", "insurance")

		res := pm.PriceCalc(paramMap["money_received"], paramMap["bouns"], paramMap["fixed_cost"], paramMap["loan"], paramMap["private"], paramMap["insurance"])

		// GetPriceInfoApi 関数を呼び出し
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusOK, c.Writer.Status())

		// レスポンスの JSON データを取得
		response := utils.Response[PriceInfo]{
			Result: []PriceInfo{
				{
					LeftAmount:  res.LeftAmount,
					TotalAmount: res.TotalAmount,
				},
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, 150, response.Result[0].LeftAmount)
		assert.Equal(t, 1870, response.Result[0].TotalAmount)

		t.Logf("response.PriceInfo.LeftAmount: %d", response.Result[0].LeftAmount)
		t.Logf("response.PriceInfo.TotalAmount: %d", response.Result[0].TotalAmount)
		t.Logf("err: %v", err)
	})

	t.Run("バリデーションエラー money_received", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=test&bouns=100&fixed_cost=100&loan=50&private=50&insurance=30", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "money_received",
					Message: "月の収入は整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー bouns", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=100&bouns=test&fixed_cost=100&loan=50&private=50&insurance=30", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "bouns",
					Message: "ボーナスは整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー fixed_cost", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=100&bouns=100&fixed_cost=test&loan=50&private=50&insurance=30", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "fixed_cost",
					Message: "月の収入は整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー loan", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=100&bouns=100&fixed_cost=100&loan=test&private=50&insurance=30", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "loan",
					Message: "ローンは整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー private", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=100&bouns=100&fixed_cost=100&loan=100&private=test&insurance=30", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "private",
					Message: "プライベートは整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー insurance", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?money_received=100&bouns=100&fixed_cost=100&loan=100&private=100&insurance=test", nil)

		// GetPriceInfoApi 関数を呼び出し
		pm := apiPriceManagementFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		pm.GetPriceInfoApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "insurance",
					Message: "保険は整数値のみです。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})
}
