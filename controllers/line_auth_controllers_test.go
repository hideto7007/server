package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"server/config"
	"server/controllers/common"
	"server/models"
	"server/templates"
	"server/utils"
	"testing"

	mock_common "server/mock/common"
	mock_config "server/mock/config"
	mock_controllers_common "server/mock/controllers/common"
	mock_templates "server/mock/templates"
	mock_utils "server/mock/utils"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLineRedirect(t *testing.T) {

	gin.SetMode(gin.TestMode)
	mockRedirectURL := "https://access.line.me/oauth2/v2.1/authorize?client_id=mock_client_id"
	state := "state"

	t.Run("LineSignIn リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().LineAuthURL(gomock.Any()).Return(mockRedirectURL, state)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			LineConfig:           mockLineService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/line/signin", nil)
		lineManager.LineSignIn(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})

	t.Run("LineSignUp リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().LineAuthURL(gomock.Any()).Return(mockRedirectURL, state)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			LineConfig:           mockLineService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/line/signup", nil)
		lineManager.LineSignUp(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})

	t.Run("LineDelete リダイレクト成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().LineAuthURL(gomock.Any()).Return(mockRedirectURL, state)

		lineManager := LineManager{
			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService: templates.NewEmailTemplateManager(),
			LineConfig:           mockLineService,
		}
		// Ginのリクエストを設定
		c.Request = httptest.NewRequest(http.MethodGet, "/auth/line/delete", nil)
		lineManager.LineDelete(c)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)        // HTTPステータス確認
		assert.Equal(t, mockRedirectURL, w.Header().Get("Location")) // リダイレクトURL確認
	})
}

func TestLineSignInCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	ResMock := []models.ExternalAuthData{
		{
			UserId:   1,
			UserName: "test@example.com",
		},
	}

	t.Run("LineSignInCallback LineAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		response := utils.ErrorResponse{
			ErrorMsg: "テストエラー",
		}

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(500, &config.LineUserInfo{}, response)

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		resMock := []models.ExternalAuthData{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback トークン生成に失敗 1", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		// モックの挙動を定義
		mockUtilsFetcher.EXPECT().
			NewToken(gomock.Any(), gomock.Any()).
			Return("", fmt.Errorf("トークン生成エラー"))

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback トークン生成に失敗 2", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback メールテンプレート生成エラー(サインイン)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback メール送信エラー(サインイン)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

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

	t.Run("LineSignInCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signin/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResp := &config.LineTokenResponse{
			AccessToken: "token",
			IdToken:     "id_token",
		}

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   mockResp,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignInCallback(c)

		// ステータスコードの確認
		// assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

		var responseBody utils.ResponseWithSlice[common.GoogleUserInfo]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSlice[common.GoogleUserInfo]{
			Result: []common.GoogleUserInfo{
				{
					UserId:   userInfo.UserId,
					UserName: userInfo.UserName,
				},
			},
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}

func TestLineSignUpCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("LineSignUpCallback LineAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		response := utils.ErrorResponse{
			ErrorMsg: "テストエラー",
		}

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(500, &config.LineUserInfo{}, response)

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignUpCallback(c)

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

	t.Run("LineSignUpCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"PostSignUp",
			func(_ *models.SignDataFetcher, data models.RequestSignUpData) error {
				return fmt.Errorf("sql登録失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusConflict, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "既に登録されたメールアドレスです。",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineSignUpCallback メールテンプレート生成エラー(登録)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(登録): メールテンプレートエラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineSignUpCallback メール送信エラー(登録)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			UserName: "test@example.com",
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(登録): メール送信エラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineSignUpCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/signup/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResp := &config.LineTokenResponse{
			AccessToken: "token",
			IdToken:     "id_token",
		}

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   mockResp,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
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

		lineManager := LineManager{
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineSignUpCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSingle[string]{
			Result: "Line外部認証の登録成功しました。",
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}

func TestLineDeleteCallback(t *testing.T) {

	gin.SetMode(gin.TestMode)

	MockResp := &config.LineTokenResponse{
		AccessToken: "token",
		IdToken:     "id_token",
	}

	t.Run("LineDeleteCallback LineAuthCommonエラー返却", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockHttpService のモックを作成
		mockHttpService := mock_common.NewMockHttpService(ctrl)

		response := utils.ErrorResponse{
			ErrorMsg: "テストエラー",
		}

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(500, &config.LineUserInfo{}, response)

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               config.NewLineManager(mockHttpService),
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

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

	t.Run("LineDeleteCallback 無効なトークンのため削除できません", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(fmt.Errorf("取得エラー"))

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "取得エラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineDeleteCallback DB取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(nil)

		resMock := []models.ExternalAuthData{}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, fmt.Errorf("sql取得失敗")
			})
		defer patches.Reset()

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

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

	t.Run("LineDeleteCallback DB削除エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(nil)

		resMock := []models.ExternalAuthData{
			{
				UserId:   1,
				UserName: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, data models.RequestSignInDeleteData) error {
				return fmt.Errorf("DB削除エラー")
			})
		defer patches1.Reset()

		lineManager := LineManager{
			UtilsFetcher:             utils.NewUtilsFetcher(utils.JwtSecret),
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "サインインの削除に失敗しました。",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineDeleteCallback メールテンプレート生成エラー(削除)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(nil)

		resMock := []models.ExternalAuthData{
			{
				UserId:   1,
				UserName: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, data models.RequestSignInDeleteData) error {
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
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メールテンプレート生成エラー(削除): メールテンプレートエラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineDeleteCallback メール送信エラー(削除)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// EmailTemplateService のモックを作成
		mockEmailTemplateService := mock_templates.NewMockEmailTemplateService(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(nil)

		resMock := []models.ExternalAuthData{
			{
				UserId:   1,
				UserName: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, data models.RequestSignInDeleteData) error {
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
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     mockEmailTemplateService,
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メール送信エラー(削除): メール送信エラー",
		}
		assert.Equal(t, responseBody.ErrorMsg, expectedErrorMessage.ErrorMsg)
	})

	t.Run("LineDeleteCallback result 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/auth/line/delete/callback?code=test-code", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// gomock のコントローラ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userInfo := &config.LineUserInfo{
			Id:          "1234",
			UserId:      1,
			UserName:    "test@example.com",
			DisplayName: "test",
			LineToken:   MockResp,
		}

		response := utils.ErrorResponse{}

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)
		// ControllersCommonService のモックを作成
		mockControllersCommonService := mock_controllers_common.NewMockControllersCommonService(ctrl)
		// mockLineConfig
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		mockControllersCommonService.EXPECT().
			LineAuthCommon(gomock.Any(), gomock.Any()).
			Return(200, userInfo, response)

		mockLineService.EXPECT().
			RevokeLineAccessToken(gomock.Any()).
			Return(nil)

		resMock := []models.ExternalAuthData{
			{
				UserId:   1,
				UserName: "test@example.com",
			},
		}

		patches := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"GetExternalAuth",
			func(_ *models.SignDataFetcher, UserName string) ([]models.ExternalAuthData, error) {
				return resMock, nil
			})
		defer patches.Reset()

		patches1 := ApplyMethod(
			reflect.TypeOf(&models.SignDataFetcher{}),
			"DeleteSignIn",
			func(_ *models.SignDataFetcher, data models.RequestSignInDeleteData) error {
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
			UtilsFetcher:             mockUtilsFetcher,
			EmailTemplateService:     templates.NewEmailTemplateManager(),
			LineConfig:               mockLineService,
			ControllersCommonService: mockControllersCommonService,
		}
		lineManager.LineDeleteCallback(c)

		// ステータスコードの確認
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)

		var responseBody utils.ResponseWithSingle[string]
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		expectedOk := utils.ResponseWithSingle[string]{
			Result: "サインイン削除に成功",
		}
		assert.Equal(t, responseBody.Result, expectedOk.Result)
	})
}
