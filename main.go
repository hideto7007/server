package main

import (
	// "net/http"

	"github.com/gin-contrib/cors" // corsミドルウェアのインポート
	"github.com/gin-gonic/gin"

	"server/routes"
)

func main() {
	r := gin.Default()

	// CORSミドルウェアを追加
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers"}
	r.Use(cors.New(config))

	// CORSのオプションリクエストに対するハンドラーを追加
	// r.OPTIONS("/api/price", func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers")
	// 	c.Status(http.StatusOK)
	// })

	routes.SetupRoutes(r)
	r.Run(":8080")
}
