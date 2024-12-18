package test_utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/utils"
	"sort"
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