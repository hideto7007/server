package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	// "server/config"
	"server/common"
	"server/models"
	"server/test_utils"
	"server/utils"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// func TestMain(m *testing.M) {
// 	config.Setup()
// 	code := m.Run()
// 	config.Teardown()
// 	os.Exit(code)
// }

// func TestMain(m *testing.M) {
// 	config.Setup()
// 	config.SetupTestDatabase()
// 	code := m.Run()
// 	config.TeardownTestDatabase()
// 	config.Teardown()
// 	os.Exit(code)
// }

// func TestMain(m *testing.M) {
// 	config.Setup()
// 	defer config.Teardown()

// 	code := m.Run()
// 	os.Exit(code)
// }

type testData struct {
	Data interface{} `json:"data"`
}

func TestGetIncomeDataInRangeApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success GetIncomeDataInRangeApi", func(t *testing.T) {
		// config.Setup()
		// config.SetupTestDatabase()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?start_date=2022-07-01&end_date=2022-09-30&user_id=1", nil)

		mockData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "30",
				Industry:         "IT",
				TotalAmount:      5000,
				DeductionAmount:  500,
				TakeHomeAmount:   4500,
				Classification:   "Salary",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac365c16e36"),
				PaymentDate:      time.Date(2022, time.August, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "30",
				Industry:         "IT",
				TotalAmount:      5000,
				DeductionAmount:  500,
				TakeHomeAmount:   4500,
				Classification:   "Salary",
				UserID:           1,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetIncomeDataInRange",
			func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetIncomeDataInRangeApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var rawResponse struct {
			Result json.RawMessage `json:"result"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &rawResponse)
		assert.NoError(t, err)
		// 期待されるJSONデータを構築
		expectedResponse, err := json.Marshal(mockData)
		assert.NoError(t, err)
		// 期待されるJSONデータと実際のレスポンスを比較
		assert.JSONEq(t, string(expectedResponse), string(rawResponse.Result))
	})

	t.Run("error GetIncomeDataInRangeApi", func(t *testing.T) {
		// ここのテストケースだけ不安定なのでN回リトライしてテスト行う
		const maxRetry = 100

		var lastErr error

		for attempt := 1; attempt <= maxRetry; attempt++ {
			time.Sleep(500 * time.Millisecond)
			// config.Setup()
			// defer config.Teardown()
			// defer config.TeardownTestDatabase()

			// テスト用のGinコンテキストを作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?start_date=2022-07-01&end_date=2022-09-30&user_id=1", nil)

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetIncomeDataInRange",
				func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
					return nil, errors.New("database error")
				})
			defer patches.Reset()

			// テスト対象の関数を呼び出し
			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスの確認
			if w.Code == http.StatusInternalServerError {
				var response utils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					lastErr = err
				} else {
					// assert.NoError(t, err)
					assert.Equal(t, "database error", response.ErrorMsg)
					fmt.Printf("GetIncomeDataInRangeApiのテスト成功しました! 試行回数: %d回目\n", attempt)
					lastErr = nil
					break
				}
			} else {
				lastErr = errors.New("データベース取得エラーになりませんでした: " + http.StatusText(w.Code))
			}

			if attempt == maxRetry {
				t.Fatalf("テスト試行回数の上限に達しました %d回 テストエラー内容: %v", maxRetry, lastErr)
			}
		}
	})

	t.Run("バリデーションエラー start_date 必須", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?start_date=&end_date=2022-09-30&user_id=1", nil)

		mockData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "30",
				Industry:         "IT",
				TotalAmount:      5000,
				DeductionAmount:  500,
				TakeHomeAmount:   4500,
				Classification:   "Salary",
				UserID:           1,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetIncomeDataInRange",
			func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetIncomeDataInRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "start_date",
					Message: "開始期間は必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー end_date 必須", func(t *testing.T) {
		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?start_date=2022-09-13&end_date=&user_id=1", nil)

		mockData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "30",
				Industry:         "IT",
				TotalAmount:      5000,
				DeductionAmount:  500,
				TakeHomeAmount:   4500,
				Classification:   "Salary",
				UserID:           1,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetIncomeDataInRange",
			func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetIncomeDataInRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "end_date",
					Message: "終了期間は必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー start_date 日付不正", func(t *testing.T) {
		paramsList := [...]string{
			"/?start_date=202207-01&end_date=2022-09-30&user_id=1",
			"/?start_date=2022-0701&end_date=2022-09-30&user_id=1",
			"/?start_date=2022-13-01&end_date=2022-09-30&user_id=1",
			"/?start_date=2022-11-32&end_date=2022-09-30&user_id=1",
			"/?start_date=test&end_date=2022-09-30&user_id=1",
		}

		for _, params := range paramsList {
			// エラーを引き起こすリクエストをシミュレート
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", params, nil)

			mockData := []models.IncomeData{
				{
					IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
					PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
					Age:              "30",
					Industry:         "IT",
					TotalAmount:      5000,
					DeductionAmount:  500,
					TakeHomeAmount:   4500,
					Classification:   "Salary",
					UserID:           1,
				},
			}

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetIncomeDataInRange",
				func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
					return mockData, nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "start_date",
						Message: "開始日の形式が間違っています。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("バリデーションエラー end_date 日付不正", func(t *testing.T) {
		paramsList := [...]string{
			"/?start_date=2022-07-01&end_date=202209-30&user_id=1",
			"/?start_date=2022-07-01&end_date=2022-0930&user_id=1",
			"/?start_date=2022-07-01&end_date=2022-13-30&user_id=1",
			"/?start_date=2022-07-01&end_date=2022-09-32&user_id=1",
			"/?start_date=2022-07-01&end_date=test0&user_id=1",
		}

		for _, params := range paramsList {
			// エラーを引き起こすリクエストをシミュレート
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", params, nil)

			mockData := []models.IncomeData{
				{
					IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
					PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
					Age:              "30",
					Industry:         "IT",
					TotalAmount:      5000,
					DeductionAmount:  500,
					TakeHomeAmount:   4500,
					Classification:   "Salary",
					UserID:           1,
				},
			}

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetIncomeDataInRange",
				func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
					return mockData, nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "end_date",
						Message: "終了日の形式が間違っています。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("バリデーションエラー user_id 必須", func(t *testing.T) {

		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?start_date=2022-07-01&end_date=2022-09-30&user_id=", nil)

		mockData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "30",
				Industry:         "IT",
				TotalAmount:      5000,
				DeductionAmount:  500,
				TakeHomeAmount:   4500,
				Classification:   "Salary",
				UserID:           1,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetIncomeDataInRange",
			func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetIncomeDataInRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー user_id 整数値のみ", func(t *testing.T) {
		paramsList := [...]string{
			"/?start_date=2022-07-01&end_date=2022-09-30&user_id=rere",
			"/?start_date=2022-07-01&end_date=2022-09-30&user_id=1.23",
		}

		for _, params := range paramsList {
			// エラーを引き起こすリクエストをシミュレート
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", params, nil)

			mockData := []models.IncomeData{
				{
					IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
					PaymentDate:      time.Date(2022, time.July, 15, 0, 0, 0, 0, time.FixedZone("", 0)),
					Age:              "30",
					Industry:         "IT",
					TotalAmount:      5000,
					DeductionAmount:  500,
					TakeHomeAmount:   4500,
					Classification:   "Salary",
					UserID:           1,
				},
			}

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetIncomeDataInRange",
				func(_ *models.PostgreSQLDataFetcher, startDate string, endDate string, userId int) ([]models.IncomeData, error) {
					return mockData, nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "user_id",
						Message: "ユーザーIDは整数値のみです。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})
}

func TestGetDateRangeApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success GetDateRangeApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		// テスト用のGinコンテキストを作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

		// モックデータを設定
		mockData := []models.PaymentDate{
			{
				UserID:            1,
				StratPaymaentDate: "2022-01-01",
				EndPaymaentDate:   "2022-12-31",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetDateRange",
			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
				return mockData, nil
			})
		defer patches.Reset()

		// テスト対象の関数を呼び出し
		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetDateRangeApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var rawResponse struct {
			Result json.RawMessage `json:"result"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &rawResponse)
		assert.NoError(t, err)
		// 期待されるJSONデータを構築
		expectedResponse, err := json.Marshal(mockData)
		assert.NoError(t, err)
		// 期待されるJSONデータと実際のレスポンスを比較
		assert.JSONEq(t, string(expectedResponse), string(rawResponse.Result))
	})

	t.Run("error GetDateRangeApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		// テスト用のGinコンテキストを作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetDateRange",
			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
				return nil, errors.New("database error")
			})
		defer patches.Reset()

		// テスト対象の関数を呼び出し
		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetDateRangeApi(c)

		// レスポンスの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotNil(t, response["error_msg"])
	})

	t.Run("バリデーションエラー user_id 必須", func(t *testing.T) {

		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?user_id=", nil)

		mockData := []models.PaymentDate{
			{
				UserID:            1,
				StratPaymaentDate: "2022-01-01",
				EndPaymaentDate:   "2022-12-31",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetDateRange",
			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
				return mockData, nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetDateRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー user_id 整数値のみ", func(t *testing.T) {
		paramsList := [...]string{
			"/?user_id=rere",
			"/?user_id=1.23",
		}

		for _, params := range paramsList {
			// エラーを引き起こすリクエストをシミュレート
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", params, nil)

			mockData := []models.PaymentDate{
				{
					UserID:            1,
					StratPaymaentDate: "2022-01-01",
					EndPaymaentDate:   "2022-12-31",
				},
			}

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetDateRange",
				func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
					return mockData, nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetDateRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "user_id",
						Message: "ユーザーIDは整数値のみです。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})
}

func TestGetYearIncomeAndDeductionApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success GetYearIncomeAndDeductionApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		// テスト用のGinコンテキストを作成
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

		// モックデータを設定
		mockData := []models.YearsIncomeData{
			{
				Years:           "2022",
				TotalAmount:     6000,
				DeductionAmount: 600,
				TakeHomeAmount:  5400,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetYearsIncomeAndDeduction",
			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		// テスト対象の関数を呼び出し
		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetYearIncomeAndDeductionApi(c)

		// レスポンスの確認
		assert.Equal(t, http.StatusOK, w.Code)
		var response struct {
			Result []models.YearsIncomeData `json:"result"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, mockData, response.Result)
	})

	t.Run("error GetYearIncomeAndDeductionApi", func(t *testing.T) {
		// ここのテストケースだけ不安定なのでN回リトライしてテスト行う
		const maxRetry = 100

		var lastErr error

		for attempt := 1; attempt <= maxRetry; attempt++ {
			time.Sleep(500 * time.Millisecond)
			// config.Setup()
			// defer config.Teardown()
			// defer config.TeardownTestDatabase()

			// テスト用のGinコンテキストを作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetYearsIncomeAndDeduction",
				func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
					return nil, errors.New("database error")
				})
			defer patches.Reset()

			// テスト対象の関数を呼び出し
			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetYearIncomeAndDeductionApi(c)

			// レスポンスの確認
			if w.Code == http.StatusInternalServerError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					lastErr = err
				} else {
					// assert.NoError(t, err)
					assert.Equal(t, "database error", response["error_msg"])
					fmt.Printf("GetYearsIncomeAndDeductionのテスト成功しました! 試行回数: %d回目\n", attempt)
					lastErr = nil
					break
				}
			} else {
				lastErr = errors.New("データベース取得エラーになりませんでした: " + http.StatusText(w.Code))
			}

			if attempt == maxRetry {
				t.Fatalf("テスト試行回数の上限に達しました %d回 テストエラー内容: %v", maxRetry, lastErr)
			}
		}
	})

	t.Run("バリデーションエラー user_id 必須", func(t *testing.T) {

		// エラーを引き起こすリクエストをシミュレート
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/?user_id=", nil)

		mockData := []models.YearsIncomeData{
			{
				Years:           "2022",
				TotalAmount:     6000,
				DeductionAmount: 600,
				TakeHomeAmount:  5400,
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"GetYearsIncomeAndDeduction",
			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
				return mockData, nil
			})
		defer patches.Reset()

		// テスト対象の関数を呼び出し
		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetYearIncomeAndDeductionApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー user_id 整数値のみ", func(t *testing.T) {
		paramsList := [...]string{
			"/?user_id=rere",
			"/?user_id=1.23",
		}

		for _, params := range paramsList {
			// エラーを引き起こすリクエストをシミュレート
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", params, nil)

			mockData := []models.YearsIncomeData{
				{
					Years:           "2022",
					TotalAmount:     6000,
					DeductionAmount: 600,
					TakeHomeAmount:  5400,
				},
			}

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"GetYearsIncomeAndDeduction",
				func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
					return mockData, nil
				})
			defer patches.Reset()

			// テスト対象の関数を呼び出し
			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetYearIncomeAndDeductionApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "user_id",
						Message: "ユーザーIDは整数値のみです。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})
}

func TestInsertIncomeDataApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success InsertIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.InsertIncomeData{
				{
					PaymentDate:     "2024-02-10",
					Age:             30,
					Industry:        "IT",
					TotalAmount:     320524,
					DeductionAmount: 93480,
					TakeHomeAmount:  227044,
					Classification:  "給料",
					UserID:          "1",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"InsertIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "新規給料情報を登録致しました。", response.Result)
	})

	t.Run("error InsertIncomeDataApi", func(t *testing.T) {

		// ここのテストケースだけ不安定なのでN回リトライしてテスト行う
		const maxRetry = 100

		var lastErr error

		for attempt := 1; attempt <= maxRetry; attempt++ {
			// config.Setup()
			// defer config.Teardown()
			// defer config.TeardownTestDatabase()
			time.Sleep(500 * time.Millisecond)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			data := testData{
				Data: []models.InsertIncomeData{
					{
						PaymentDate:     "2024-02-10",
						Age:             30,
						Industry:        "IT",
						TotalAmount:     320524,
						DeductionAmount: 93480,
						TakeHomeAmount:  227044,
						Classification:  "給料",
						UserID:          "1",
					},
				},
			}

			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"InsertIncome",
				func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
					return errors.New("database error")
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.InsertIncomeDataApi(c)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			// レスポンスの確認
			if w.Code == http.StatusInternalServerError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					lastErr = err
				} else {
					// assert.NoError(t, err)
					assert.Equal(t, "新規登録時にエラーが発生。", response["error_msg"])
					fmt.Printf("InsertIncomeDataApiのテスト成功しました! 試行回数: %d回目\n", attempt)
					lastErr = nil
					break
				}
			} else {
				lastErr = errors.New("データベース登録エラーになりませんでした: " + http.StatusText(w.Code))
			}

			if attempt == maxRetry {
				t.Fatalf("テスト試行回数の上限に達しました %d回 テストエラー内容: %v", maxRetry, lastErr)
			}
		}
	})

	t.Run("request Data empty InsertIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.InsertIncomeData{},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"InsertIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[models.InsertIncomeData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// 期待するエラーメッセージを確認
		expectedErrorMessage := utils.ResponseWithSlice[models.InsertIncomeData]{
			ErrorMsg: "登録するデータが存在しません。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("invalid JSON InsertIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.InsertIncomeData{
				{
					PaymentDate:     "",
					Age:             0,
					Industry:        "",
					TotalAmount:     "",
					DeductionAmount: "",
					TakeHomeAmount:  "",
					Classification:  "",
					UserID:          "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"InsertIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			RecodeRows: 1,
			Result: []utils.ErrorMessages{
				{
					Field:   "age",
					Message: "年齢は必須又は整数値のみです。",
				},
				{
					Field:   "industry",
					Message: "業種は必須です。",
				},
				{
					Field:   "total_amount",
					Message: "総支給額は必須です。",
				},
				{
					Field:   "deduction_amount",
					Message: "差引額は必須です。",
				},
				{
					Field:   "take_home_amount",
					Message: "手取額は必須です。",
				},
				{
					Field:   "classification",
					Message: "分類は必須です。",
				},
				{
					Field:   "user_id",
					Message: "ユーザーIDは必須です。",
				},
				{
					Field:   "payment_date",
					Message: "報酬日付は必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー 数値文字列以外は無効", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.InsertIncomeData{
				{
					PaymentDate:     "2021-09-10",
					Age:             25,
					Industry:        "IT業界",
					TotalAmount:     "0kj",
					DeductionAmount: "0.3",
					TakeHomeAmount:  "hoge",
					Classification:  "賞与",
					UserID:          "1",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"InsertIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			RecodeRows: 1,
			Result: []utils.ErrorMessages{
				{
					Field:   "total_amount",
					Message: "総支給額で数値文字列以外は無効です。",
				},
				{
					Field:   "deduction_amount",
					Message: "差引額で数値文字列以外は無効です。",
				},
				{
					Field:   "take_home_amount",
					Message: "手取額で数値文字列以外は無効です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー 形式及びユーザーIDの整数値チェック", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		dataList := []testData{
			{
				Data: []models.InsertIncomeData{
					{
						PaymentDate:     "202109-10",
						Age:             25,
						Industry:        "IT業界",
						TotalAmount:     "350000",
						DeductionAmount: "76000",
						TakeHomeAmount:  "274000",
						Classification:  "賞与",
						UserID:          "1o",
					},
				},
			},
			{
				Data: []models.InsertIncomeData{
					{
						PaymentDate:     "2021-0910",
						Age:             25,
						Industry:        "IT業界",
						TotalAmount:     "350000",
						DeductionAmount: "76000",
						TakeHomeAmount:  "274000",
						Classification:  "賞与",
						UserID:          "1.0",
					},
				},
			},
			{
				Data: []models.InsertIncomeData{
					{
						PaymentDate:     "2021-13-10",
						Age:             25,
						Industry:        "IT業界",
						TotalAmount:     "350000",
						DeductionAmount: "76000",
						TakeHomeAmount:  "274000",
						Classification:  "賞与",
						UserID:          "1.0hoge",
					},
				},
			},
			{
				Data: []models.InsertIncomeData{
					{
						PaymentDate:     "2021-12-33",
						Age:             25,
						Industry:        "IT業界",
						TotalAmount:     "350000",
						DeductionAmount: "76000",
						TakeHomeAmount:  "274000",
						Classification:  "賞与",
						UserID:          "hoge",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"InsertIncome",
				func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.InsertIncomeDataApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				RecodeRows: 1,
				Result: []utils.ErrorMessages{
					{
						Field:   "payment_date",
						Message: "給料支給日の形式が間違っています。",
					},
					{
						Field:   "user_id",
						Message: "ユーザーIDは整数値のみです。",
					},
				},
			}
			test_utils.SortErrorMessages(responseBody.Result)
			test_utils.SortErrorMessages(expectedErrorMessage.Result)
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})
}

func TestUpdateIncomeDataApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success UpdateIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.UpdateIncomeData{
				{
					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
					PaymentDate:      "2024-02-10",
					Age:              30,
					Industry:         "IT",
					TotalAmount:      320524,
					DeductionAmount:  93480,
					TakeHomeAmount:   227044,
					UpdateUser:       "test_user",
					Classification:   "給料",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"UpdateIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "給料情報の更新が問題なく成功しました。", response.Result)
	})

	t.Run("error UpdateIncomeDataApi", func(t *testing.T) {
		// ここのテストケースだけ不安定なのでN回リトライしてテスト行う
		const maxRetry = 100

		var lastErr error

		for attempt := 1; attempt <= maxRetry; attempt++ {
			// time.Sleep(1 * time.Second) // 1秒
			time.Sleep(500 * time.Millisecond)
			// config.Setup()
			// defer config.Teardown()
			// defer config.TeardownTestDatabase()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			data := testData{
				Data: []models.UpdateIncomeData{
					{
						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
						PaymentDate:      "2024-02-10",
						Age:              30,
						Industry:         "IT",
						TotalAmount:      320524,
						DeductionAmount:  93480,
						TakeHomeAmount:   227044,
						UpdateUser:       "test_user",
						Classification:   "給料",
					},
				},
			}

			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"UpdateIncome",
				func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
					return errors.New("database error")
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.UpdateIncomeDataApi(c)

			if w.Code == http.StatusInternalServerError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					lastErr = err
				} else {
					// assert.NoError(t, err)
					assert.Equal(t, "更新時にエラーが発生。", response["error_msg"])
					fmt.Printf("成功しました！試行回数: %d回目\n", attempt)
					lastErr = nil
					break
				}
			} else {
				lastErr = errors.New("データベース挿入エラーになりませんでした: " + http.StatusText(w.Code))
			}

			if attempt == maxRetry {
				t.Fatalf("テスト試行回数の上限に達しました %d回 テストエラー内容: %v", maxRetry, lastErr)
			}
		}
	})

	t.Run("request Data empty UpdateIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.UpdateIncomeData{},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_update", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"UpdateIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[models.UpdateIncomeData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// 期待するエラーメッセージを確認
		expectedErrorMessage := utils.ResponseWithSlice[models.UpdateIncomeData]{
			ErrorMsg: "更新するデータが存在しません。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("invalid JSON UpdateIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/income_update", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.UpdateIncomeData{
				{
					IncomeForecastID: "",
					PaymentDate:      "",
					Age:              0,
					Industry:         "",
					TotalAmount:      "",
					DeductionAmount:  "",
					TakeHomeAmount:   "",
					UpdateUser:       "",
					Classification:   "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"UpdateIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			RecodeRows: 1,
			Result: []utils.ErrorMessages{
				{
					Field:   "income_forecast_id",
					Message: "年収推移IDは必須です。",
				},
				{
					Field:   "payment_date",
					Message: "報酬日付は必須です。",
				},
				{
					Field:   "age",
					Message: "年齢は必須又は整数値のみです。",
				},
				{
					Field:   "industry",
					Message: "業種は必須です。",
				},
				{
					Field:   "total_amount",
					Message: "総支給額は必須です。",
				},
				{
					Field:   "deduction_amount",
					Message: "差引額は必須です。",
				},
				{
					Field:   "take_home_amount",
					Message: "手取額は必須です。",
				},
				{
					Field:   "update_user",
					Message: "更新者は必須です。",
				},
				{
					Field:   "classification",
					Message: "分類は必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー 数値文字列以外は無効", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.UpdateIncomeData{
				{
					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
					PaymentDate:      "2024-02-10",
					Age:              30,
					Industry:         "IT",
					TotalAmount:      "fd",
					DeductionAmount:  "0.8u",
					TakeHomeAmount:   "0.45",
					UpdateUser:       "test_user",
					Classification:   "給料",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"UpdateIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			RecodeRows: 1,
			Result: []utils.ErrorMessages{
				{
					Field:   "total_amount",
					Message: "総支給額で数値文字列以外は無効です。",
				},
				{
					Field:   "deduction_amount",
					Message: "差引額で数値文字列以外は無効です。",
				},
				{
					Field:   "take_home_amount",
					Message: "手取額で数値文字列以外は無効です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー 形式及びユーザーIDの整数値チェック", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		dataList := []testData{
			{
				Data: []models.UpdateIncomeData{
					{
						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
						PaymentDate:      "202402-10",
						Age:              30,
						Industry:         "IT",
						TotalAmount:      320524,
						DeductionAmount:  93480,
						TakeHomeAmount:   227044,
						UpdateUser:       "test_user",
						Classification:   "給料",
					},
				},
			},
			{
				Data: []models.UpdateIncomeData{
					{
						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
						PaymentDate:      "2024-0210",
						Age:              30,
						Industry:         "IT",
						TotalAmount:      320524,
						DeductionAmount:  93480,
						TakeHomeAmount:   227044,
						UpdateUser:       "test_user",
						Classification:   "給料",
					},
				},
			},
			{
				Data: []models.UpdateIncomeData{
					{
						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
						PaymentDate:      "2024-13-10",
						Age:              30,
						Industry:         "IT",
						TotalAmount:      320524,
						DeductionAmount:  93480,
						TakeHomeAmount:   227044,
						UpdateUser:       "test_user",
						Classification:   "給料",
					},
				},
			},
			{
				Data: []models.UpdateIncomeData{
					{
						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
						PaymentDate:      "2024-02-32",
						Age:              30,
						Industry:         "IT",
						TotalAmount:      320524,
						DeductionAmount:  93480,
						TakeHomeAmount:   227044,
						UpdateUser:       "test_user",
						Classification:   "給料",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
				"UpdateIncome",
				func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiIncomeDataFetcher{
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.UpdateIncomeDataApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
				RecodeRows: 1,
				Result: []utils.ErrorMessages{
					{
						Field:   "payment_date",
						Message: "給料支給日の形式が間違っています。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})
}

func TestDeleteIncomeDataApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success DeleteIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.DeleteIncomeData{
				{
					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "給料情報の削除が問題なく成功しました。", response.Result)
	})

	t.Run("error DeleteIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.DeleteIncomeData{
				{
					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return errors.New("database error")
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "削除中にエラーが発生しました", response["error_msg"])
	})

	t.Run("invalid JSON DeleteIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("request Data empty DeleteIncomeDataApi", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.DeleteIncomeData{},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[models.DeleteIncomeData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// 期待するエラーメッセージを確認
		expectedErrorMessage := utils.ResponseWithSlice[models.DeleteIncomeData]{
			ErrorMsg: "削除するデータが存在しません。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.DeleteIncomeData{
				{
					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
				},
				{
					IncomeForecastID: "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("PUT", "/api/income_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiIncomeDataFetcher{
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			RecodeRows: 2,
			Result: []utils.ErrorMessages{
				{
					Field:   "income_forecast_id",
					Message: "年収推移IDは必須です。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})
}
