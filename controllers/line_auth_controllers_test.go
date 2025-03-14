package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	mock_templates "server/mock/templates"
	mock_utils "server/mock/utils"
	"server/models"
	"server/templates"
	"server/test_utils"
	"server/utils"
	"testing"

	// mock_config "server/mock/config"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLineSignInCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ResMock := []models.ExternalAuthData{
		{
			UserId:    1,
			UserEmail: "test@example.com",
		},
	}

	t.Run("LineSignInCallback バリデーション必須チェック", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=&user_name=",
			nil,
			map[string]string{
				"user_email": "",
				"user_name":  "",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "メールアドレスは必須です。",
				},
				{
					Field:   "user_name",
					Message: "ユーザー名は必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignInCallback バリデーションメールアドレス形式不正", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=hoge&user_name=test",
			nil,
			map[string]string{
				"user_email": "hoge",
				"user_name":  "test",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignInCallback DB取得エラー", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		resMock := []models.ExternalAuthData{}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "ユーザー情報取得に失敗しました。",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignInCallback トークン生成に失敗 1", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return ResMock, nil
			})
		defer patches.Reset()

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// レスポンスボディの確認
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "新規トークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignInCallback トークン生成に失敗 2", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return ResMock, nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("token", nil)

		mockUtilsFetcher.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		// レスポンスボディの確認
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "リフレッシュトークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignInCallback メールテンプレート生成エラー(サインイン)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return ResMock, nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("new_token", nil)

		mockUtilsFetcher.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("refresh_token", nil)

		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockEmailTemplateService.EXPECT().
			PostSignInTemplate(gomock.Any(), gomock.Any()).
			Return("", "", fmt.Errorf("テンプレート生成エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: mockEmailTemplateService,
		}
		lineManager.LineSignInCallback(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(サインイン): テンプレート生成エラー",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignInCallback メール送信エラー(サインイン)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return ResMock, nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("new_token", nil)

		mockUtilsFetcher.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("refresh_token", nil)

		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("メール送信エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メール送信エラー(サインイン): メール送信エラー",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignInCallback result 成功", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signin/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resMock := []models.ExternalAuthData{
			{
				UserId:    3,
				UserEmail: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("new_token", nil)

		mockUtilsFetcher.EXPECT().
			RefreshToken(gomock.Any(), gomock.Any()).
			Return("refresh_token", nil)

		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ResponseData[SignInResult]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// assert.Equal(t, len(responseBody.Token), 120)

		expectedOk := utils.ResponseData[SignInResult]{
			Result: SignInResult{
				UserId:       3,
				UserEmail:    "test@example.com",
				UserPassword: "",
			},
		}
		assert.Equal(t, responseBody.Result.UserId, expectedOk.Result.UserId)
		assert.Equal(t, responseBody.Result.UserEmail, expectedOk.Result.UserEmail)
		assert.Equal(t, responseBody.Result.UserPassword, expectedOk.Result.UserPassword)
	})
}

func TestLineSignUpCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("LineSignUpCallback バリデーション必須チェック", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=&user_name=",
			nil,
			map[string]string{
				"user_email": "",
				"user_name":  "",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignUpCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "メールアドレスは必須です。",
				},
				{
					Field:   "user_name",
					Message: "ユーザー名は必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignUpCallback バリデーションメールアドレス形式不正", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=hoge&user_name=test",
			nil,
			map[string]string{
				"user_email": "hoge",
				"user_name":  "test",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignUpCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineSignUpCallback DB取得エラー", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return fmt.Errorf("sql登録失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusConflict, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "既に登録されたメールアドレスです。",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignUpCallback メールテンプレート生成エラー(登録)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockEmailTemplateService.EXPECT().
			PostSignUpTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("件名", "本文", fmt.Errorf("メールテンプレートエラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: mockEmailTemplateService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(登録): メールテンプレートエラー",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignUpCallback メール送信エラー(登録)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockEmailTemplateService.EXPECT().
			PostSignUpTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("件名", "本文", nil)

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("メール送信エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: mockEmailTemplateService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メール送信エラー(登録): メール送信エラー",
		}
		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineSignUpCallback result 成功", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/signup/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedResponse := utils.ResponseData[string]{
			Result: "line外部認証の登録が成功しました。",
		}
		assert.Equal(t, responseBody.Result, expectedResponse.Result)
	})
}

func TestLineDeleteCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("LineDeleteCallback バリデーション必須チェック", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=&user_name=",
			nil,
			map[string]string{
				"user_email": "",
				"user_name":  "",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineDeleteCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "メールアドレスは必須です。",
				},
				{
					Field:   "user_name",
					Message: "ユーザー名は必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineDeleteCallback バリデーションメールアドレス形式不正", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=hoge&user_name=test",
			nil,
			map[string]string{
				"user_email": "hoge",
				"user_name":  "test",
			},
		)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineDeleteCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var responseBody utils.ErrorValidationResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorValidationResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "user_email",
					Message: "正しいメールアドレス形式である必要があります。",
				},
			},
		}
		test_utils.SortErrorMessages(responseBody.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, responseBody, expectedErrorMessage)
	})

	t.Run("LineDeleteCallback DB取得エラー", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		resMock := []models.ExternalAuthData{}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "sql取得失敗",
		}

		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineDeleteCallback DB削除エラー", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		resMock := []models.ExternalAuthData{
			{
				UserId:    1,
				UserEmail: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, userId int, data models.RequestSignInDeleteData) error {
				return fmt.Errorf("DB削除エラー")
			})
		defer patches1.Reset()

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "DB削除エラー",
		}

		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineDeleteCallback メールテンプレート生成エラー(削除)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)

		resMock := []models.ExternalAuthData{
			{
				UserId:    1,
				UserEmail: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, userId int, data models.RequestSignInDeleteData) error {
				return nil
			})
		defer patches1.Reset()

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockEmailTemplateService.EXPECT().
			DeleteSignInTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("件名", "本文", fmt.Errorf("メールテンプレートエラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: mockEmailTemplateService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メールテンプレート生成エラー(削除): メールテンプレートエラー",
		}

		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineDeleteCallback メール送信エラー(削除)", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)

		resMock := []models.ExternalAuthData{
			{
				UserId:    1,
				UserEmail: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, userId int, data models.RequestSignInDeleteData) error {
				return nil
			})
		defer patches1.Reset()

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockEmailTemplateService.EXPECT().
			DeleteSignInTemplate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("件名", "本文", nil)

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("メール送信エラー"))

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: mockEmailTemplateService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedErrorMessage := utils.ErrorMessageResponse{
			Result: "メール送信エラー(削除): メール送信エラー",
		}

		assert.Equal(t, responseBody.Result, expectedErrorMessage.Result)
	})

	t.Run("LineDeleteCallback result 成功", func(t *testing.T) {

		w, c := test_utils.CreateTestRequest(
			"GET", "/api/line/delete/callback?user_email=test@example.com&user_name=test",
			nil,
			map[string]string{
				"user_email": "test@example.com",
				"user_name":  "test",
			},
		)

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		resMock := []models.ExternalAuthData{
			{
				UserId:    1,
				UserEmail: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, userId int, data models.RequestSignInDeleteData) error {
				return nil
			})
		defer patches1.Reset()

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			DateTimeStr(gomock.Any(), gomock.Any()).
			Return("2024年12月2日")

		mockUtilsFetcher.EXPECT().
			SendMail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		lineManager := LineManager{
			UtilsFetcher:         mockUtilsFetcher,
			EmailTemplateService: templates.NewEmailTemplateManager(),
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody utils.ErrorMessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.Nil(t, err)

		expectedResponse := utils.ResponseData[string]{
			Result: "line外部認証の削除が成功しました。",
		}

		assert.Equal(t, responseBody.Result, expectedResponse.Result)
	})
}
