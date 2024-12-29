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

	mock_common "server/mock/common"
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
		test_utils.SortErrorMessages(response.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
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
			UserName:      "",
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
			UserName:      "test@example.com",
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

func TestLineAuthCommon(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockRedirectURL := "http://localhost:8080/auth/line/callback"
	LinePramsVaild := LinePrams{}
	var EnptyUserInfo *config.LineUserInfo
	LinePrams := LinePrams{
		Code:        "1234",
		State:       "dummy_state",
		RedirectUri: mockRedirectURL,
	}

	t.Run("LineAuthCommon バリデーション 必須", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockMockHttpService のモックを作成
		mockMockHttpService := mock_common.NewMockHttpService(ctrl)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		lineManager := ControllersCommonManager{
			LineConfig: config.NewLineManager(mockMockHttpService),
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePramsVaild)

		// ステータスコードの確認
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, lineUserInfo, EnptyUserInfo)

		expectedErrorMessage := utils.ErrorResponse{
			Result: []utils.ErrorMessages{
				{
					Field:   "code",
					Message: "コードは必須です。",
				},
				{
					Field:   "state",
					Message: "ステートは必須です。",
				},
				{
					Field:   "redirect_uri",
					Message: "リダイレクトは必須です。",
				},
			},
		}
		test_utils.SortErrorMessages(response.Result)
		test_utils.SortErrorMessages(expectedErrorMessage.Result)
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("LineAuthCommon 無効なステートです。エラー系", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockMockHttpService のモックを作成
		mockMockHttpService := mock_common.NewMockHttpService(ctrl)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		lineManager := ControllersCommonManager{
			LineConfig: config.NewLineManager(mockMockHttpService),
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, lineUserInfo, EnptyUserInfo)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "無効なステートです。",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("LineAuthCommon 無効なLineアクセストークンです。エラー系", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  utils.OauthState,
			Value: "dummy_state",
		})
		var mockToken *config.LineTokenResponse

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().
			GetLineAccessToken(gomock.Any(), gomock.Any()).
			Return(mockToken, fmt.Errorf("トークン取得エラー"))

		lineManager := ControllersCommonManager{
			LineConfig: mockLineService,
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, lineUserInfo, EnptyUserInfo)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "無効なLineアクセストークンです。",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("LineAuthCommon Lineのユーザー情報取得に失敗しました。エラー系", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  utils.OauthState,
			Value: "dummy_state",
		})

		mockResp := &config.LineTokenResponse{
			AccessToken: "token",
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().
			GetLineAccessToken(gomock.Any(), gomock.Any()).
			Return(mockResp, nil)

		mockLineService.EXPECT().
			GetLineUserInfo(gomock.Any()).
			Return(nil, fmt.Errorf("ユーザー取得エラー"))

		lineManager := ControllersCommonManager{
			LineConfig: mockLineService,
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, lineUserInfo, EnptyUserInfo)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "Lineのユーザー情報取得に失敗しました。",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("LineAuthCommon メールアドレス取得失敗", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  utils.OauthState,
			Value: "dummy_state",
		})
		// ダミーのトークンとレスポンス
		mockResp := &config.LineTokenResponse{
			AccessToken: "token",
			IdToken:     "id_token",
		}

		userInfo := &config.LineUserInfo{
			Id:          "1234test",
			UserId:      1,
			UserName:    "",
			DisplayName: "test",
		}

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().
			GetLineAccessToken(gomock.Any(), gomock.Any()).
			Return(mockResp, nil)

		mockLineService.EXPECT().
			GetLineUserInfo(gomock.Any()).
			Return(userInfo, nil)

		mockLineService.EXPECT().
			GetEmail(gomock.Any()).
			Return("", fmt.Errorf("メールアドレス取得エラー"))

		lineManager := ControllersCommonManager{
			LineConfig: mockLineService,
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, lineUserInfo, EnptyUserInfo)

		expectedErrorMessage := utils.ErrorResponse{
			ErrorMsg: "メールアドレス取得エラー",
		}
		assert.Equal(t, response, expectedErrorMessage)
	})

	t.Run("LineAuthCommon 成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// パラメータなしのリクエストを送信して、不正なリクエストをシミュレート
		c.Request = httptest.NewRequest("GET", "/line_auth_common", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.AddCookie(&http.Cookie{
			Name:  utils.OauthState,
			Value: "dummy_state",
		})

		// ダミーのトークンとレスポンス
		mockResp := &config.LineTokenResponse{
			AccessToken: "token",
			IdToken:     "id_token",
		}

		userInfo := &config.LineUserInfo{
			Id:          "1234test",
			UserId:      1,
			UserName:    "",
			DisplayName: "test",
		}

		email := "test@example.com"

		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// mockLineService のモックを作成
		mockLineService := mock_config.NewMockLineConfig(ctrl)

		// モックの挙動を定義
		mockLineService.EXPECT().
			GetLineAccessToken(gomock.Any(), gomock.Any()).
			Return(mockResp, nil)

		mockLineService.EXPECT().
			GetLineUserInfo(gomock.Any()).
			Return(userInfo, nil)

		mockLineService.EXPECT().
			GetEmail(gomock.Any()).
			Return("test@example.com", nil)

		lineManager := ControllersCommonManager{
			LineConfig: mockLineService,
		}
		status, lineUserInfo, response := lineManager.LineAuthCommon(c, LinePrams)

		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, status)
		expectedUserInfo := userInfo
		expectedUserInfo.UserName = email
		expectedUserInfo.LineToken = mockResp
		assert.Equal(t, lineUserInfo, expectedUserInfo)

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
