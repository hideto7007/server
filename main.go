package main

import (
	"log"
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

	r := gin.Default()
	log.Println("start server...")
	corsMiddleware := middleware.CORSMiddleware()

	// CORSミドルウェアを設定
	r.Use(corsMiddleware)
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204) // No Content（Preflightリクエストに対する適切なレスポンス）
	})
	routes.SetupRoutes(r)

	r.Run(":8080")
}
