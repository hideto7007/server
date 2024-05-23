// config/test_config.go
package config

import (
	"database/sql"
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
	testDB, err = sql.Open("postgres", "user=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// テーブルのリセットなどを行う
	_, err = testDB.Exec(`TRUNCATE TABLE income_data, payment_dates, other_tables RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatalf("Failed to reset test database: %v", err)
	}
}

func TeardownTestDatabase() {
	if testDB != nil {
		testDB.Close()
	}
}
