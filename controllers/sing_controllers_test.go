package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"

	// "server/config"

	"server/common"
	"server/models"
	"server/test_utils"
	"server/utils"
	"testing"

	mock_utils "server/mock/utils"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostSingInApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("TestPostSingInApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
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

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
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

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
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
						UserName:     "test@example.com",
						UserPassword: "Test12!",
					},
				},
			},
			{
				Data: []models.RequestSingInData{
					{
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

			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.PostSingInApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
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

	// 処理削除したのでスキップ
	// t.Run("TestPostSingInApi result件数0件", func(t *testing.T) {

	// 	data := testData{
	// 		Data: []models.RequestSingInData{
	// 			{
	// 				UserName:     "test@example.com",
	// 				UserPassword: "Test123456!!",
	// 			},
	// 		},
	// 	}

	// 	resMock := []models.SingInData{}

	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	body, _ := json.Marshal(data)
	// 	c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
	// 	c.Request.Header.Set("Content-Type", "application/json")

	// 	patches := ApplyMethod(
	// 		reflect.TypeOf(&models.SingDataFetcher{}),
	// 		"GetSingIn",
	// 		func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
	// 			return resMock, nil
	// 		})
	// 	defer patches.Reset()

	// 	fetcher := apiSingDataFetcher{
	// 		UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
	// 		CommonFetcher: common.NewCommonFetcher(),
	// 	}
	// 	fetcher.PostSingInApi(c)

	// 	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 	var responseBody utils.ResponseWithSlice[requestSingInData]
	// 	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	// 	assert.NoError(t, err)

	// 	expectedErrorMessage := utils.ResponseWithSlice[requestSingInData]{
	// 		ErrorMsg: "サインインに失敗しました。",
	// 	}
	// 	assert.Equal(t, responseBody, expectedErrorMessage)
	// })

	t.Run("TestPostSingInApi result 成功", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInData{
				{
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

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ResponseWithSlice[SingInResult]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// assert.Equal(t, len(responseBody.Token), 120)

		expectedOk := utils.ResponseWithSlice[SingInResult]{
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

	t.Run("TestPostSingInApi sql取得失敗しエラー発生", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInData{
				{
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
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSlice[SingInResult]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedError := utils.ResponseWithSlice[SingInResult]{
			ErrorMsg: "sql取得失敗",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedError.ErrorMsg)
	})

	t.Run("TestPostSingInApi トークン生成に失敗 1", func(t *testing.T) {
		data := testData{
			Data: []models.RequestSingInData{
				{
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

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// GetSingIn のモック化
		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches1.Reset()

		// モックを使って API を呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[requestSingInData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[requestSingInData]{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi トークン生成に失敗 2", func(t *testing.T) {
		data := testData{
			Data: []models.RequestSingInData{
				{
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

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("token", nil)

		mockUtilsFetcher.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// GetSingIn のモック化
		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"GetSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInData) ([]models.SingInData, error) {
				return resMock, nil
			})
		defer patches1.Reset()

		// モックを使って API を呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingInApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[requestSingInData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[requestSingInData]{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})
}

func TestGetRefreshTokenApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("TestGetRefreshTokenApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/api/refresh_token", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// UtilsFetcher のモックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// レスポンスボディの確認
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

	t.Run("TestPostSingInApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// UtilsFetcher のモックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

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
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// UtilsFetcher のモックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

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
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestPostSingInApi バリデーション 数値文字列以外", func(t *testing.T) {

		paramsList := [2]string{
			"/api/refresh_token?user_id=test",
			"/api/refresh_token?user_id=1.23",
		}

		for _, params := range paramsList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			// リクエストにパラメータを設定
			c.Request = httptest.NewRequest("GET", params, nil)
			c.Request.Header.Set("Content-Type", "application/json")

			// UtilsFetcher のモックを使ってAPIを呼び出し
			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.GetRefreshTokenApi(c)

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
			test_utils.SortErrorMessages(responseBody.Result)
			test_utils.SortErrorMessages(expectedErrorMessage.Result)
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("TestGetRefreshTokenApi リフレッシュトークンなし", func(t *testing.T) {

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=1", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  "AuthToken",
			Value: "dummy_token",
		})

		// モックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[RequestRefreshToken]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンがありません。再ログインしてください。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestGetRefreshTokenApi リフレッシュトークンが無効", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			ParseWithClaims(gomock.Any()).
			Return(nil, fmt.Errorf("トークン無効"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=1", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  "AuthToken",
			Value: "dummy_token",
		})

		// 2つ目のCookieを設定
		c.Request.AddCookie(&http.Cookie{
			Name:  "RefreshAuthToken",
			Value: "dummy_token",
		})

		// モックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[RequestRefreshToken]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "リフレッシュトークンが無効です。再ログインしてください。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestGetRefreshTokenApi 無効なリフレッシュトークン", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの MapClaims を定義
		mockClaims := jwt.MapClaims{
			"UserId": float64(1),
			"exp":    float64(time.Now().Add(time.Hour).Unix()),
		}

		mockToken := &jwt.Token{
			Claims: jwt.MapClaims{
				"UserId": float64(1),
				"exp":    float64(time.Now().Add(time.Hour).Unix()),
			},
			Valid: false,
		}

		// ParseWithClaimsのモックを定義
		mockUtilsFetcher.EXPECT().
			ParseWithClaims(gomock.Any()).
			Return(mockToken, nil)

		mockUtilsFetcher.EXPECT().
			MapClaims(gomock.Any()).
			Return(mockClaims, false)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=1", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  "AuthToken",
			Value: "dummy_token",
		})

		// 2つ目のCookieを設定
		c.Request.AddCookie(&http.Cookie{
			Name:  "RefreshAuthToken",
			Value: "dummy_token",
		})

		// モックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[RequestRefreshToken]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "無効なリフレッシュトークン。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestGetRefreshTokenApi 新しいアクセストークンの生成に失敗", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの MapClaims を定義
		mockClaims := jwt.MapClaims{
			"UserId": float64(1),
			"exp":    float64(time.Now().Add(time.Hour).Unix()),
		}

		mockToken := &jwt.Token{
			Claims: jwt.MapClaims{
				"UserId": float64(1),
				"exp":    float64(time.Now().Add(time.Hour).Unix()),
			},
			Valid: true,
		}

		// ParseWithClaimsのモックを定義
		mockUtilsFetcher.EXPECT().
			ParseWithClaims(gomock.Any()).
			Return(mockToken, nil)

		mockUtilsFetcher.EXPECT().
			MapClaims(gomock.Any()).
			Return(mockClaims, true)

		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=1", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  "AuthToken",
			Value: "dummy_token",
		})

		// 2つ目のCookieを設定
		c.Request.AddCookie(&http.Cookie{
			Name:  "RefreshAuthToken",
			Value: "dummy_token",
		})

		// モックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSlice[RequestRefreshToken]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[RequestRefreshToken]{
			ErrorMsg: "新しいアクセストークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("TestGetRefreshTokenApi 成功", func(t *testing.T) {

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの MapClaims を定義
		mockClaims := jwt.MapClaims{
			"UserId": float64(1),
			"exp":    float64(time.Now().Add(time.Hour).Unix()),
		}

		mockToken := &jwt.Token{
			Claims: jwt.MapClaims{
				"UserId": float64(1),
				"exp":    float64(time.Now().Add(time.Hour).Unix()),
			},
			Valid: true,
		}

		// ParseWithClaimsのモックを定義
		mockUtilsFetcher.EXPECT().
			ParseWithClaims(gomock.Any()).
			Return(mockToken, nil)

		mockUtilsFetcher.EXPECT().
			MapClaims(gomock.Any()).
			Return(mockClaims, true)

		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("new_token", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// リクエストにパラメータを設定
		c.Request = httptest.NewRequest("GET", "/api/refresh_token?user_id=1", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  "AuthToken",
			Value: "dummy_token",
		})

		// 2つ目のCookieを設定
		c.Request.AddCookie(&http.Cookie{
			Name:  "RefreshAuthToken",
			Value: "dummy_token",
		})

		// モックを使ってAPIを呼び出し
		fetcher := apiSingDataFetcher{
			UtilsFetcher:  mockUtilsFetcher,
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.GetRefreshTokenApi(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOK := utils.ResponseWithSingle[string]{
			Result: "新しいアクセストークンが発行されました。",
		}
		assert.Equal(t, responseBody.Result, expectedOK.Result)
	})
}

func TestPostSingUpApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("PostSingUpApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingUpApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("PostSingUpApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingUpData{
				{
					NickName:     "",
					UserName:     "",
					UserPassword: "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PostSingUp",
			func(_ *models.SingDataFetcher, data models.RequestSingUpData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingUpApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "nick_name",
					Message: "ニックネームは必須です。",
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

	t.Run("PostSingUpApi バリデーション メールアドレス不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingUpData{
				{
					NickName:     "test",
					UserName:     "test",
					UserPassword: "Test12345!",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PostSingUp",
			func(_ *models.SingDataFetcher, data models.RequestSingUpData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingUpApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_name",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("PostSingUpApi バリデーション パスワード不正", func(t *testing.T) {

		dataList := []testData{
			{
				Data: []models.RequestSingUpData{
					{
						NickName:     "test",
						UserName:     "test@example.com",
						UserPassword: "Test12!",
					},
				},
			},
			{
				Data: []models.RequestSingUpData{
					{
						NickName:     "test",
						UserName:     "test@example.com",
						UserPassword: "Test123456",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.SingDataFetcher{}),
				"PostSingUp",
				func(_ *models.SingDataFetcher, data models.RequestSingUpData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.PostSingUpApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
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

	t.Run("PostSingUpApi sql取得で失敗しサインアップ失敗になる", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingUpData{
				{
					NickName:     "test",
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PostSingUp",
			func(_ *models.SingDataFetcher, data models.RequestSingUpData) error {
				return fmt.Errorf("sql登録失敗")
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingUpApi(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSlice[requestSingInData]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[requestSingInData]{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("PostSingUpApi result 成功", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingUpData{
				{
					NickName:     "test",
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singup", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PostSingUp",
			func(_ *models.SingDataFetcher, data models.RequestSingUpData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PostSingUpApi(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSingle[string]{
			Result: "サインアップに成功",
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}

func TestPutSingInEditApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("PutSingInEditApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PutSingInEditApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("PutSingInEditApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingInEditData{
				{
					UserId:       "",
					UserName:     "",
					UserPassword: "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PutSingInEdit",
			func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PutSingInEditApi(c)

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
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("PutSingInEditApi バリデーション 数値文字列以外", func(t *testing.T) {

		dataList := []testData{
			{
				Data: []models.RequestSingInEditData{
					{
						UserId:       "test",
						UserName:     "",
						UserPassword: "",
					},
				},
			},
			{
				Data: []models.RequestSingInEditData{
					{
						UserId:       "1.25",
						UserName:     "",
						UserPassword: "",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.SingDataFetcher{}),
				"PutSingInEdit",
				func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.PutSingInEditApi(c)

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
			test_utils.SortErrorMessages(responseBody.Result)
			test_utils.SortErrorMessages(expectedErrorMessage.Result)
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("PutSingInEditApi バリデーション メールアドレス不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingInEditData{
				{
					UserId:       "1",
					UserName:     "test",
					UserPassword: "Test12345!",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PutSingInEdit",
			func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PutSingInEditApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_name",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("PutSingInEditApi バリデーション パスワード不正", func(t *testing.T) {

		dataList := []testData{
			{
				Data: []models.RequestSingInEditData{
					{
						UserId:       "1",
						UserName:     "test@example.com",
						UserPassword: "Test12!",
					},
				},
			},
			{
				Data: []models.RequestSingInEditData{
					{
						UserId:       "2",
						UserName:     "test@example.com",
						UserPassword: "Test123456",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.SingDataFetcher{}),
				"PutSingInEdit",
				func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.PutSingInEditApi(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var responseBody utils.ResponseWithSlice[utils.ErrorMessages]
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedErrorMessage := utils.ResponseWithSlice[utils.ErrorMessages]{
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

	t.Run("PutSingInEditApi sql取得で失敗しサインイン情報編集に失敗になる", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInEditData{
				{
					UserId:       "1",
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PutSingInEdit",
			func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
				return fmt.Errorf("sql更新失敗")
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PutSingInEditApi(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSlice[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインイン情報編集に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("PutSingInEditApi result 成功", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInEditData{
				{
					UserId:       "1",
					UserName:     "test@example.com",
					UserPassword: "Test123456!!",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_edit", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"PutSingInEdit",
			func(_ *models.SingDataFetcher, data models.RequestSingInEditData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.PutSingInEditApi(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSingle[string]{
			Result: "サインイン編集に成功",
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}

func TestDeleteSingInApi(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("DeleteSingInApi JSON不正", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		invalidJSON := `{"data": [`

		c.Request = httptest.NewRequest("POST", "/api/singin_delete", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteSingInApi(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error_msg"], "unexpected EOF")
	})

	t.Run("DeleteSingInApi バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := testData{
			Data: []models.RequestSingInDeleteData{
				{
					UserId: "",
				},
			},
		}

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"DeleteSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInDeleteData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteSingInApi(c)

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
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("DeleteSingInApi バリデーション 数値文字列以外", func(t *testing.T) {

		dataList := []testData{
			{
				Data: []models.RequestSingInDeleteData{
					{
						UserId: "test",
					},
				},
			},
			{
				Data: []models.RequestSingInDeleteData{
					{
						UserId: "1.25",
					},
				},
			},
		}

		for _, data := range dataList {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(data)
			c.Request = httptest.NewRequest("POST", "/api/singin_delete", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			patches := ApplyMethod(
				reflect.TypeOf(&models.SingDataFetcher{}),
				"DeleteSingIn",
				func(_ *models.SingDataFetcher, data models.RequestSingInDeleteData) error {
					return nil
				})
			defer patches.Reset()

			fetcher := apiSingDataFetcher{
				UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
				CommonFetcher: common.NewCommonFetcher(),
			}
			fetcher.DeleteSingInApi(c)

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
			test_utils.SortErrorMessages(responseBody.Result)
			test_utils.SortErrorMessages(expectedErrorMessage.Result)
			assert.Equal(t, responseBody, expectedErrorMessage)
		}
	})

	t.Run("DeleteSingInApi sql取得で失敗しサインインの削除失敗になる", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInDeleteData{
				{
					UserId: "1",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"DeleteSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInDeleteData) error {
				return fmt.Errorf("sql削除失敗")
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteSingInApi(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSlice[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ResponseWithSlice[string]{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("DeleteSingInApi result 成功", func(t *testing.T) {

		data := testData{
			Data: []models.RequestSingInDeleteData{
				{
					UserId: "1",
				},
			},
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(data)
		c.Request = httptest.NewRequest("POST", "/api/singin_delete", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		patches := ApplyMethod(
			reflect.TypeOf(&models.SingDataFetcher{}),
			"DeleteSingIn",
			func(_ *models.SingDataFetcher, data models.RequestSingInDeleteData) error {
				return nil
			})
		defer patches.Reset()

		fetcher := apiSingDataFetcher{
			UtilsFetcher:  utils.NewUtilsFetcher(utils.JwtSecret),
			CommonFetcher: common.NewCommonFetcher(),
		}
		fetcher.DeleteSingInApi(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSingle[string]{
			Result: "サインイン削除に成功",
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}
