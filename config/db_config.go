// config/db_config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DataSourceName string = getDataBaseSource()

func getDataBaseSource() string {
	// .env ファイルを読み込む
	// 環境変数 ENV が "test" の場合は .env ファイルの読み込みをスキップ

	if os.Getenv("PSQL_ENV") != "test" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	dsn := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_HOST"),
		os.Getenv("PSQL_PORT"),
		os.Getenv("PSQL_SSLMODEL"),
	)

	log.Println(dsn)

	return dsn

	// TODO
	// 上記の設定はローカルのみ接続するようになっている。グローバルにするには、ssh接続を追加する必要がある
}
