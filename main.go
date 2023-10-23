package main

import (
	"server/routes"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 許可するオリジンを指定
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		// c.Header("Access-Control-Allow-Origin", "https://incomeflowpro.net/")
		// 許可するHTTPメソッドを指定
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 許可するカスタムヘッダーを指定
		c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
		// クライアントがクッキーを送信できるように設定
		c.Header("Access-Control-Allow-Credentials", "true")
		// 通常のリクエストに対してはミドルウェアを通過させる
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// CORSミドルウェアを設定
	r.Use(CORSMiddleware())

	routes.SetupRoutes(r)
	r.Run(":8080")
}
