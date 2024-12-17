package controllers

import (
	"net/http"
	"net/http/httptest"

	"server/templates"
	"server/utils"
	"testing"

	mock_config "server/mock/config"

	// . "github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"

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

// func TestGoogleSignInCallback(t *testing.T) {

// 	gin.SetMode(gin.TestMode)

// 	mockRedirectURL := "http://localhost:8080/auth/google/callback"
// 	GooglePramsVaild := GooglePrams{}
// 	GooglePrams := GooglePrams{
// 		Code:        "1234",
// 		RedirectUri: mockRedirectURL,
// 	}

// 	// oauth2.Configのモックを生成
// 	MockOauthConfig := &oauth2.Config{
// 		ClientID:     "mock_client_id",
// 		ClientSecret: "mock_client_secret",
// 		RedirectURL:  mockRedirectURL,
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.email",
// 			"https://www.googleapis.com/auth/userinfo.profile",
// 		},
// 		Endpoint: google.Endpoint,
// 	}

// 	t.Run("GoogleSignInCallback バリデーション 必須", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
// 		c.Request = httptest.NewRequest("GET", "/auth/google/signin/callback?code=test-code", nil)
// 		c.Request.Header.Set("Content-Type", "application/json")

// 		// gomock のコントローラ作成
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		mockGoogleService := mock_confidence.NewMockGoogleConfig(ctrl)

// 		response := utils.ErrorResponse{
// 			ErrorMsg: "テストエラー",
// 		}

// 		mockGoogleService.EXPECT().
// 			GoogleAuthCommon(gomock.Any(), gomock.Any()).
// 			Return(500, GoogleUserInfo{}, response)

// 		googleManager := GoogleManager{
// 			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
// 			EmailTemplateService: templates.NewEmailTemplateManager(),
// 			GoogleConfig:        config.NewGoogleManager(),
// 		}
// 		status, googleUserInfo, response := googleManager.GoogleSignInCallback(c)

// 		// ステータスコードの確認
// 		assert.Equal(t, http.StatusBadRequest, status)
// 		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

// 		expectedErrorMessage := utils.ErrorResponse{
// 			Result: []utils.ErrorMessages{
// 				{
// 					Field:   "code",
// 					Message: "コードは必須です。",
// 				},
// 				{
// 					Field:   "redirect_uri",
// 					Message: "リダイレクトは必須です。",
// 				},
// 			},
// 		}
// 		assert.Equal(t, response, expectedErrorMessage)
// 	})

// 	t.Run("GoogleSignInCallback 認証コードからトークン取得エラー", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
// 		c.Request = httptest.NewRequest("GET", "auth/google/signin/callback", nil)
// 		c.Request.Header.Set("Content-Type", "application/json")
// 		var mockToken *oauth2.Token

// 		// gomock のコントローラを作成
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		// mockGoogleService のモックを作成
// 		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

// 		// モックの挙動を定義
// 		mockGoogleService.EXPECT().
// 			GoogleOauth(gomock.Any()).
// 			Return(MockOauthConfig)

// 		mockGoogleService.EXPECT().
// 			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockToken, fmt.Errorf("トークン取得エラー"))

// 		googleManager := GoogleManager{
// 			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
// 			EmailTemplateService: templates.NewEmailTemplateManager(),
// 			GoogleConfig:        mockGoogleService,
// 		}
// 		status, googleUserInfo, response := googleManager.GoogleSignInCallback(c)

// 		// ステータスコードの確認
// 		assert.Equal(t, http.StatusInternalServerError, status)
// 		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

// 		expectedErrorMessage := utils.ErrorResponse{
// 			ErrorMsg: "トークン取得エラー",
// 		}
// 		assert.Equal(t, response, expectedErrorMessage)
// 	})

// 	t.Run("GoogleSignInCallback トークンを使ってユーザー情報を取得エラー", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
// 		c.Request = httptest.NewRequest("GET", "auth/google/signin/callback", nil)
// 		c.Request.Header.Set("Content-Type", "application/json")
// 		// ダミーのトークンとレスポンス
// 		mockToken := &oauth2.Token{
// 			AccessToken:  "test-access-token",
// 			RefreshToken: "test-refresh-token",
// 			TokenType:    "Bearer",
// 		}

// 		var mockClient *http.Client

// 		mockResp := &http.Response{
// 			StatusCode: http.StatusOK,
// 		}

// 		// gomock のコントローラを作成
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		// mockGoogleService のモックを作成
// 		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

// 		// モックの挙動を定義
// 		mockGoogleService.EXPECT().
// 			GoogleOauth(gomock.Any()).
// 			Return(MockOauthConfig)

// 		mockGoogleService.EXPECT().
// 			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockToken, nil)

// 		mockGoogleService.EXPECT().
// 			Client(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockClient)

// 		mockGoogleService.EXPECT().
// 			Get(gomock.Any(), gomock.Any()).
// 			Return(mockResp, fmt.Errorf("ユーザー取得エラー"))

// 		googleManager := GoogleManager{
// 			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
// 			EmailTemplateService: templates.NewEmailTemplateManager(),
// 			GoogleConfig:        mockGoogleService,
// 		}
// 		status, googleUserInfo, response := googleManager.GoogleSignInCallback(c)

// 		// ステータスコードの確認
// 		assert.Equal(t, http.StatusInternalServerError, status)
// 		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

// 		expectedErrorMessage := utils.ErrorResponse{
// 			ErrorMsg: "ユーザー取得エラー",
// 		}
// 		assert.Equal(t, response, expectedErrorMessage)
// 	})

// 	t.Run("GoogleSignInCallback デコード時エラー", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
// 		c.Request = httptest.NewRequest("GET", "auth/google/signin/callback", nil)
// 		c.Request.Header.Set("Content-Type", "application/json")
// 		// ダミーのトークンとレスポンス
// 		mockToken := &oauth2.Token{
// 			AccessToken:  "test-access-token",
// 			RefreshToken: "test-refresh-token",
// 			TokenType:    "Bearer",
// 		}

// 		var mockClient *http.Client

// 		mockResp := &http.Response{
// 			StatusCode: http.StatusOK,
// 			Body:       io.NopCloser(strings.NewReader("test")),
// 		}

// 		// gomock のコントローラを作成
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		// mockGoogleService のモックを作成
// 		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

// 		// モックの挙動を定義
// 		mockGoogleService.EXPECT().
// 			GoogleOauth(gomock.Any()).
// 			Return(MockOauthConfig)

// 		mockGoogleService.EXPECT().
// 			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockToken, nil)

// 		mockGoogleService.EXPECT().
// 			Client(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockClient)

// 		mockGoogleService.EXPECT().
// 			Get(gomock.Any(), gomock.Any()).
// 			Return(mockResp, nil)

// 		googleManager := GoogleManager{
// 			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
// 			EmailTemplateService: templates.NewEmailTemplateManager(),
// 			GoogleConfig:        mockGoogleService,
// 		}
// 		status, googleUserInfo, response := googleManager.GoogleSignInCallback(c)

// 		// ステータスコードの確認
// 		assert.Equal(t, http.StatusInternalServerError, status)
// 		userInfo := GoogleUserInfo{
// 			ID:            "",
// 			UserId:        0,
// 			Email:         "",
// 			VerifiedEmail: false,
// 			Name:          "",
// 			GivenName:     "",
// 			FamilyName:    "",
// 			Picture:       "",
// 			Locale:        "",
// 			Token:         mockToken,
// 		}
// 		assert.Equal(t, googleUserInfo, userInfo)

// 		expectedErrorMessage := utils.ErrorResponse{
// 			ErrorMsg: "invalid character 'e' in literal true (expecting 'r')",
// 		}
// 		assert.Equal(t, response, expectedErrorMessage)
// 	})

// 	t.Run("GoogleSignInCallback 成功", func(t *testing.T) {
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
// 		c.Request = httptest.NewRequest("GET", "auth/google/signin/callback", nil)
// 		c.Request.Header.Set("Content-Type", "application/json")
// 		// ダミーのトークンとレスポンス
// 		mockToken := &oauth2.Token{
// 			AccessToken:  "test-access-token",
// 			RefreshToken: "test-refresh-token",
// 			TokenType:    "Bearer",
// 		}

// 		var mockClient *http.Client

// 		userInfo := GoogleUserInfo{
// 			ID:            "1234",
// 			UserId:        1,
// 			Email:         "test@example.com",
// 			VerifiedEmail: true,
// 			Name:          "test",
// 			GivenName:     "test",
// 			FamilyName:    "test",
// 			Picture:       "111",
// 			Locale:        "222",
// 			Token:         mockToken,
// 		}
// 		// JSONエンコード
// 		userInfoBytes, err := json.Marshal(userInfo)
// 		if err != nil {
// 			t.Fatalf("failed to marshal userInfo: %v", err)
// 		}

// 		mockResp := &http.Response{
// 			StatusCode: http.StatusOK,
// 			Body:       io.NopCloser(bytes.NewReader(userInfoBytes)),
// 		}

// 		// gomock のコントローラを作成
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		// mockGoogleService のモックを作成
// 		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

// 		// モックの挙動を定義
// 		mockGoogleService.EXPECT().
// 			GoogleOauth(gomock.Any()).
// 			Return(MockOauthConfig)

// 		mockGoogleService.EXPECT().
// 			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockToken, nil)

// 		mockGoogleService.EXPECT().
// 			Client(gomock.Any(), gomock.Any(), gomock.Any()).
// 			Return(mockClient)

// 		mockGoogleService.EXPECT().
// 			Get(gomock.Any(), gomock.Any()).
// 			Return(mockResp, nil)

// 		googleManager := GoogleManager{
// 			UtilsFetcher:         utils.NewUtilsFetcher(utils.JwtSecret),
// 			EmailTemplateService: templates.NewEmailTemplateManager(),
// 			GoogleConfig:        mockGoogleService,
// 		}
// 		status, googleUserInfo, response := googleManager.GoogleSignInCallback(c)

// 		// ステータスコードの確認
// 		assert.Equal(t, http.StatusOK, status)
// 		expectedUserInfo := userInfo
// 		assert.Equal(t, googleUserInfo, expectedUserInfo)

// 		expectedErrorMessage := utils.ErrorResponse{}
// 		assert.Equal(t, response, expectedErrorMessage)
// 	})
// }
