package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testHttpData struct {
		Data interface{} `json:"data"`
	}

	RequestTestData struct {
		TestName string `json:"test_name"`
	}
)

func TestRequest(t *testing.T) {
	t.Run("Request リクエスト作成失敗", func(t *testing.T) {
		// 無効なURLを渡す
		invalidURL := "http://\x7f\x7f\x7f"

		client := NewHTTPClient()

		// テスト実行
		resp, err := client.Request("GET", invalidURL, nil, nil)

		// エラーが発生していることを確認
		assert.NotNil(t, err)
		assert.Nil(t, resp)

		// エラーメッセージを確認
		expectedErrorMessage := "リクエスト作成失敗"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	})
	t.Run("Request リクエスト送信失敗", func(t *testing.T) {
		var testHttp HttpService = NewHTTPClient()
		req, err := testHttp.Request("test", "test", Headers, nil)

		assert.Nil(t, req)
		// エラーメッセージが含まれていることを確認
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "リクエスト送信失敗")
		assert.Contains(t, err.Error(), "unsupported protocol scheme")
	})
	t.Run("Request リクエスト成功", func(t *testing.T) {
		// モックサーバーを作成
		data := `{"status":"success"}`
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// リクエストヘッダーを確認
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// リクエストボディを読み取る
			body, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"key":"value"}`, string(body))

			// レスポンスを返す
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(data))
		}))
		defer mockServer.Close()

		client := NewHTTPClient()

		method := http.MethodPost
		url := mockServer.URL
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		body := bytes.NewReader([]byte(`{"key":"value"}`))

		// テスト実行
		resp, err := client.Request(method, url, headers, body)

		// エラーがないことを確認
		assert.Nil(t, err)

		// レスポンスコードを確認
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// レスポンスボディを確認
		responseBody, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, data, string(responseBody))
	})
}

func TestHttpMethod(t *testing.T) {
	// モックサーバーを起動
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "success"}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	TestURL := mockServer.URL // モックサーバーのURLを使用
	Data := testHttpData{
		Data: []RequestTestData{
			{
				TestName: "test",
			},
		},
	}

	Body, _ := json.Marshal(Data)

	t.Run("Getメソッドテスト", func(t *testing.T) {
		var testHttp HttpService = NewHTTPClient()
		req, err := testHttp.Get(TestURL, Headers)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, req.StatusCode)
	})
	t.Run("Postメソッドテスト", func(t *testing.T) {
		var testHttp HttpService = NewHTTPClient()
		req, err := testHttp.Post(TestURL, Headers, Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, req.StatusCode)
	})
	t.Run("Putメソッドテスト", func(t *testing.T) {
		var testHttp HttpService = NewHTTPClient()
		req, err := testHttp.Put(TestURL, Headers, Body)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, req.StatusCode)
	})
	t.Run("Deleteメソッドテスト", func(t *testing.T) {
		var testHttp HttpService = NewHTTPClient()
		req, err := testHttp.Delete(TestURL, Headers)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, req.StatusCode)
	})
}
