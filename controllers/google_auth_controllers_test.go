package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"server/config"
	controllers_common "server/controllers/common"
	"server/models"
	"server/templates"
	"server/utils"
	"testing"

	mock_config "server/mock/config"
	mock_controllers_common "server/mock/controllers/common"
	mock_templates "server/mock/templates"
	mock_utils "server/mock/utils"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {

	gin.SetMode(gin.TestMode)
	mockRedirectURL := "https://accounts.google.com/o/oauth2/auth?client_id=mock_client_id"

	t.Run("GoogleSignIn リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().GoogleAuthURL(gomock.Any()).Return(mockRedirectURL)

		googleManager := GoogleManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			GoogleConfig:         mockGoogleService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/signin", nil)
		googleManager.GoogleSignIn(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})

	t.Run("GoogleSignUp リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().GoogleAuthURL(gomock.Any()).Return(mockRedirectURL)

		googleManager := GoogleManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			GoogleConfig:         mockGoogleService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/signup", nil)
		googleManager.GoogleSignUp(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})

	t.Run("GoogleDelete リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().GoogleAuthURL(gomock.Any()).Return(mockRedirectURL)

		googleManager := GoogleManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			GoogleConfig:         mockGoogleService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/google/delete", nil)
		googleManager.GoogleDelete(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})
}

func TestGoogleSignInCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ResMock := []models.ExternalAuthData{
		{
			UserId:   1,
			UserName: "test@example.com",
		},
	}

	t.Run("GoogleSignInCallback GoogleAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		response := utils.ErrorResponse{
			ErrorMsg: "テストエラー",
		}

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(500, controllers_common.GoogleUserInfo{}, response)

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "テストエラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			Email: "test@example.com",
		}

		response := utils.ErrorResponse{}

		resMock := []models.ExternalAuthData{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "sql取得失敗",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback トークン生成に失敗 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			Email: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return ResMock, nil
			})
		defer patches.Reset()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "新規トークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback トークン生成に失敗 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			Email: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "リフレッシュトークンの生成に失敗しました。",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback メールテンプレート生成エラー(サインイン)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			Email: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(サインイン): テンプレート生成エラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback メール送信エラー(サインイン)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			Email: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(サインイン): メール送信エラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("GoogleSignInCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockToken := &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			TokenType:    "Bearer",
		}

		userInfo := controllers_common.GoogleUserInfo{
			ID:            "1234",
			UserId:        1,
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "test",
			GivenName:     "test",
			FamilyName:    "test",
			Picture:       "111",
			Locale:        "222",
			Token:         mockToken,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
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
			Return(nil)

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody controllers_common.GoogleUserInfo
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := userInfo
		assert.Equal(t, responseBody, expectedOk)
	})
}
