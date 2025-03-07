package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/common"
	"server/utils"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateErrorMessage はテスト用のエラーメッセージ構造体を生成する関数
func CreateErrorMessage(field string, message string) map[string]interface{} {
	return map[string]interface{}{
		"field":   field,
		"message": message,
	}
}

func SortErrorMessages(sortData []utils.ErrorMessages) {
	sort.SliceStable(
		sortData, func(i, j int) bool {
			return sortData[i].Field < sortData[j].Field
		},
	)
}

func SetupMockHTTPServer(statusCode int, query, responseBody string) (*http.Client, string, func()) {
	// モックHTTPサーバーの作成
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	}))

	// モックHTTPクライアントの作成
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				// モックサーバーのURLを設定
				return url.Parse(mockServer.URL)
			},
		},
	}

	// モックサーバーURL
	mockURL := mockServer.URL + query

	// クリーンアップ関数を返す
	cleanup := func() {
		mockServer.Close()
	}

	return client, mockURL, cleanup
}

func QueryUnescape(location string) (string, error) {
	// URLをパースする
	// スキームを手動で追加
	if strings.HasPrefix(location, ":///") {
		location = "http" + location[1:]
	}
	parsedURL, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	// クエリーパラメータを取得する
	queryParams := parsedURL.Query()

	// errorキーの値を取得する
	errorMessage := queryParams.Get("error")

	// デコード
	decoded, err := url.QueryUnescape(errorMessage)
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func RedirectSuccess(location string) (int, string, error) {
	// URLをパースする
	// スキームを手動で追加
	if strings.HasPrefix(location, ":///") {
		location = "http" + location[1:]
	}
	parsedURL, err := url.Parse(location)
	if err != nil {
		return 0, "", err
	}

	// クエリーパラメータを取得する
	queryParams := parsedURL.Query()

	userId := queryParams.Get("user_id")
	userEmail := queryParams.Get("user_email")

	common := common.NewCommonFetcher()
	id, _ := common.StrToInt(userId)

	return id, userEmail, nil
}

// 共通のリクエスト作成ヘルパー関数
func CreateTestRequest(method, path string, data interface{}, params map[string]string, headers ...map[string]string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// リクエストボディをJSON化（データがnilの場合は空）
	var body io.Reader
	switch v := data.(type) {
	case string:
		body = bytes.NewBufferString(v)
	case nil:
		body = nil
	default:
		// JSON 変換
		jsonData, _ := json.Marshal(v)
		body = bytes.NewBuffer(jsonData)
	}

	c.Request = httptest.NewRequest(method, path, body)

	defaultHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			defaultHeaders[key] = value
		}
	}

	for key, value := range defaultHeaders {
		c.Request.Header.Set(key, value)
	}

	// params を処理（空でないことを確認）
	if len(params) > 0 {
		for key, value := range params {
			c.Params = append(c.Params, gin.Param{Key: key, Value: value})
		}
	}

	return w, c
}
