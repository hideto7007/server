package test_utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/common"
	"server/utils"
	"sort"
	"strings"
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
	userName := queryParams.Get("user_name")

	common := common.NewCommonFetcher()
	id, _ := common.StrToInt(userId)

	return id, userName, nil
}
