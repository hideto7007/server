// middleware/middleware.go
package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"server/enum"
	"server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "トークンが必要です。"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// トークンの検証
		_, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "トークンの有効期限が切れています"})
			c.Abort()
			return
		}

		// トークンのクレームを取得
		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	if exp, ok := claims["exp"].(float64); ok {
		// 		if time.Unix(int64(exp), 0).Before(time.Now()) {
		// 			c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "トークンの有効期限が切れています"})
		// 			c.Abort()
		// 			return
		// 		}
		// 	}
		// } else {
		// 	c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "有効期限が切れて無効なトークンです"})
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{
			os.Getenv("REACT_CLIENT"),
			os.Getenv("VUE_CLIENT"),
			os.Getenv("DOCKER_CLIENT"),
			os.Getenv("SWAGGER_CLIENT1"),
			os.Getenv("SWAGGER_CLIENT2"),
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Length",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}

	// 設定内容をログに出力
	log.Println("CORS Configuration:")
	log.Printf("Allowed Origins: %v", config.AllowOrigins)
	log.Printf("Allowed Methods: %v", config.AllowMethods)
	log.Printf("Allowed Headers: %v", config.AllowHeaders)
	log.Printf("Allow Credentials: %v", config.AllowCredentials)
	log.Printf("MaxAge: %v", config.MaxAge)

	return cors.New(config)
}
