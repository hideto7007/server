package controllers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"server/config"
	controllers_common "server/controllers/common"
	"server/models"
	"server/templates"
	"server/test_utils"
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
			UserId:    1,
			UserEmail: "test@example.com",
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
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "外部認証情報取得に失敗しました。")
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
			UserEmail: "test@example.com",
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
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
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
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "ユーザー情報取得に失敗しました。")
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
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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
			Return("", fmt.Errorf("トークン生成エラー"))

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "新規トークンの生成に失敗しました。")
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
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "リフレッシュトークンの生成に失敗しました。")
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
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
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
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
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
			UserEmail:     "test@example.com",
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
			Return(nil)

		mockControllersCommonService.EXPECT().
			RedirectSignIn(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("http://localhost:8080/test?user_id=1&user_email=test@example.com")

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignInCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		userId, userEmail, err := test_utils.RedirectSuccess(location)
		assert.Nil(t, err)
		assert.Equal(t, userId, userInfo.UserId)
		assert.Equal(t, userEmail, userInfo.UserEmail)
	})
}

func TestGoogleSignUpCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("GoogleSignUpCallback GoogleAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signup/callback?code=test-code", nil)
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
		googleManager.GoogleSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "外部認証情報取得に失敗しました。")
	})

	t.Run("GoogleSignUpCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return fmt.Errorf("sql登録失敗")
			})
		defer patches.Reset()

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "既に登録されたメールアドレスです。")
	})

	t.Run("GoogleSignUpCallback メールテンプレート生成エラー(登録)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
	})

	t.Run("GoogleSignUpCallback メール送信エラー(登録)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
	})

	t.Run("GoogleSignUpCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/signup/callback?code=test-code", nil)
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
			UserEmail:     "test@example.com",
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

		mockControllersCommonService.EXPECT().
			RedirectSignIn(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("http://localhost:8080/test")

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	})
}

func TestGoogleDeleteCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	MockToken := &oauth2.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
	}

	t.Run("GoogleDeleteCallback GoogleAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
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
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "外部認証情報取得に失敗しました。")
	})

	t.Run("GoogleDeleteCallback 無効なトークンのため削除できません", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
			Token:     MockToken,
		}
		response := utils.ErrorResponse{}

		resp := &http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
		}

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resp, fmt.Errorf("取得エラー"))

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "無効なトークンのため削除できません。")
	})

	t.Run("GoogleDeleteCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
			Token:     MockToken,
		}

		response := utils.ErrorResponse{}

		resq := &http.Response{
			StatusCode: http.StatusOK,
		}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resq, nil)

		resMock := []models.ExternalAuthData{}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserEmail string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
	})

	t.Run("GoogleDeleteCallback DB削除エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
			Token:     MockToken,
		}

		response := utils.ErrorResponse{}

		resq := &http.Response{
			StatusCode: http.StatusOK,
		}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resq, nil)

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

		googleManager := GoogleManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "削除中にエラーが発生しました。")
	})

	t.Run("GoogleDeleteCallback メールテンプレート生成エラー(削除)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
			Token:     MockToken,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		resq := &http.Response{
			StatusCode: http.StatusOK,
		}

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resq, nil)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
	})

	t.Run("GoogleDeleteCallback メール送信エラー(削除)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			UserEmail: "test@example.com",
			Name:      "test",
			Token:     MockToken,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		resq := &http.Response{
			StatusCode: http.StatusOK,
		}

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resq, nil)

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

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		msg, err := test_utils.QueryUnescape(location)
		assert.Nil(t, err)

		assert.Equal(t, msg, "予期せぬエラーが発生しました。")
	})

	t.Run("GoogleDeleteCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/google/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := controllers_common.GoogleUserInfo{
			ID:            "1234",
			UserId:        1,
			UserEmail:     "test@example.com",
			VerifiedEmail: true,
			Name:          "test",
			GivenName:     "test",
			FamilyName:    "test",
			Picture:       "111",
			Locale:        "222",
			Token:         MockToken,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		resq := &http.Response{
			StatusCode: http.StatusOK,
		}

		mockControllersCommonService.EXPECT().
			GoogleAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockControllersCommonService.EXPECT().
			GetRevoke(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(resq, nil)

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

		mockControllersCommonService.EXPECT().
			RedirectSignIn(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("http://localhost:8080/test")

		googleManager := GoogleManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			GoogleConfig:             config.NewGoogleManager(),
			ControllersCommonService: mockControllersCommonService,
		}
		googleManager.GoogleDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	})
}
