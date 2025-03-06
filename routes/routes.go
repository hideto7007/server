package routes

import (
	"server/common"
	"server/config"
	"server/controllers"
	controllers_common "server/controllers/common"
	"server/middleware"
	"server/templates"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	httpClient := common.NewHTTPClient()

	// APiインターフェイスのインスタンス定義
	var signAPI controllers.SignDataFetcher = controllers.NewSignDataFetcher(
		utils.NewUtilsFetcher(utils.JwtSecret),
		common.NewCommonFetcher(),
		templates.NewEmailTemplateManager(),
		config.NewRedisManager(),
	)
	var googleApi controllers.GoogleService = controllers.NewGoogleService(
		config.NewGoogleManager(),
		controllers_common.NewControllersCommonManager(
			config.NewGoogleManager(),
			config.NewLineManager(httpClient),
		),
		templates.NewEmailTemplateManager(),
		utils.NewUtilsFetcher(utils.JwtSecret),
	)
	var lineApi controllers.LineService = controllers.NewLineService(
		config.NewLineManager(httpClient),
		controllers_common.NewControllersCommonManager(
			config.NewGoogleManager(),
			config.NewLineManager(httpClient),
		),
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
	r.GET("auth/google/signin", googleApi.GoogleSignIn)
	r.GET("auth/google/signup", googleApi.GoogleSignUp)
	r.GET("auth/google/delete", googleApi.GoogleDelete)
	r.GET("auth/google/signin/callback", googleApi.GoogleSignInCallback)
	r.GET("auth/google/signup/callback", googleApi.GoogleSignUpCallback)
	r.GET("auth/google/delete/callback", googleApi.GoogleDeleteCallback)
	// line認証
	r.GET("auth/line/signin", lineApi.LineSignIn)
	r.GET("auth/line/signup", lineApi.LineSignUp)
	r.GET("auth/line/delete", lineApi.LineDelete)
	r.GET("auth/line/signin/callback", lineApi.LineSignInCallback)
	r.GET("auth/line/signup/callback", lineApi.LineSignUpCallback)
	r.GET("auth/line/delete/callback", lineApi.LineDeleteCallback)

	// ルートの設定
	Routes := r.Group("/api")
	{
		Routes.POST("/signin", signAPI.PostSignInApi)
		Routes.GET("/refresh_token", signAPI.GetRefreshTokenApi)
		Routes.POST("/temporary_signup", signAPI.TemporaryPostSignUpApi)
		Routes.GET("/retry_auth_email", signAPI.RetryAuthEmail)
		Routes.POST("/signup", signAPI.PostSignUpApi)
		Routes.PUT("/signin_edit/:user_id", signAPI.PutSignInEditApi)
		Routes.DELETE("/signin_delete", signAPI.DeleteSignInApi)
		Routes.GET("/signout", signAPI.SignOutApi)
		Routes.GET("/register_email_check_notice", signAPI.RegisterEmailCheckNotice)
		// tokenIdからUserIdを取得していて、トークン漏洩防止のためパラメータにUserIdは含めない
		Routes.PUT("/new_password_update", signAPI.NewPasswordUpdate)

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
