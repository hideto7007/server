// config/db_config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DataSourceName string = getDataBaseSource()

func getDataBaseSource() string {
	// .env ファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	globalDbSouce := os.Getenv("GLOBALDBSOURCE")
	localDbSouce := os.Getenv("LOCALDBSOURCE")

	if globalDbSouce != "" {
		return globalDbSouce
	}
	return localDbSouce

	// TODO
	// 上記の設定はローカルのみ接続するようになっている。グローバルにするには、ssh接続を追加する必要がある
}
