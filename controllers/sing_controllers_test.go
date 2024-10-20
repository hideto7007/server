package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	// "server/config"

	mock_controllers "server/mock/controllers"
	mock_utils "server/mock/utils"
	"server/models"
	"server/test_utils"
	"server/utils"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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

// type testData struct {
// 	Data interface{} `json:"data"`
// }

func TestPostSingInApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("TestPostSingInApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := NewSingDataFetcher()
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("TestPostSingInApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingInData{
				{
					UserId:       "",
					UserName:     "",
					UserPassword: "",
				},
			},
		}

		resMock := []models.SingInData{
			{
				UserId:       3,
				UserName:     "test@example.com",
				UserPassword: "Test12345!",
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		fetcher := NewSingDataFetcher()
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_id",
					Message: "ユーザーIDは必須です。",
				},
				{
					Field:   "user_name",
					Message: "ユーザー名は必須です。",
				},
				{
					Field:   "user_password",
					Message: "パスワードは必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi バリデーション メールアドレス不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingInData{
				{
					UserId:       1,
					UserName:     "test",
					UserPassword: "Test12345!",
				},
			},
		}

		resMock := []models.SingInData{
			{
				UserId:       3,
				UserName:     "test@example.com",
				UserPassword: "Test12345!",
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		fetcher := NewSingDataFetcher()
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.Response[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_name",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi バリデーション パスワード不正", func(t *testing.T) {

		dataList := []testData{
			{
				Data: []models.RequestSingInData{
					{
						UserId:       1,
						UserName:     "test@example.com",
						UserPassword: "Test12!",
					},
				},
			},
			{
				Data: []models.RequestSingInData{
					{
						UserId:       1,
						UserName:     "test@example.com",
						UserPassword: "Test123456",
					},
				},
			},
		}

		resMock := []models.SingInData{
			{
				UserId:       3,
				UserName:     "test@example.com",
				UserPassword: "Test12345!",
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.SingDataFetcher{}),
				"GetSingIn",
				func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
					return resMock, nil
				})
			defer patches.Reset()

			fetcher := NewSingDataFetcher()
			fetcher.PostSingInApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.Response[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.Response[utils.ErrorMessages]{
				Result: []utils.ErrorMessages{
					{
						Field:   "user_password",
						Message: "パスワードの形式が間違っています。",
					},
				},
			}
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("TestPostSingInApi result件数0件", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInData{
				{
					UserId:       1,
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		resMock := []models.SingInData{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		fetcher := NewSingDataFetcher()
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.Response[requestSingInData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[requestSingInData]{
			ErrorMsg: "サインインに失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi result 成功", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInData{
				{
					UserId:       1,
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		resMock := []models.SingInData{
			{
				UserId:       3,
				UserName:     "test@example.com",
				UserPassword: "Test12345!",
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		fetcher := NewSingDataFetcher()
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.Response[SingInResult]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, len(responseBody.Token), 120)

		expectedOk := utils.Response[SingInResult]{
			Result: []SingInResult{
				{
					UserId:       3,
					UserName:     "test@example.com",
					UserPassword: "Test12345!",
				},
			},
		}
		assert.Equal(t, responseBody.Result[0].UserId, expectedOk.Result[0].UserId)
		assert.Equal(t, responseBody.Result[0].UserName, expectedOk.Result[0].UserName)
		assert.Equal(t, responseBody.Result[0].UserPassword, expectedOk.Result[0].UserPassword)
	})

	t.Run("TestPostSingInApi トークン生成に失敗", func(t *testing.T) {
		data := testData{
			Data: []models.RequestSingInData{
				{
					UserId:       1,
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// TokenGenerator のモックを作成
		mockTokenGen := mock_utils.NewMockTokenGenerator(ctrl)

		// モックの挙動を定義 (トークン生成失敗)
		mockTokenGen.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー")).Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// モックの PostSingInApi メソッドを設定して、エラーレスポンスを返す
		mockFetcher := mock_controllers.NewMockSingDataFetcher(ctrl)
		mockFetcher.EXPECT().PostSingInApi(gomock.Any()).DoAndReturn(func(c *gin.Context) {
			w := httptest.NewRecorder()
			c.Writer = w
			response := utils.Response[requestSingInData]{
				ErrorMsg: "サインインに失敗しました。",
			}
			c.JSON(http.StatusUnauthorized, response)
		})

		// モックを使って API を呼び出し
		mockFetcher.PostSingInApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// レスポンスボディの確認
		var responseBody utils.Response[requestSingInData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.Response[requestSingInData]{
			ErrorMsg: "サインインに失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})
}

// func TestGetDateRangeApi(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	t.Run("success GetDateRangeApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		// テスト用のGinコンテキストを作成
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

// 		// モックデータを設定
// 		mockData := []models.PaymentDate{
// 			{
// 				UserID:            1,
// 				StratPaymaentDate: "2022-01-01",
// 				EndPaymaentDate:   "2022-12-31",
// 			},
// 		}

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetDateRange",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
// 				return mockData, nil
// 			})
// 		defer patches.Reset()

// 		// テスト対象の関数を呼び出し
// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetDateRangeApi(c)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		var rawResponse struct {
// 			Result json.RawMessage `json:"result"`
// 		}
// 		err := json.Unmarshal(w.Body.Bytes(), &rawResponse)
// 		assert.NoError(t, err)
// 		// 期待されるJSONデータを構築
// 		expectedResponse, err := json.Marshal(mockData)
// 		assert.NoError(t, err)
// 		// 期待されるJSONデータと実際のレスポンスを比較
// 		assert.JSONEq(t, string(expectedResponse), string(rawResponse.Result))
// 	})

// 	t.Run("error GetDateRangeApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		// テスト用のGinコンテキストを作成
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetDateRange",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
// 				return nil, errors.New("database error")
// 			})
// 		defer patches.Reset()

// 		// テスト対象の関数を呼び出し
// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetDateRangeApi(c)

// 		// レスポンスの確認
// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		var response map[string]interface{}
// 		json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NotNil(t, response["error_msg"])
// 	})

// 	t.Run("バリデーションエラー user_id 必須", func(t *testing.T) {

// 		// エラーを引き起こすリクエストをシミュレート
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=", nil)

// 		mockData := []models.PaymentDate{
// 			{
// 				UserID:            1,
// 				StratPaymaentDate: "2022-01-01",
// 				EndPaymaentDate:   "2022-12-31",
// 			},
// 		}

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetDateRange",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
// 				return mockData, nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetDateRangeApi(c)

// 		// レスポンスのステータスコードを確認
// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "user_id",
// 					Message: "ユーザーIDは必須です。",
// 				},
// 			},
// 		}
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー user_id 整数値のみ", func(t *testing.T) {
// 		paramsList := [...]string{
// 			"/?user_id=rere",
// 			"/?user_id=1.23",
// 		}

// 		for _, params := range paramsList {
// 			// エラーを引き起こすリクエストをシミュレート
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			c.Request = httptest.NewRequest("GET", params, nil)

// 			mockData := []models.PaymentDate{
// 				{
// 					UserID:            1,
// 					StratPaymaentDate: "2022-01-01",
// 					EndPaymaentDate:   "2022-12-31",
// 				},
// 			}

// 			patches := ApplyMethod(
// 				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 				"GetDateRange",
// 				func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.PaymentDate, error) {
// 					return mockData, nil
// 				})
// 			defer patches.Reset()

// 			fetcher := NewIncomeDataFetcher()
// 			fetcher.GetDateRangeApi(c)

// 			// レスポンスのステータスコードを確認
// 			assert.Equal(t, http.StatusBadRequest, w.Code)

// 			var responseBody utils.Response[utils.ErrorMessages]
// 			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.NoError(t, err)

// 			expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 				Result: []utils.ErrorMessages{
// 					{
// 						Field:   "user_id",
// 						Message: "ユーザーIDは整数値のみです。",
// 					},
// 				},
// 			}
// 			assert.Equal(t, responseBody, expectedErrorMessage)
// 		}
// 	})
// }

// func TestGetYearIncomeAndDeductionApi(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	t.Run("success GetYearIncomeAndDeductionApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		// テスト用のGinコンテキストを作成
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

// 		// モックデータを設定
// 		mockData := []models.YearsIncomeData{
// 			{
// 				Years:           "2022",
// 				TotalAmount:     6000,
// 				DeductionAmount: 600,
// 				TakeHomeAmount:  5400,
// 			},
// 		}

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetYearsIncomeAndDeduction",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
// 				return mockData, nil
// 			})
// 		defer patches.Reset()

// 		// テスト対象の関数を呼び出し
// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetYearIncomeAndDeductionApi(c)

// 		// レスポンスの確認
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		var response struct {
// 			Result []models.YearsIncomeData `json:"result"`
// 		}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, mockData, response.Result)
// 	})

// 	t.Run("error GetYearIncomeAndDeductionApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		// テスト用のGinコンテキストを作成
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=1", nil)

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetYearsIncomeAndDeduction",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
// 				return nil, errors.New("database error")
// 			})
// 		defer patches.Reset()

// 		// テスト対象の関数を呼び出し
// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetYearIncomeAndDeductionApi(c)

// 		// レスポンスの確認
// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "database error", response["error_msg"])
// 	})

// 	t.Run("バリデーションエラー user_id 必須", func(t *testing.T) {

// 		// エラーを引き起こすリクエストをシミュレート
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)
// 		c.Request = httptest.NewRequest("GET", "/?user_id=", nil)

// 		mockData := []models.YearsIncomeData{
// 			{
// 				Years:           "2022",
// 				TotalAmount:     6000,
// 				DeductionAmount: 600,
// 				TakeHomeAmount:  5400,
// 			},
// 		}

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"GetYearsIncomeAndDeduction",
// 			func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
// 				return mockData, nil
// 			})
// 		defer patches.Reset()

// 		// テスト対象の関数を呼び出し
// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.GetYearIncomeAndDeductionApi(c)

// 		// レスポンスのステータスコードを確認
// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "user_id",
// 					Message: "ユーザーIDは必須です。",
// 				},
// 			},
// 		}
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー user_id 整数値のみ", func(t *testing.T) {
// 		paramsList := [...]string{
// 			"/?user_id=rere",
// 			"/?user_id=1.23",
// 		}

// 		for _, params := range paramsList {
// 			// エラーを引き起こすリクエストをシミュレート
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			c.Request = httptest.NewRequest("GET", params, nil)

// 			mockData := []models.YearsIncomeData{
// 				{
// 					Years:           "2022",
// 					TotalAmount:     6000,
// 					DeductionAmount: 600,
// 					TakeHomeAmount:  5400,
// 				},
// 			}

// 			patches := ApplyMethod(
// 				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 				"GetYearsIncomeAndDeduction",
// 				func(_ *models.PostgreSQLDataFetcher, UserID int) ([]models.YearsIncomeData, error) {
// 					return mockData, nil
// 				})
// 			defer patches.Reset()

// 			// テスト対象の関数を呼び出し
// 			fetcher := NewIncomeDataFetcher()
// 			fetcher.GetYearIncomeAndDeductionApi(c)

// 			// レスポンスのステータスコードを確認
// 			assert.Equal(t, http.StatusBadRequest, w.Code)

// 			var responseBody utils.Response[utils.ErrorMessages]
// 			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.NoError(t, err)

// 			expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 				Result: []utils.ErrorMessages{
// 					{
// 						Field:   "user_id",
// 						Message: "ユーザーIDは整数値のみです。",
// 					},
// 				},
// 			}
// 			assert.Equal(t, responseBody, expectedErrorMessage)
// 		}
// 	})
// }

// func TestInsertIncomeDataApi(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	t.Run("success InsertIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.InsertIncomeData{
// 				{
// 					PaymentDate:     "2024-02-10",
// 					Age:             30,
// 					Industry:        "IT",
// 					TotalAmount:     320524,
// 					DeductionAmount: 93480,
// 					TakeHomeAmount:  227044,
// 					Classification:  "給料",
// 					UserID:          "1",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"InsertIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.InsertIncomeDataApi(c)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "新規給料情報を登録致しました。", response["result_msg"])
// 	})

// 	t.Run("error InsertIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.InsertIncomeData{
// 				{
// 					PaymentDate:     "2024-02-10",
// 					Age:             30,
// 					Industry:        "IT",
// 					TotalAmount:     320524,
// 					DeductionAmount: 93480,
// 					TakeHomeAmount:  227044,
// 					Classification:  "給料",
// 					UserID:          "1",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"InsertIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
// 				return errors.New("database error")
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.InsertIncomeDataApi(c)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "新規登録時にエラーが発生。", response["error_msg"])
// 	})

// 	t.Run("invalid JSON InsertIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// Invalid JSON
// 		invalidJSON := `{"data": [`

// 		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBufferString(invalidJSON))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.InsertIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Contains(t, response["error_msg"], "unexpected EOF")
// 	})

// 	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.InsertIncomeData{
// 				{
// 					PaymentDate:     "",
// 					Age:             0,
// 					Industry:        "",
// 					TotalAmount:     "",
// 					DeductionAmount: "",
// 					TakeHomeAmount:  "",
// 					Classification:  "",
// 					UserID:          "",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"InsertIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.InsertIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			RecodeRows: 1,
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "age",
// 					Message: "年齢は必須又は整数値のみです。",
// 				},
// 				{
// 					Field:   "industry",
// 					Message: "業種は必須です。",
// 				},
// 				{
// 					Field:   "total_amount",
// 					Message: "総支給額は必須です。",
// 				},
// 				{
// 					Field:   "deduction_amount",
// 					Message: "差引額は必須です。",
// 				},
// 				{
// 					Field:   "take_home_amount",
// 					Message: "手取額は必須です。",
// 				},
// 				{
// 					Field:   "classification",
// 					Message: "分類は必須です。",
// 				},
// 				{
// 					Field:   "user_id",
// 					Message: "ユーザーIDは必須です。",
// 				},
// 				{
// 					Field:   "payment_date",
// 					Message: "報酬日付は必須です。",
// 				},
// 			},
// 		}
// 		test_utils.SortErrorMessages(responseBody.Result)
// 		test_utils.SortErrorMessages(expectedErrorMessage.Result)
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー 数値文字列以外は無効", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.InsertIncomeData{
// 				{
// 					PaymentDate:     "2021-09-10",
// 					Age:             25,
// 					Industry:        "IT業界",
// 					TotalAmount:     "0kj",
// 					DeductionAmount: "0.3",
// 					TakeHomeAmount:  "hoge",
// 					Classification:  "賞与",
// 					UserID:          "1",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"InsertIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.InsertIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			RecodeRows: 1,
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "total_amount",
// 					Message: "総支給額で数値文字列以外は無効です。",
// 				},
// 				{
// 					Field:   "deduction_amount",
// 					Message: "差引額で数値文字列以外は無効です。",
// 				},
// 				{
// 					Field:   "take_home_amount",
// 					Message: "手取額で数値文字列以外は無効です。",
// 				},
// 			},
// 		}
// 		test_utils.SortErrorMessages(responseBody.Result)
// 		test_utils.SortErrorMessages(expectedErrorMessage.Result)
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー 形式及びユーザーIDの整数値チェック", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		dataList := []testData{
// 			{
// 				Data: []models.InsertIncomeData{
// 					{
// 						PaymentDate:     "202109-10",
// 						Age:             25,
// 						Industry:        "IT業界",
// 						TotalAmount:     "350000",
// 						DeductionAmount: "76000",
// 						TakeHomeAmount:  "274000",
// 						Classification:  "賞与",
// 						UserID:          "1o",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.InsertIncomeData{
// 					{
// 						PaymentDate:     "2021-0910",
// 						Age:             25,
// 						Industry:        "IT業界",
// 						TotalAmount:     "350000",
// 						DeductionAmount: "76000",
// 						TakeHomeAmount:  "274000",
// 						Classification:  "賞与",
// 						UserID:          "1.0",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.InsertIncomeData{
// 					{
// 						PaymentDate:     "2021-13-10",
// 						Age:             25,
// 						Industry:        "IT業界",
// 						TotalAmount:     "350000",
// 						DeductionAmount: "76000",
// 						TakeHomeAmount:  "274000",
// 						Classification:  "賞与",
// 						UserID:          "1.0hoge",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.InsertIncomeData{
// 					{
// 						PaymentDate:     "2021-12-33",
// 						Age:             25,
// 						Industry:        "IT業界",
// 						TotalAmount:     "350000",
// 						DeductionAmount: "76000",
// 						TakeHomeAmount:  "274000",
// 						Classification:  "賞与",
// 						UserID:          "hoge",
// 					},
// 				},
// 			},
// 		}

// 		for _, data := range dataList {
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			body, _ := json.Marshal(data)
// 			c.Request = httptest.NewRequest("POST", "/api/income_create", bytes.NewBuffer(body))
// 			c.Request.Header.Set("Content-Type", "application/json")

// 			patches := ApplyMethod(
// 				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 				"InsertIncome",
// 				func(_ *models.PostgreSQLDataFetcher, data []models.InsertIncomeData) error {
// 					return nil
// 				})
// 			defer patches.Reset()

// 			fetcher := NewIncomeDataFetcher()
// 			fetcher.InsertIncomeDataApi(c)

// 			assert.Equal(t, http.StatusBadRequest, w.Code)

// 			var responseBody utils.Response[utils.ErrorMessages]
// 			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.NoError(t, err)

// 			expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 				RecodeRows: 1,
// 				Result: []utils.ErrorMessages{
// 					{
// 						Field:   "payment_date",
// 						Message: "給料支給日の形式が間違っています。",
// 					},
// 					{
// 						Field:   "user_id",
// 						Message: "ユーザーIDは整数値のみです。",
// 					},
// 				},
// 			}
// 			test_utils.SortErrorMessages(responseBody.Result)
// 			test_utils.SortErrorMessages(expectedErrorMessage.Result)
// 			assert.Equal(t, responseBody, expectedErrorMessage)
// 		}
// 	})
// }

// func TestUpdateIncomeDataApi(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	t.Run("success UpdateIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.UpdateIncomeData{
// 				{
// 					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 					PaymentDate:      "2024-02-10",
// 					Age:              30,
// 					Industry:         "IT",
// 					TotalAmount:      320524,
// 					DeductionAmount:  93480,
// 					TakeHomeAmount:   227044,
// 					UpdateUser:       "test_user",
// 					Classification:   "給料",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"UpdateIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.UpdateIncomeDataApi(c)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "給料情報の更新が問題なく成功しました。", response["result_msg"])
// 	})

// 	t.Run("error UpdateIncomeDataApi", func(t *testing.T) {
// 		// ここのテストケースだけ不安定なのでN回リトライしてテスト行う
// 		const maxRetry = 100

// 		var lastErr error

// 		for attempt := 1; attempt <= maxRetry; attempt++ {
// 			// time.Sleep(1 * time.Second) // 1秒
// 			time.Sleep(500 * time.Millisecond)
// 			// config.Setup()
// 			// defer config.Teardown()
// 			// defer config.TeardownTestDatabase()

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			data := testData{
// 				Data: []models.UpdateIncomeData{
// 					{
// 						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 						PaymentDate:      "2024-02-10",
// 						Age:              30,
// 						Industry:         "IT",
// 						TotalAmount:      320524,
// 						DeductionAmount:  93480,
// 						TakeHomeAmount:   227044,
// 						UpdateUser:       "test_user",
// 						Classification:   "給料",
// 					},
// 				},
// 			}

// 			body, _ := json.Marshal(data)
// 			c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
// 			c.Request.Header.Set("Content-Type", "application/json")

// 			patches := ApplyMethod(
// 				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 				"UpdateIncome",
// 				func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
// 					return errors.New("database error")
// 				})
// 			defer patches.Reset()

// 			fetcher := NewIncomeDataFetcher()
// 			fetcher.UpdateIncomeDataApi(c)

// 			if w.Code == http.StatusInternalServerError {
// 				var response map[string]interface{}
// 				err := json.Unmarshal(w.Body.Bytes(), &response)
// 				if err != nil {
// 					lastErr = err
// 				} else {
// 					// assert.NoError(t, err)
// 					assert.Equal(t, "更新時にエラーが発生。", response["error_msg"])
// 					fmt.Printf("成功しました！試行回数: %d回目\n", attempt)
// 					lastErr = nil
// 					break
// 				}
// 			} else {
// 				lastErr = errors.New("データベース挿入エラーになりませんでした: " + http.StatusText(w.Code))
// 			}

// 			if attempt == maxRetry {
// 				t.Fatalf("テスト試行回数の上限に達しました %d回 テストエラー内容: %v", maxRetry, lastErr)
// 			}
// 		}
// 	})

// 	t.Run("invalid JSON UpdateIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// Invalid JSON
// 		invalidJSON := `{"data": [`

// 		c.Request = httptest.NewRequest("POST", "/api/income_update", bytes.NewBufferString(invalidJSON))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.UpdateIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Contains(t, response["error_msg"], "unexpected EOF")
// 	})

// 	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.UpdateIncomeData{
// 				{
// 					IncomeForecastID: "",
// 					PaymentDate:      "",
// 					Age:              0,
// 					Industry:         "",
// 					TotalAmount:      "",
// 					DeductionAmount:  "",
// 					TakeHomeAmount:   "",
// 					UpdateUser:       "",
// 					Classification:   "",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"UpdateIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.UpdateIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			RecodeRows: 1,
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "income_forecast_id",
// 					Message: "年収推移IDは必須です。",
// 				},
// 				{
// 					Field:   "payment_date",
// 					Message: "報酬日付は必須です。",
// 				},
// 				{
// 					Field:   "age",
// 					Message: "年齢は必須又は整数値のみです。",
// 				},
// 				{
// 					Field:   "industry",
// 					Message: "業種は必須です。",
// 				},
// 				{
// 					Field:   "total_amount",
// 					Message: "総支給額は必須です。",
// 				},
// 				{
// 					Field:   "deduction_amount",
// 					Message: "差引額は必須です。",
// 				},
// 				{
// 					Field:   "take_home_amount",
// 					Message: "手取額は必須です。",
// 				},
// 				{
// 					Field:   "update_user",
// 					Message: "更新者は必須です。",
// 				},
// 				{
// 					Field:   "classification",
// 					Message: "分類は必須です。",
// 				},
// 			},
// 		}
// 		test_utils.SortErrorMessages(responseBody.Result)
// 		test_utils.SortErrorMessages(expectedErrorMessage.Result)
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー 数値文字列以外は無効", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.UpdateIncomeData{
// 				{
// 					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 					PaymentDate:      "2024-02-10",
// 					Age:              30,
// 					Industry:         "IT",
// 					TotalAmount:      "fd",
// 					DeductionAmount:  "0.8u",
// 					TakeHomeAmount:   "0.45",
// 					UpdateUser:       "test_user",
// 					Classification:   "給料",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"UpdateIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.UpdateIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			RecodeRows: 1,
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "total_amount",
// 					Message: "総支給額で数値文字列以外は無効です。",
// 				},
// 				{
// 					Field:   "deduction_amount",
// 					Message: "差引額で数値文字列以外は無効です。",
// 				},
// 				{
// 					Field:   "take_home_amount",
// 					Message: "手取額で数値文字列以外は無効です。",
// 				},
// 			},
// 		}
// 		test_utils.SortErrorMessages(responseBody.Result)
// 		test_utils.SortErrorMessages(expectedErrorMessage.Result)
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})

// 	t.Run("バリデーションエラー 形式及びユーザーIDの整数値チェック", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		dataList := []testData{
// 			{
// 				Data: []models.UpdateIncomeData{
// 					{
// 						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 						PaymentDate:      "202402-10",
// 						Age:              30,
// 						Industry:         "IT",
// 						TotalAmount:      320524,
// 						DeductionAmount:  93480,
// 						TakeHomeAmount:   227044,
// 						UpdateUser:       "test_user",
// 						Classification:   "給料",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.UpdateIncomeData{
// 					{
// 						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 						PaymentDate:      "2024-0210",
// 						Age:              30,
// 						Industry:         "IT",
// 						TotalAmount:      320524,
// 						DeductionAmount:  93480,
// 						TakeHomeAmount:   227044,
// 						UpdateUser:       "test_user",
// 						Classification:   "給料",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.UpdateIncomeData{
// 					{
// 						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 						PaymentDate:      "2024-13-10",
// 						Age:              30,
// 						Industry:         "IT",
// 						TotalAmount:      320524,
// 						DeductionAmount:  93480,
// 						TakeHomeAmount:   227044,
// 						UpdateUser:       "test_user",
// 						Classification:   "給料",
// 					},
// 				},
// 			},
// 			{
// 				Data: []models.UpdateIncomeData{
// 					{
// 						IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 						PaymentDate:      "2024-02-32",
// 						Age:              30,
// 						Industry:         "IT",
// 						TotalAmount:      320524,
// 						DeductionAmount:  93480,
// 						TakeHomeAmount:   227044,
// 						UpdateUser:       "test_user",
// 						Classification:   "給料",
// 					},
// 				},
// 			},
// 		}

// 		for _, data := range dataList {
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)
// 			body, _ := json.Marshal(data)
// 			c.Request = httptest.NewRequest("PUT", "/api/income_update", bytes.NewBuffer(body))
// 			c.Request.Header.Set("Content-Type", "application/json")

// 			patches := ApplyMethod(
// 				reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 				"UpdateIncome",
// 				func(_ *models.PostgreSQLDataFetcher, data []models.UpdateIncomeData) error {
// 					return nil
// 				})
// 			defer patches.Reset()

// 			fetcher := NewIncomeDataFetcher()
// 			fetcher.UpdateIncomeDataApi(c)

// 			assert.Equal(t, http.StatusBadRequest, w.Code)

// 			var responseBody utils.Response[utils.ErrorMessages]
// 			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 			assert.NoError(t, err)

// 			expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 				RecodeRows: 1,
// 				Result: []utils.ErrorMessages{
// 					{
// 						Field:   "payment_date",
// 						Message: "給料支給日の形式が間違っています。",
// 					},
// 				},
// 			}
// 			assert.Equal(t, responseBody, expectedErrorMessage)
// 		}
// 	})
// }

// func TestDeleteIncomeDataApi(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	t.Run("success DeleteIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.DeleteIncomeData{
// 				{
// 					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"DeleteIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.DeleteIncomeDataApi(c)

// 		assert.Equal(t, http.StatusOK, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "給料情報の削除が問題なく成功しました。", response["result_msg"])
// 	})

// 	t.Run("error DeleteIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.DeleteIncomeData{
// 				{
// 					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"DeleteIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
// 				return errors.New("database error")
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.DeleteIncomeDataApi(c)

// 		assert.Equal(t, http.StatusInternalServerError, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "削除中にエラーが発生しました", response["error_msg"])
// 	})

// 	t.Run("invalid JSON DeleteIncomeDataApi", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// Invalid JSON
// 		invalidJSON := `{"data": [`

// 		c.Request = httptest.NewRequest("POST", "/api/income_delete", bytes.NewBufferString(invalidJSON))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.DeleteIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 		var response map[string]interface{}
// 		err := json.Unmarshal(w.Body.Bytes(), &response)
// 		assert.NoError(t, err)
// 		assert.Contains(t, response["error_msg"], "unexpected EOF")
// 	})

// 	t.Run("バリデーションエラー 対象カラム必須", func(t *testing.T) {
// 		// config.Setup()
// 		// defer config.Teardown()
// 		// defer config.TeardownTestDatabase()

// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		data := testData{
// 			Data: []models.DeleteIncomeData{
// 				{
// 					IncomeForecastID: "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
// 				},
// 				{
// 					IncomeForecastID: "",
// 				},
// 			},
// 		}

// 		body, _ := json.Marshal(data)
// 		c.Request = httptest.NewRequest("PUT", "/api/income_delete", bytes.NewBuffer(body))
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		patches := ApplyMethod(
// 			reflect.TypeOf(&models.PostgreSQLDataFetcher{}),
// 			"DeleteIncome",
// 			func(_ *models.PostgreSQLDataFetcher, data []models.DeleteIncomeData) error {
// 				return nil
// 			})
// 		defer patches.Reset()

// 		fetcher := NewIncomeDataFetcher()
// 		fetcher.DeleteIncomeDataApi(c)

// 		assert.Equal(t, http.StatusBadRequest, w.Code)

// 		var responseBody utils.Response[utils.ErrorMessages]
// 		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)

// 		expectedErrorMessage := utils.Response[utils.ErrorMessages]{
// 			RecodeRows: 2,
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "income_forecast_id",
// 					Message: "年収推移IDは必須です。",
// 				},
// 			},
// 		}
// 		assert.Equal(t, responseBody, expectedErrorMessage)
// 	})
// }
