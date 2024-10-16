// controllers/annual_income_management_controllers.go
package controllers

import (
	"net/http"
	"server/common"
	"server/config"
	"server/models" // モデルのインポート
	"server/utils"
	"server/validation"

	"github.com/gin-gonic/gin"
)

type (
	IncomeDataFetcher interface {
		GetIncomeDataInRangeApi(c *gin.Context)
		GetDateRangeApi(c *gin.Context)
		GetYearIncomeAndDeductionApi(c *gin.Context)
		InsertIncomeDataApi(c *gin.Context)
		UpdateIncomeDataApi(c *gin.Context)
		DeleteIncomeDataApi(c *gin.Context)
	}

	apiGetIncomeDataFetcher struct{}

	requestInsertIncomeData struct {
		Data []models.InsertIncomeData `json:"data"`
	}

	requestUpdateIncomeData struct {
		Data []models.UpdateIncomeData `json:"data"`
	}

	requestDeleteIncomeData struct {
		Data []models.DeleteIncomeData `json:"data"`
	}
)

func NewIncomeDataFetcher() IncomeDataFetcher {
	return &apiGetIncomeDataFetcher{}
}

// GetIncomeDataInRangeApi は登録された給料及び賞与の金額を指定期間で返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetIncomeDataInRangeApi(c *gin.Context) {
	var common common.CommonFetcher = common.NewCommonFetcher()
	// パラメータから日付の始まりと終わりを取得
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userIdPrams := c.Query("user_id")

	validator := validation.RequestYearIncomeAndDeductiontData{
		UserId:    userIdPrams,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	userId, _ := common.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	incomeData, err := dbFetcher.GetIncomeDataInRange(startDate, endDate, userId)

	if err != nil {
		response := utils.Response{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: incomeData,
	}
	c.JSON(http.StatusOK, response)
}

// GetDateRangeApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetDateRangeApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	var common common.CommonFetcher = common.NewCommonFetcher()
	userIdPrams := c.Query("user_id")

	validator := validation.RequestDateRangeData{
		UserId: userIdPrams,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	userId, _ := common.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	paymentDate, err := dbFetcher.GetDateRange(userId)

	if err != nil {
		response := utils.Response{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: paymentDate,
	}
	c.JSON(http.StatusOK, response)
}

// GetYearIncomeAndDeductionApi は各年ごとの収入、差引額、手取を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) GetYearIncomeAndDeductionApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	var common common.CommonFetcher = common.NewCommonFetcher()
	userIdPrams := c.Query("user_id")

	validator := validation.RequestDateRangeData{
		UserId: userIdPrams,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		c.JSON(http.StatusBadRequest, errMsgList)
		return
	}

	userId, _ := common.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	yearIncomeData, err := dbFetcher.GetYearsIncomeAndDeduction(userId)

	if err != nil {
		response := utils.Response{
			ErrorMsg: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: yearIncomeData,
	}
	c.JSON(http.StatusOK, response)
}

// InsertIncomeDataApi は新規登録
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) InsertIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestInsertIncomeData
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	for idx, data := range requestData.Data {
		validator := validation.RequestInsertIncomeData{
			PaymentDate:     data.PaymentDate,
			Age:             data.Age,
			Industry:        data.Industry,
			TotalAmount:     common.AnyToStr(data.TotalAmount),
			DeductionAmount: common.AnyToStr(data.DeductionAmount),
			TakeHomeAmount:  common.AnyToStr(data.TakeHomeAmount),
			Classification:  data.Classification,
			UserId:          data.UserID,
		}
		// TODO:エラーになったリクエストデータを全て出力するのか？
		// それとも、エラーが発生したレコードだけ出力するのか、考える
		if valid, errMsgList := validator.Validate(); !valid {
			errMsgList[0].RecodeRows = idx + 1
			c.JSON(http.StatusBadRequest, errMsgList)
			return
		}
	}

	// 収入データベースへ新しいデータ登録
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.InsertIncome(requestData.Data); err != nil {
		response := utils.Response{
			ErrorMsg: "新規登録時にエラーが発生。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: "新規給料情報を登録致しました。",
	}
	c.JSON(http.StatusOK, response)
}

// UpdateIncomeDataApi は更新
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) UpdateIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestUpdateIncomeData
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	for idx, data := range requestData.Data {
		validator := validation.RequestUpdateIncomeData{
			IncomeForecastID: data.IncomeForecastID,
			PaymentDate:      data.PaymentDate,
			Age:              data.Age,
			Industry:         data.Industry,
			TotalAmount:      common.AnyToStr(data.TotalAmount),
			DeductionAmount:  common.AnyToStr(data.DeductionAmount),
			TakeHomeAmount:   common.AnyToStr(data.TakeHomeAmount),
			UpdateUser:       data.UpdateUser,
			Classification:   data.Classification,
		}
		// TODO:エラーになったリクエストデータを全て出力するのか？
		// それとも、エラーが発生したレコードだけ出力するのか、考える
		if valid, errMsgList := validator.Validate(); !valid {
			errMsgList[0].RecodeRows = idx + 1
			c.JSON(http.StatusBadRequest, errMsgList)
			return
		}
	}

	// 収入データベースの更新
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.UpdateIncome(requestData.Data); err != nil {
		response := utils.Response{
			ErrorMsg: "更新時にエラーが発生。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: "給料情報の更新が問題なく成功しました。",
	}
	c.JSON(http.StatusOK, response)
}

// DeleteIncomeDataApi は削除
// 引数:
//   - c: Ginコンテキスト
//

func (af *apiGetIncomeDataFetcher) DeleteIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestDeleteIncomeData
	// JSONのバインドエラーチェックは無視する。後にバリデーションチェックを行うため
	c.ShouldBindJSON(&requestData)

	for idx, data := range requestData.Data {
		validator := validation.RequestDeleteIncomeData{
			IncomeForecastID: data.IncomeForecastID,
		}
		if valid, errMsgList := validator.Validate(); !valid {
			errMsgList[0].RecodeRows = idx + 1
			c.JSON(http.StatusBadRequest, errMsgList)
			return
		}
	}

	// 収入データベースの指定されたIDの削除
	dbFetcher, _, _ := models.NewPostgreSQLDataFetcher(config.DataSourceName)
	if err := dbFetcher.DeleteIncome(requestData.Data); err != nil {
		response := utils.Response{
			ErrorMsg: "削除中にエラーが発生しました",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.Response{
		Result: "給料情報の削除が問題なく成功しました。",
	}
	c.JSON(http.StatusOK, response)
}
