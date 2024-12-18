package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"server/config"
	"server/test_utils"
	"server/utils"
	"strings"
	"testing"

	mock_config "server/mock/config"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestGoogleAuthCommon(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRedirectURL := "http://localhost:8080/auth/google/callback"
	GooglePramsVaild := GooglePrams{}
	GooglePrams := GooglePrams{
		Code:        "1234",
		RedirectUri: mockRedirectURL,
	}

	// oauth2.Configのモックを生成
	MockOauthConfig := &oauth2.Config{
		ClientID:     "mock_client_id",
		ClientSecret: "mock_client_secret",
		RedirectURL:  mockRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	t.Run("GoogleAuthCommon バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/google_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		googleManager := ControllersCommonManager{
			GoogleConfig: config.NewGoogleManager(),
		}
		status, googleUserInfo, response := googleManager.GoogleAuthCommon(c, GooglePramsVaild)

		// ステータスコードの確認
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

		expectedErrorMessage := utils.ErrorResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "code",
					Message: "コードは必須です。",
				},
				{
					Field:   "redirect_uri",
					Message: "リダイレクトは必須です。",
				},
			},
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("GoogleAuthCommon 認証コードからトークン取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/google_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		var mockToken *oauth2.Token

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().
			GoogleOauth(gomock.Any()).
			Return(MockOauthConfig)

		mockGoogleService.EXPECT().
			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockToken, fmt.Errorf("トークン取得エラー"))

		googleManager := ControllersCommonManager{
			GoogleConfig: mockGoogleService,
		}
		status, googleUserInfo, response := googleManager.GoogleAuthCommon(c, GooglePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "トークン取得エラー",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("GoogleAuthCommon トークンを使ってユーザー情報を取得エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/google_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		// ダミーのトークンとレスポンス
		mockToken := &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			TokenType:    "Bearer",
		}

		var mockClient *http.Client

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().
			GoogleOauth(gomock.Any()).
			Return(MockOauthConfig)

		mockGoogleService.EXPECT().
			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockToken, nil)

		mockGoogleService.EXPECT().
			Client(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockClient)

		mockGoogleService.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Return(mockResp, fmt.Errorf("ユーザー取得エラー"))

		googleManager := ControllersCommonManager{
			GoogleConfig: mockGoogleService,
		}
		status, googleUserInfo, response := googleManager.GoogleAuthCommon(c, GooglePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, googleUserInfo, GoogleUserInfo{})

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "ユーザー取得エラー",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("GoogleAuthCommon デコード時エラー", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/google_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		// ダミーのトークンとレスポンス
		mockToken := &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			TokenType:    "Bearer",
		}

		var mockClient *http.Client

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("test")),
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().
			GoogleOauth(gomock.Any()).
			Return(MockOauthConfig)

		mockGoogleService.EXPECT().
			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockToken, nil)

		mockGoogleService.EXPECT().
			Client(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockClient)

		mockGoogleService.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Return(mockResp, nil)

		googleManager := ControllersCommonManager{
			GoogleConfig: mockGoogleService,
		}
		status, googleUserInfo, response := googleManager.GoogleAuthCommon(c, GooglePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		userInfo := GoogleUserInfo{
			ID:            "",
			UserId:        0,
			Email:         "",
			VerifiedEmail: false,
			Name:          "",
			GivenName:     "",
			FamilyName:    "",
			Picture:       "",
			Locale:        "",
			Token:         mockToken,
		}
		assert.Equal(t, googleUserInfo, userInfo)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "invalid character 'e' in literal true (expecting 'r')",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("GoogleAuthCommon 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/google_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		// ダミーのトークンとレスポンス
		mockToken := &oauth2.Token{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
			TokenType:    "Bearer",
		}

		var mockClient *http.Client

		userInfo := GoogleUserInfo{
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
		// JSONエンコード
		userInfoBytes, err := json.Marshal(userInfo)
		if err != nil {
			t.Fatalf("failed to marshal userInfo: %v", err)
		}

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(userInfoBytes)),
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockGoogleService のモックを作成
		mockGoogleService := mock_config.NewMockGoogleConfig(ctrl)

		// モックの挙動を定義
		mockGoogleService.EXPECT().
			GoogleOauth(gomock.Any()).
			Return(MockOauthConfig)

		mockGoogleService.EXPECT().
			Exchange(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockToken, nil)

		mockGoogleService.EXPECT().
			Client(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockClient)

		mockGoogleService.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Return(mockResp, nil)

		googleManager := ControllersCommonManager{
			GoogleConfig: mockGoogleService,
		}
		status, googleUserInfo, response := googleManager.GoogleAuthCommon(c, GooglePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, status)
		expectedUserInfo := userInfo
		assert.Equal(t, googleUserInfo, expectedUserInfo)

		expectedErrorMessage := utils.ErrorResponse{}
		assert.Equal(t, response, expectedErrorMessage)
	})
}

func TestGetRevoke(t *testing.T) {

	gin.SetMode(gin.TestMode)

	AccessToken := "test_token"

	t.Run("GetRevoke 異常系", func(t *testing.T) {
		// モックHTTPサーバを作成
		client, mockURL, cleanup := test_utils.SetupMockHTTPServer(
			http.StatusInternalServerError,
			"/dummry?",
			"エラー",
		)
		defer cleanup()

		googleManager := ControllersCommonManager{
			GoogleConfig: config.NewGoogleManager(),
		}
		resp, err := googleManager.GetRevoke(client, mockURL, AccessToken)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NoError(t, err) // HTTPリクエスト自体は成功するためerrはnilになる
	})

	t.Run("GetRevoke 正常系", func(t *testing.T) {

		// モックHTTPサーバを作成
		client, mockURL, cleanup := test_utils.SetupMockHTTPServer(
			http.StatusOK,
			"/dummry?",
			"成功",
		)
		defer cleanup()

		googleManager := ControllersCommonManager{
			GoogleConfig: config.NewGoogleManager(),
		}
		resp, err := googleManager.GetRevoke(client, mockURL, AccessToken)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, err, nil)
	})
}
