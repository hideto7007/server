package main

import (
	"log"
	"os"
	"server/middleware"
	"server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	log.Println("start server...")
	log.Println("REACT_CLIENT : " + os.Getenv("REACT_CLIENT"))
	log.Println("VUE_CLIENT : ", os.Getenv("VUE_CLIENT"))
	log.Println("DOCKER_CLIENT : ", os.Getenv("DOCKER_CLIENT"))
	log.Println("SWAGGER_CLIENT : ", os.Getenv("SWAGGER_CLIENT"))
	corsMiddleware := middleware.CORSMiddleware()

	// CORSミドルウェアを設定
	r.Use(corsMiddleware)
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204) // No Content（Preflightリクエストに対する適切なレスポンス）
	})
	routes.SetupRoutes(r)

	r.Run(":8080")
}
