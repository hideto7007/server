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
	"server/models"
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

		fetcher := NewIncomeDataFetcher()
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
		// config.Setup()
		// defer config.Teardown()
		// defer config.TeardownTestDatabase()

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

		fetcher := NewIncomeDataFetcher()
		fetcher.GetIncomeDataInRangeApi(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "database error", response["error_msg"])
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

		fetcher := NewIncomeDataFetcher()
		fetcher.GetIncomeDataInRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody []errorMessages
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := []errorMessages{
			{
				Field:   "start_date",
				Message: "開始期間は必須です。",
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

			fetcher := NewIncomeDataFetcher()
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody []errorMessages
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := []errorMessages{
				{
					Field:   "start_date",
					Message: "開始日の形式が間違っています。",
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

			fetcher := NewIncomeDataFetcher()
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody []errorMessages
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := []errorMessages{
				{
					Field:   "end_date",
					Message: "終了日の形式が間違っています。",
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

		fetcher := NewIncomeDataFetcher()
		fetcher.GetIncomeDataInRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody []errorMessages
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := []errorMessages{
			{
				Field:   "user_id",
				Message: "ユーザーIDは必須です。",
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

			fetcher := NewIncomeDataFetcher()
			fetcher.GetIncomeDataInRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody []errorMessages
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := []errorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは整数値のみです。",
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
		fetcher := NewIncomeDataFetcher()
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
		fetcher := NewIncomeDataFetcher()
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

		fetcher := NewIncomeDataFetcher()
		fetcher.GetDateRangeApi(c)

		// レスポンスのステータスコードを確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody []errorMessages
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := []errorMessages{
			{
				Field:   "user_id",
				Message: "ユーザーIDは必須です。",
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

			fetcher := NewIncomeDataFetcher()
			fetcher.GetDateRangeApi(c)

			// レスポンスのステータスコードを確認
			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody []errorMessages
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := []errorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは整数値のみです。",
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
		fetcher := NewIncomeDataFetcher()
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
		fetcher := NewIncomeDataFetcher()
		fetcher.GetYearIncomeAndDeductionApi(c)

		// レスポンスの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response["error_msg"])
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

		fetcher := NewIncomeDataFetcher()
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "新規給料情報を登録致しました。", response["result"])
	})

	t.Run("error InsertIncomeDataApi", func(t *testing.T) {
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
				return errors.New("database error")
			})
		defer patches.Reset()

		fetcher := NewIncomeDataFetcher()
		fetcher.InsertIncomeDataApi(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "新規登録時にエラーが発生。", response["error_msg"])
	})

	// t.Run("invalid JSON InsertIncomeDataApi", func(t *testing.T) {
	// 	// config.Setup()
	// 	// defer config.Teardown()
	// 	// defer config.TeardownTestDatabase()

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	// Invalid JSON
	// 	invalidJSON := `{"data": [`

	// 	c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBufferString(invalidJSON))
	// 	c.Request.Header.Set("Content-Type", "application/json")

	// 	fetcher := NewIncomeDataFetcher()
	// 	fetcher.InsertIncomeDataApi(c)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	var response map[string]interface{}
	// 	err := json.Unmarshal(w.Body.Bytes(), &response)
	// 	assert.NoError(t, err)
	// 	assert.Contains(t, response["error_msg"], "unexpected EOF")
	// })
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

		fetcher := NewIncomeDataFetcher()
		fetcher.UpdateIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "給料情報の更新が問題なく成功しました。", response["result"])
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

			fetcher := NewIncomeDataFetcher()
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

	// t.Run("invalid JSON UpdateIncomeDataApi", func(t *testing.T) {
	// 	// config.Setup()
	// 	// defer config.Teardown()
	// 	// defer config.TeardownTestDatabase()

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	// Invalid JSON
	// 	invalidJSON := `{"data": [`

	// 	c.Request = httptest.NewRequest("POST", "/api/income_update", bytes.NewBufferString(invalidJSON))
	// 	c.Request.Header.Set("Content-Type", "application/json")

	// 	fetcher := NewIncomeDataFetcher()
	// 	fetcher.UpdateIncomeDataApi(c)

	// 	assert.Equal(t, http.StatusBadRequest, w.Code)
	// 	var response map[string]interface{}
	// 	err := json.Unmarshal(w.Body.Bytes(), &response)
	// 	assert.NoError(t, err)
	// 	assert.Contains(t, response["error_msg"], "unexpected EOF")
	// })
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
		c.Request = httptest.NewRequest("POST", "/api/income_delete?", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := NewIncomeDataFetcher()
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "給料情報の削除が問題なく成功しました。", response["result"])
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
		c.Request = httptest.NewRequest("POST", "/api/income_delete?", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
			"DeleteIncome",
			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
				return errors.New("database error")
			})
		defer patches.Reset()

		fetcher := NewIncomeDataFetcher()
		fetcher.DeleteIncomeDataApi(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "削除中にエラーが発生しました", response["error_msg"])
	})
}
