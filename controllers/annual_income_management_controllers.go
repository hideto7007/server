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

	requestInsertIncomeData struct {
		Data []models.InsertIncomeData `json:"data"`
	}

	requestUpdateIncomeData struct {
		Data []models.UpdateIncomeData `json:"data"`
	}

	requestDeleteIncomeData struct {
		Data []models.DeleteIncomeData `json:"data"`
	}

	apiIncomeDataFetcher struct {
		CommonFetcher common.CommonFetcher
	}
)

func NewIncomeDataFetcher(CommonFetcher common.CommonFetcher) IncomeDataFetcher {
	return &apiIncomeDataFetcher{
		CommonFetcher: CommonFetcher,
	}
}

// GetIncomeDataInRangeApi は登録された給料及び賞与の金額を指定期間で返すAPI
//
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) GetIncomeDataInRangeApi(c *gin.Context) {
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
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId, _ := aid.CommonFetcher.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	incomeData, err := dbFetcher.GetIncomeDataInRange(startDate, endDate, userId)

	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[[]models.IncomeData]{
		Result: incomeData,
	}
	c.JSON(http.StatusOK, response)
}

// GetDateRangeApi は登録されている最も古い日付と最も新しい日付を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) GetDateRangeApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userIdPrams := c.Query("user_id")

	validator := validation.RequestDateRangeData{
		UserId: userIdPrams,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId, _ := aid.CommonFetcher.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	paymentDate, err := dbFetcher.GetDateRange(userId)

	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[[]models.PaymentDate]{
		Result: paymentDate,
	}
	c.JSON(http.StatusOK, response)
}

// GetYearIncomeAndDeductionApi は各年ごとの収入、差引額、手取を取得するAPI
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) GetYearIncomeAndDeductionApi(c *gin.Context) {
	// パラメータからユーザー情報取得
	userIdPrams := c.Query("user_id")

	validator := validation.RequestDateRangeData{
		UserId: userIdPrams,
	}

	if valid, errMsgList := validator.Validate(); !valid {
		response := utils.ErrorValidationResponse{
			Result: errMsgList,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userId, _ := aid.CommonFetcher.StrToInt(userIdPrams)

	// データベースから指定範囲のデータを取得
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	yearIncomeData, err := dbFetcher.GetYearsIncomeAndDeduction(userId)

	if err != nil {
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[[]models.YearsIncomeData]{
		Result: yearIncomeData,
	}
	c.JSON(http.StatusOK, response)
}

// InsertIncomeDataApi は新規登録
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) InsertIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestInsertIncomeData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(requestData.Data) == 0 {
		response := utils.ErrorMessageResponse{
			Result: "登録するデータが存在しません。",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	for idx, data := range requestData.Data {
		validator := validation.RequestInsertIncomeData{
			PaymentDate:     data.PaymentDate,
			Age:             data.Age,
			Industry:        data.Industry,
			TotalAmount:     common.AnyToStr(data.TotalAmount),
			DeductionAmount: common.AnyToStr(data.DeductionAmount),
			TakeHomeAmount:  common.AnyToStr(data.TakeHomeAmount),
			Classification:  data.Classification,
			UserId:          common.AnyToStr(data.UserID),
		}
		// TODO:エラーになったリクエストデータを全て出力するのか？
		// それとも、エラーが発生したレコードだけ出力するのか、考える
		if valid, errMsgList := validator.Validate(); !valid {
			response := utils.ErrorValidationResponse{
				RecodeRows: idx + 1,
				Result:     errMsgList,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// 収入データベースへ新しいデータ登録
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	if err := dbFetcher.InsertIncome(requestData.Data); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "新規登録時にエラーが発生。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[string]{
		Result: "新規給料情報を登録致しました。",
	}
	c.JSON(http.StatusOK, response)
}

// UpdateIncomeDataApi は更新
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) UpdateIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestUpdateIncomeData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(requestData.Data) == 0 {
		response := utils.ErrorMessageResponse{
			Result: "更新するデータが存在しません。",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

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
			response := utils.ErrorValidationResponse{
				RecodeRows: idx + 1,
				Result:     errMsgList,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// 収入データベースの更新
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	if err := dbFetcher.UpdateIncome(requestData.Data); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "更新時にエラーが発生。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[string]{
		Result: "給料情報の更新が問題なく成功しました。",
	}
	c.JSON(http.StatusOK, response)
}

// DeleteIncomeDataApi は削除
// 引数:
//   - c: Ginコンテキスト
//

func (aid *apiIncomeDataFetcher) DeleteIncomeDataApi(c *gin.Context) {
	// JSONデータを受け取るための構造体を定義
	var requestData requestDeleteIncomeData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// エラーメッセージを出力して確認
		response := utils.ErrorMessageResponse{
			Result: err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(requestData.Data) == 0 {
		response := utils.ErrorMessageResponse{
			Result: "削除するデータが存在しません。",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	for idx, data := range requestData.Data {
		validator := validation.RequestDeleteIncomeData{
			IncomeForecastID: data.IncomeForecastID,
		}
		if valid, errMsgList := validator.Validate(); !valid {
			response := utils.ErrorValidationResponse{
				RecodeRows: idx + 1,
				Result:     errMsgList,
			}
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// 収入データベースの指定されたIDの削除
	dbFetcher, _, _ := models.NewAnnualIncomeDataFetcher(config.GetDataBaseSource())
	if err := dbFetcher.DeleteIncome(requestData.Data); err != nil {
		response := utils.ErrorMessageResponse{
			Result: "削除中にエラーが発生しました。",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// JSONレスポンスを返す
	response := utils.ResponseData[string]{
		Result: "給料情報の削除が問題なく成功しました。",
	}
	c.JSON(http.StatusOK, response)
}
