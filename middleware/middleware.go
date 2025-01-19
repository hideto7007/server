// middleware/middleware.go
package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"server/config"
	"server/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func JWTAuthMiddleware(utilsFetcher utils.UtilsDataFetcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		requestID, _ := c.Get("request_id")
		authToken, err1 := c.Cookie(utils.AuthToken)
		if authToken == "" || err1 != nil {
			response := utils.ErrorResponse{
				ErrorMsg: "トークンが必要です。",
			}
			logrus.WithField("request_id", requestID).Error(err1.Error())
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
			response := utils.ErrorResponse{
				ErrorMsg: "無効なトークンです",
			}
			logrus.WithField("request_id", requestID).Error(err.Error())
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
			config.GlobalEnv.RedirectPath,
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

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := uuid.New().String() // リクエストIDを生成
		c.Set("request_id", requestID)

		c.Writer.Header().Set("X-Request-ID", requestID)
		// リクエストボディを読み取る
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// リクエストボディを再設定（再利用可能にする）
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// リクエスト実行前ログ
		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"params":     c.Request.URL.Query(),
			"body":       requestBody,
			"client_ip":  c.ClientIP(),
		}).Info("Request started")

		c.Next()

		// レスポンス後ログ
		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"status":     c.Writer.Status(),
			"latency":    time.Since(startTime).Milliseconds(),
		}).Info("Request completed")

		// Gin フレームワークでエラーを記録する際にここでログに書き込む
		// c.Error(err)等で
		// 現状、このような仕様でエラーレスポンスを管理してなくログに書き込まれることはない
		// 今後、Gin　フレームワークでログ記録することがあればこちらを使用する
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// エラーログの出力
				logrus.WithFields(logrus.Fields{
					"request_id": c.GetString("request_id"),
					"client_ip":  c.ClientIP(),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"params":     c.Request.URL.Query(),
					"status":     c.Writer.Status(),
					"error":      e.Err.Error(),
					"stack":      string(debug.Stack()), // オプションでスタックトレース
				}).Error("APIエラー発生")
			}
		}
	}
}
