package routes

import (
	"server/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORSミドルウェアの設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}                   // 許可するオリジンを指定
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // 許可するHTTPメソッドを指定
	config.AllowHeaders = []string{"Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers"}

	// CORSミドルウェアをルーターに追加
	r.Use(cors.New(config))

	// ルートの設定
	Routes := r.Group("/api")
	{
		Routes.GET("/price", controllers.GetPriceInfoApi)
		Routes.GET("/income_data", controllers.GetIncomeDataInRangeApi)
		// 他のエンドポイントのルーティングもここで設定
	}
}
