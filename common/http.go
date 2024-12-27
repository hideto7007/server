package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	HttpService interface {
		Request(method, url string, headers map[string]string, body io.Reader) (*http.Response, error)
		Get(url string, headers map[string]string) (*http.Response, error)
		Post(url string, headers map[string]string, body []byte) (*http.Response, error)
		Put(url string, headers map[string]string, body []byte) (*http.Response, error)
		Delete(url string, headers map[string]string) (*http.Response, error)
	}

	HttpServiceManager struct{}

	// HTTPClient はHTTPリクエストを扱う構造体
	HTTPClient struct {
		client *http.Client
	}
)

var Headers = map[string]string{
	"Content-Type": "application/json",
}

// タイムアウト設定
var timeout time.Duration = 10 * time.Second

// NewHTTPClient はHTTPクライアントを初期化する
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Request は汎用HTTPリクエスト関数
func (hc *HTTPClient) Request(method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	// 新しいリクエストを作成
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗: %w", err)
	}

	// ヘッダーを設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// リクエストを送信
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエスト送信失敗: %w", err)
	}

	return resp, nil
}

// Get はGETリクエストのラッパー
func (hc *HTTPClient) Get(url string, headers map[string]string) (*http.Response, error) {
	return hc.Request(http.MethodGet, url, headers, nil)
}

// Post はPOSTリクエストのラッパー
func (hc *HTTPClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return hc.Request(http.MethodPost, url, headers, bytes.NewReader(body))
}

// Put はPUTリクエストのラッパー
func (hc *HTTPClient) Put(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return hc.Request(http.MethodPut, url, headers, bytes.NewReader(body))
}

// Delete はDELETEリクエストのラッパー
func (hc *HTTPClient) Delete(url string, headers map[string]string) (*http.Response, error) {
	return hc.Request(http.MethodDelete, url, headers, nil)
}
