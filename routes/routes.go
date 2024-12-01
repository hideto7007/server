package routes

import (
	"server/common"
	"server/controllers"
	"server/middleware"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// APiインターフェイスのインスタンス定義
	var signAPI controllers.SignDataFetcher = controllers.NewSignDataFetcher(
		utils.NewUtilsFetcher(utils.JwtSecret),
		common.NewCommonFetcher(),
	)
	var priceAPI controllers.PriceManagementFetcher = controllers.NewPriceManagementFetcher(
		common.NewCommonFetcher(),
	)
	var incomeAPI controllers.IncomeDataFetcher = controllers.NewIncomeDataFetcher(
		common.NewCommonFetcher(),
	)

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

// {
//     "data": [
//         {
//             "income_forecast_id": "7b941edb-b7a2-e1e7-6466-ce53d1c8bcff",
//             "payment_date": "2024-02-10",
//             "age": 30,
//             "industry": "test開発",
//             "total_amount": 320524,
//             "deduction_amount": 93480,
//             "take_home_amount": 227044,
//             "classification": "給料",
//             "user_id": 1
//         },
//         {
//             "income_forecast_id": "af16418d-85d2-7945-bef8-bc50d3adbf82",
//             "payment_date": "2024-03-10",
//             "age": 30,
//             "industry": "test開発",
//             "total_amount": 320524,
//             "deduction_amount": 93480,
//             "take_home_amount": 227044,
//             "classification": "給料",
//             "user_id": 1
//         },
//         {
//             "income_forecast_id": "2c33ff50-d48a-094b-cc6a-bafa618dd96d",
//             "payment_date": "2024-04-10",
//             "age": 30,
//             "industry": "test開発",
//             "total_amount": 320524,
//             "deduction_amount": 93480,
//             "take_home_amount": 227044,
//             "classification": "給料",
//             "user_id": 1
//         }
//     ]
// }
