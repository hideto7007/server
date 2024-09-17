package main

import (
	"log"
	"os"
	"server/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// .env ファイルを読み込む（必要であれば）
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }s
		// 許可したいアクセス元
		AllowOrigins: []string{
			os.Getenv("REACT_CLIENT"),
			os.Getenv("VUE_CLIENT"),
			os.Getenv("DOCKER_CLIENT"),
		},
		// 許可したいHTTPメソッド
		AllowMethods: []string{
			"GET",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報の要否
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	})
}

func main() {
	r := gin.Default()
	log.Println("start server...")

	// CORSミドルウェアを設定
	r.Use(CORSMiddleware())

	routes.SetupRoutes(r)
	r.Run(":8080")
}
