package controllers

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "os"
    "reflect"
    "server/models"
    "testing"
    "time"

    . "github.com/agiledragon/gomonkey/v2"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func setup() {
    os.Setenv("DB_HOST", "localhost") // テスト用のデータベースホストを設定
}

func teardown() {
    os.Unsetenv("DB_HOST")
}

func TestGetIncomeDataInRangeApi(t *testing.T) {
    setup()
    defer teardown()

    gin.SetMode(gin.TestMode)

    t.Run("success GetIncomeDataInRangeApi", func(t *testing.T) {
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

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "GetIncomeDataInRange", func(_ *models.PostgreSQLDataFetcher, startDate, endDate, userId string) ([]models.IncomeData, error) {
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
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/?start_date=2022-07-01&end_date=2022-09-30&user_id=1", nil)

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "GetIncomeDataInRange", func(_ *models.PostgreSQLDataFetcher, startDate, endDate, userId string) ([]models.IncomeData, error) {
            return nil, errors.New("database error")
        })
        defer patches.Reset()

        fetcher := NewIncomeDataFetcher()
        fetcher.GetIncomeDataInRangeApi(c)

        assert.Equal(t, http.StatusInternalServerError, w.Code)
        var response map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.Equal(t, "database error", response["error"])
    })
}

func TestGetDateRangeApi(t *testing.T) {
    setup()
    defer teardown()

    gin.SetMode(gin.TestMode)

    t.Run("success GetDateRangeApi", func(t *testing.T) {
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

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "GetDateRange", func(_ *models.PostgreSQLDataFetcher, UserID string) ([]models.PaymentDate, error) {
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
        // テスト用のGinコンテキストを作成
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "GetDateRange", func(_ *models.PostgreSQLDataFetcher, UserID string) ([]models.PaymentDate, error) {
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
        assert.NotNil(t, response["error"])
    })
}

func TestUpdateIncomeDataApi(t *testing.T) {
    setup()
    defer teardown()

    gin.SetMode(gin.TestMode)

    t.Run("success UpdateIncomeDataApi", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        testData := []models.UpdateIncomeData{
            {
                IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
                PaymentDate:      "2024-02-10",
                Age:              30,
                Industry:         "IT",
                TotalAmount:      320524,
                DeductionAmount:  93480,
                TakeHomeAmount:   227044,
                Classification:   "給料",
            },
        }
        var Body struct {
            Data []models.UpdateIncomeData `json:"data"`
        }
        Body.Data = testData

        body, _ := json.Marshal(Body)
        c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
        c.Request.Header.Set("Content-Type", "application/json")

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "UpdateIncome", func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
            return nil
        })
        defer patches.Reset()

        fetcher := NewIncomeDataFetcher()
        fetcher.UpdateIncomeDataApi(c)

        assert.Equal(t, http.StatusOK, w.Code)
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "給料情報の更新が問題なく成功しました。", response["message"])
    })

    t.Run("error UpdateIncomeDataApi", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        testData := []models.UpdateIncomeData{
            {
                IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
                PaymentDate:      "2024-02-10",
                Age:              30,
                Industry:         "IT",
                TotalAmount:      320524,
                DeductionAmount:  93480,
                TakeHomeAmount:   227044,
                Classification:   "給料",
            },
        }
        var Body struct {
            Data []models.UpdateIncomeData `json:"data"`
        }
        Body.Data = testData

        body, _ := json.Marshal(Body)
        c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
        c.Request.Header.Set("Content-Type", "application/json")

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "UpdateIncome", func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
            return errors.New("database error")
        })
        defer patches.Reset()

        fetcher := NewIncomeDataFetcher()
        fetcher.UpdateIncomeDataApi(c)

        assert.Equal(t, http.StatusInternalServerError, w.Code)
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "データベースへの挿入中にエラーが発生しました", response["error"])
    })
    t.Run("invalid JSON UpdateIncomeDataApi", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)

        // Invalid JSON
        invalidJSON := `{"data": [`

        c.Request = httptest.NewRequest("POST", "/api/income_update", bytes.NewBufferString(invalidJSON))
        c.Request.Header.Set("Content-Type", "application/json")

        fetcher := NewIncomeDataFetcher()
        fetcher.UpdateIncomeDataApi(c)

        assert.Equal(t, http.StatusBadRequest, w.Code)
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Contains(t, response["error"], "unexpected EOF")
    })
}

func TestDeleteIncomeDataApi(t *testing.T) {
    setup()
    defer teardown()

    gin.SetMode(gin.TestMode)

    t.Run("success DeleteIncomeDataApi", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("DELETE", "/api/income_delete?income_forecast_id=7b941edb-b7a2-e1e7-6466-ce53d1c8bcff", nil)

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "DeleteIncome", func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
            return nil
        })
        defer patches.Reset()

        fetcher := NewIncomeDataFetcher()
        fetcher.DeleteIncomeDataApi(c)

        assert.Equal(t, http.StatusOK, w.Code)
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "給料情報の削除が問題なく成功しました。", response["message"])
    })

    t.Run("error DeleteIncomeDataApi", func(t *testing.T) {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("DELETE", "/api/income_delete?income_forecast_id=7b941edb-b7a2-e1e7-6466-ce53d1c8bcff", nil)

        patches := ApplyMethod(reflect.TypeOf(&models.PostgreSQLDataFetcher{}), "DeleteIncome", func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
            return errors.New("database error")
        })
        defer patches.Reset()

        fetcher := NewIncomeDataFetcher()
        fetcher.DeleteIncomeDataApi(c)

        assert.Equal(t, http.StatusInternalServerError, w.Code)
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "データベースからの削除中にエラーが発生しました", response["error"])
    })
}

func jsonReader(data interface{}) *bytes.Reader {
    body, _ := json.Marshal(data)
    return bytes.NewReader(body)
}
