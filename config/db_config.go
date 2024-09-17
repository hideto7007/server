// config/db_config.go
package config

import (
	"os"
)

var DataSourceName string = getDataBaseSource()

func getDataBaseSource() string {
	// .env ファイルを読み込む
	// 環境変数ファイルを動的に選択
	// envFile := ".env"
	// if os.Getenv("ENV") == "test" {
	// 	envFile = ".env.test"
	// 	fmt.Print(envFile)
	// }

	// err := godotenv.Load(envFile)
	// if err != nil {
	// 	log.Fatalf("Error loading %s file", envFile)
	// }

	globalDbSouce := os.Getenv("GLOBALDBSOURCE")
	localDbSouce := os.Getenv("LOCALDBSOURCE")

	if globalDbSouce != "" {
		return globalDbSouce
	}
	return localDbSouce

	// TODO
	// 上記の設定はローカルのみ接続するようになっている。グローバルにするには、ssh接続を追加する必要がある
}
