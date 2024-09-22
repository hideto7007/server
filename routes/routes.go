package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// APiインターフェイスのインスタンス定義
	var singInAPI controllers.SingInDataFetcher = controllers.NewSingInDataFetcher()
	var priceAPI controllers.PriceManagementFetcher = controllers.NewPriceManagementFetcher()
	var incomeAPI controllers.IncomeDataFetcher = controllers.NewIncomeDataFetcher()

	// ルートの設定
	Routes := r.Group("/api")
	{
		Routes.GET("/singin", singInAPI.GetSingInApi)
		Routes.GET("/price", priceAPI.GetPriceInfoApi)
		Routes.GET("/income_data", incomeAPI.GetIncomeDataInRangeApi)
		Routes.GET("/range_date", incomeAPI.GetDateRangeApi)
		Routes.GET("/years_income_date", incomeAPI.GetYearIncomeAndDeductionApi)
		Routes.POST("/income_create", incomeAPI.InsertIncomeDataApi)
		Routes.PUT("/income_update", incomeAPI.UpdateIncomeDataApi)
		Routes.DELETE("/income_delete", incomeAPI.DeleteIncomeDataApi)
		// 他のエンドポイントのルーティングもここで設定
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
