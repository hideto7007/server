// middleware/middleware.go
package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"server/config"
	"server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(utilsFetcher utils.UtilsDataFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		authToken, err1 := c.Cookie(utils.AuthToken)
		if authToken == "" || err1 != nil {
			log.Println("ERROR : ", err1.Error())
			response := utils.ErrorResponse{
				ErrorMsg: "トークンが必要です。",
			}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authToken, "Bearer ", "", 1)

		// トークンの検証
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			log.Println("ERROR: ", err)
			response := utils.ErrorResponse{
				ErrorMsg: "無効なトークンです",
			}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// トークンのクレームを取得して、有効期限を確認
		if claimsInterface, ok := utilsFetcher.MapClaims(token); ok {
			claims, _ := claimsInterface.(jwt.MapClaims)
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					response := utils.ErrorResponse{
						ErrorMsg: "トークンの有効期限が切れています",
					}
					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
					return
				}
			} else {
				response := utils.ErrorResponse{
					ErrorMsg: "トークンの有効期限が不正です",
				}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		} else {
			response := utils.ErrorResponse{
				ErrorMsg: "トークンのクレームが不正です",
			}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// トークンが有効な場合、リクエストを次に進める
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{
			config.GlobalEnv.ReactClient,
			config.GlobalEnv.VueClient,
			config.GlobalEnv.SwaggerClient,
			config.GlobalEnv.DockerClient,
			config.GlobalEnv.GoogleAccounts,
			config.GlobalEnv.GoogleApis,
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
