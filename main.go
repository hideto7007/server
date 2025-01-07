package main

import (
	"log"
	"server/common"
	"server/config"
	"server/middleware"
	"server/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .envファイルを最初にロードする
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 環境変数の初期化
	config.InitGoogleEnvs()

	common.InitLogger(config.GlobalEnv.OutPutLoggerFile)

	r := gin.Default()
	log.Println("start server...")
	corsMiddleware := middleware.CORSMiddleware()

	// CORSミドルウェアを設定
	r.Use(corsMiddleware)
	// ログ監視
	r.Use(middleware.RequestLoggerMiddleware())
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204) // No Content（Preflightリクエストに対する適切なレスポンス）
	})
	routes.SetupRoutes(r)

	r.Run(":8080")
}
