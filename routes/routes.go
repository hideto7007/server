package routes

import (
	"server/common"
	"server/config"
	"server/controllers"
	"server/middleware"
	"server/templates"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// APiインターフェイスのインスタンス定義
	var signAPI controllers.SignDataFetcher = controllers.NewSignDataFetcher(
		utils.NewUtilsFetcher(utils.JwtSecret),
		common.NewCommonFetcher(),
		templates.NewEmailTemplateManager(),
		config.NewRedisManager(),
	)
	var googleApi controllers.GoogleService = controllers.NewGoogleService(
		config.NewGoogleManager(),
	)
	var priceAPI controllers.PriceManagementFetcher = controllers.NewPriceManagementFetcher(
		common.NewCommonFetcher(),
	)
	var incomeAPI controllers.IncomeDataFetcher = controllers.NewIncomeDataFetcher(
		common.NewCommonFetcher(),
	)

	// google認証
	r.GET("auth/google/signin", googleApi.HandleGoogleSignIn)
	r.GET("auth/google/callback", googleApi.HandleGoogleCallback)

	// ルートの設定
	Routes := r.Group("/api")
	{
		Routes.POST("/signin", signAPI.PostSignInApi)
		Routes.GET("/refresh_token", signAPI.GetRefreshTokenApi)
		Routes.POST("/temporay_signup", signAPI.TemporayPostSignUpApi)
		Routes.GET("/retry_auth_email", signAPI.RetryAuthEmail)
		Routes.POST("/signup", signAPI.PostSignUpApi)
		Routes.PUT("/signin_edit", signAPI.PutSignInEditApi)
		Routes.DELETE("/signin_delete", signAPI.DeleteSignInApi)
		Routes.GET("/signout", signAPI.SignOutApi)

		// 認証が必要なルートにミドルウェアを追加
		authRoutes := Routes.Group("/")
		authRoutes.Use(middleware.JWTAuthMiddleware(utils.UtilsDataFetcher{}))
		{
			authRoutes.GET("/price", priceAPI.GetPriceInfoApi)
			authRoutes.GET("/income_data", incomeAPI.GetIncomeDataInRangeApi)
			authRoutes.GET("/range_date", incomeAPI.GetDateRangeApi)
			authRoutes.GET("/years_income_date", incomeAPI.GetYearIncomeAndDeductionApi)
			authRoutes.POST("/income_create", incomeAPI.InsertIncomeDataApi)
			authRoutes.PUT("/income_update", incomeAPI.UpdateIncomeDataApi)
			authRoutes.POST("/income_delete", incomeAPI.DeleteIncomeDataApi)
			// 他のエンドポイントのルーティングもここで設定
		}
	}
}
