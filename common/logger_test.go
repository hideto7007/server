package common

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type LoggerEntry struct {
	Level       string `json:"level"`
	Msg         string `json:"msg"`
	Time        string `json:"time"`
	ErrorDetail string `json:"error_detail"`
}

func TestInitLogger(t *testing.T) {
	t.Run("InitLogger エラー発生しないこと", func(t *testing.T) {
		// テスト用のログファイル名を設定
		testLogFile := "test_takuwaeru.log"

		// テスト前に既存のテスト用ファイルを削除（存在する場合）
		_ = os.Remove(testLogFile)

		// InitLoggerを呼び出す
		InitLogger(testLogFile)

		// ログに書き込むテストメッセージ
		logrus.Info("Test log message")

		fp, err := os.Open(testLogFile)
		assert.NoError(t, err, "ログファイルのオープンに失敗しました")
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			var entry LoggerEntry
			err := json.Unmarshal([]byte(scanner.Text()), &entry)
			assert.NoError(t, err, "ログのJSONデコードに失敗しました")

			assert.Equal(t, entry.Level, "info")
			assert.Equal(t, entry.Msg, "Test log message")
			assert.NotEmpty(t, entry.Time)
		}
	})
	t.Run("InitLogger エラー発生すること", func(t *testing.T) {
		// テスト用のログファイル名を設定
		testLogFile := "test_takuwaeru.log"

		// テスト前に既存のテスト用ファイルを削除（存在する場合）
		_ = os.Remove(testLogFile)

		// InitLoggerを呼び出す
		InitLogger(testLogFile)

		// エラーメッセージをログに書き込む
		logrus.WithField("error_detail", "example error").Error("Test error log message")

		fp, err := os.Open(testLogFile)
		assert.NoError(t, err, "ログファイルのオープンに失敗しました")
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			var entry LoggerEntry
			err := json.Unmarshal([]byte(scanner.Text()), &entry)
			assert.NoError(t, err, "ログのJSONデコードに失敗しました")

			assert.Equal(t, entry.Level, "error")
			assert.Equal(t, entry.Msg, "Test error log message")
			assert.Equal(t, entry.ErrorDetail, "example error")
			assert.NotEmpty(t, entry.Time)
		}
	})
	t.Run("ログの中身削除", func(t *testing.T) {
		// テスト用のログファイル名を設定
		testLogFile := "test_takuwaeru.log"

		// テスト前に既存のテスト用ファイルを削除（存在する場合）
		_ = os.Remove(testLogFile)
	})
}
