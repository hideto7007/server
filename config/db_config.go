// config/db_config.go
package config

import (
	"fmt"
	"log"
)

func GetDataBaseSource() string {
	dsn := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		GlobalEnv.PsqlUser,
		GlobalEnv.PsqlDbname,
		GlobalEnv.PsqlPassword,
		GlobalEnv.PsqlHost,
		GlobalEnv.PsqlPort,
		GlobalEnv.PsqlSslModel,
	)

	log.Println(dsn)

	return dsn

	// TODO
	// 上記の設定はローカルのみ接続するようになっている。グローバルにするには、ssh接続を追加する必要がある
}
