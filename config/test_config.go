// config/test_config.go
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var originalDBHost string

func Setup() {
	originalDBHost = os.Getenv("DB_HOST")
	os.Setenv("DB_HOST", "localhost") // テスト用のデータベースホストを設定
}

func Teardown() {
	os.Setenv("DB_HOST", originalDBHost)
}

var testDB *sql.DB

func SetupTestDatabase() {
	var err error
	// 環境変数から値を取得して接続文字列を作成
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("PSQL_USER"),        // ユーザー名
		os.Getenv("PSQL_PASSWORD"),    // パスワード
		os.Getenv("PSQL_TEST_DBNAME"), // テストDB名
		os.Getenv("PSQL_SSLMODEL"),    // SSLモード
	)

	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// テーブルのリセットなどを行う
	_, err = testDB.Exec(`TRUNCATE TABLE auth, users, income_forecast_data RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatalf("Failed to reset test database: %v", dsn)
	}
}

func TeardownTestDatabase() {
	if testDB != nil {
		testDB.Close()
	}
}
