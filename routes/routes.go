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
		templates.NewEmailTemplateManager(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	var lineApi controllers.LineService = controllers.NewLineService(
		templates.NewEmailTemplateManager(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	var priceAPI controllers.PriceManagementFetcher = controllers.NewPriceManagementFetcher(
		common.NewCommonFetcher(),
	)
	var incomeAPI controllers.IncomeDataFetcher = controllers.NewIncomeDataFetcher(
		common.NewCommonFetcher(),
	)

	// google認証
	// r.GET("auth/google/signin", googleApi.GoogleSignIn)
	// r.GET("auth/google/signup", googleApi.GoogleSignUp)
	// r.GET("auth/google/delete", googleApi.GoogleDelete)
	// r.GET("auth/google/signin/callback", googleApi.GoogleSignInCallback)
	// r.GET("auth/google/signup/callback", googleApi.GoogleSignUpCallback)
	// r.GET("auth/google/delete/callback", googleApi.GoogleDeleteCallback)
	// line認証
	// r.GET("auth/line/signin", lineApi.LineSignIn)
	// r.GET("auth/line/signup", lineApi.LineSignUp)
	// r.GET("auth/line/delete", lineApi.LineDelete)
	// r.GET("auth/line/signin/callback", lineApi.LineSignInCallback)
	// r.GET("auth/line/signup/callback", lineApi.LineSignUpCallback)
	// r.GET("auth/line/delete/callback", lineApi.LineDeleteCallback)

	// ルートの設定
	Routes := r.Group("/api")
	{
		Routes.POST("/signin", signAPI.PostSignInApi)
		Routes.GET("/refresh_token", signAPI.GetRefreshTokenApi)
		Routes.POST("/temporary_signup", signAPI.TemporaryPostSignUpApi)
		Routes.GET("/retry_auth_email", signAPI.RetryAuthEmail)
		Routes.POST("/signup", signAPI.PostSignUpApi)
		Routes.PUT("/signin_edit/:user_id", signAPI.PutSignInEditApi)
		Routes.DELETE("/signin_delete/:user_id", signAPI.DeleteSignInApi) // 修正
		Routes.GET("/signout", signAPI.SignOutApi)
		Routes.GET("/register_email_check_notice", signAPI.RegisterEmailCheckNotice)
		// tokenIdからUserIdを取得していて、トークン漏洩防止のためパラメータにUserIdは含めない
		Routes.PUT("/new_password_update", signAPI.NewPasswordUpdate)
		// google認証
		Routes.GET("/google/signin/callback", googleApi.GoogleSignInCallback)
		Routes.GET("/google/signup/callback", googleApi.GoogleSignUpCallback)
		Routes.GET("/google/delete/callback", googleApi.GoogleDeleteCallback)
		// line認証
		Routes.GET("/line/signin/callback", lineApi.LineSignInCallback)
		Routes.GET("/line/signup/callback", lineApi.LineSignUpCallback)
		Routes.GET("/line/delete/callback", lineApi.LineDeleteCallback)

		// 認証が必要なルートにミドルウェアを追加
		authRoutes := Routes.Group("/")
		authRoutes.Use(middleware.JWTAuthMiddleware(utils.UtilsDataFetcher{}))
		{
			authRoutes.GET("/price", priceAPI.GetPriceInfoApi)
			authRoutes.GET("/income_data", incomeAPI.GetIncomeDataInRangeApi)
			authRoutes.GET("/range_date", incomeAPI.GetDateRangeApi)
			authRoutes.GET("/years_income_date", incomeAPI.GetYearIncomeAndDeductionApi)
			authRoutes.POST("/income_create", incomeAPI.InsertIncomeDataApi)
			// データが複数件の場合があるため、urlにキーは付与しない
			authRoutes.PUT("/income_update", incomeAPI.UpdateIncomeDataApi)
			authRoutes.POST("/income_delete", incomeAPI.DeleteIncomeDataApi)
			// 他のエンドポイントのルーティングもここで設定
		}
	}
}
