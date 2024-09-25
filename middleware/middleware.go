// middleware/middleware.go
package middleware

import (
	"fmt"
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
		fmt.Print("authHeader", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "トークンが必要です"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		// トークンの検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{enum.ERROR: "無効なトークンです"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// .env ファイルを読み込む（必要であれば）
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }
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
