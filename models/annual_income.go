// models/anuual_income.go
package models

import (
	"database/sql"
	"log"
	"server/common"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	AnuualIncomeFetcher interface {
		GetIncomeDataInRange(startDate, endDate string) ([]IncomeData, error)
		GetStartDataAndEndDate(UserID string) ([]PaymentDate, error)
		GetYearsIncomeAndDeduction(UserID string) ([]YearsIncomeData, error)
		InsertIncome(data []InsertIncomeData) error
		UpdateIncome(data []UpdateIncomeData) error
		DeleteIncome(UserID []DeleteIncomeData) error
	}

	IncomeData struct {
		IncomeForecastID uuid.UUID
		PaymentDate      time.Time
		Age              string
		Industry         string
		TotalAmount      int
		DeductionAmount  int
		TakeHomeAmount   int
		Classification   string
		UserID           int
	}

	PaymentDate struct {
		UserID            int
		StratPaymaentDate string
		EndPaymaentDate   string
	}

	YearsIncomeData struct {
		Years           string
		TotalAmount     int
		DeductionAmount int
		TakeHomeAmount  int
	}

	InsertIncomeData struct {
		PaymentDate     string `json:"payment_date"`
		Age             string `json:"age"`
		Industry        string `json:"industry"`
		TotalAmount     string `json:"total_amount"`
		DeductionAmount string `json:"deduction_amount"`
		TakeHomeAmount  string `json:"take_home_amount"`
		UpdateUser      string `json:"update_user"`
		Classification  string `json:"classification"`
		UserID          string `json:"user_id"`
	}

	UpdateIncomeData struct {
		IncomeForecastID string `json:"income_forecast_id"`
		PaymentDate      string `json:"payment_date"`
		Age              string `json:"age"`
		Industry         string `json:"industry"`
		TotalAmount      string `json:"total_amount"`
		DeductionAmount  string `json:"deduction_amount"`
		TakeHomeAmount   string `json:"take_home_amount"`
		Classification   string `json:"classification"`
	}

	DeleteIncomeData struct {
		IncomeForecastID string `json:"income_forecast_id"`
	}

	PostgreSQLDataFetcher struct{ db *sql.DB }
)

func NewPostgreSQLDataFetcher(dataSourceName string) *PostgreSQLDataFetcher {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("sql.Open error %s", err)
	}
	return &PostgreSQLDataFetcher{db: db}
}

// GetIncomeDataInRange はDBに登録された給料及び賞与の金額を指定期間で返す。
//
// 引数:
//   - StratPaymaentDate: 始まりの期間
//   - EndPaymaentDate: 終わりの期間
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) GetIncomeDataInRange(startDate, endDate string) ([]IncomeData, error) {
	var incomeData []IncomeData

	// startDate と endDate を日付型に変換
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// データベースクエリを実行
	rows, err := pf.db.Query(`
        SELECT income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, classification, user_id
        FROM incomeforecast_incomeforecastdata
        WHERE payment_date BETWEEN $1 AND $2
    `, start, end)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data IncomeData
		err := rows.Scan(
			&data.IncomeForecastID,
			&data.PaymentDate,
			&data.Age,
			&data.Industry,
			&data.TotalAmount,
			&data.DeductionAmount,
			&data.TakeHomeAmount,
			&data.Classification,
			&data.UserID,
		)

		if err != nil {
			return nil, err
		}

		incomeData = append(incomeData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incomeData, nil
}

// GetStartDataAndEndDate は対象ユーザーの情報で最も古い日付と最も新しい日付を取得して返す。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) GetStartDataAndEndDate(UserId string) ([]PaymentDate, error) {
	var paymentDate []PaymentDate

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	rows, err := pf.db.Query(`
		SELECT user_id, MIN(payment_date) as "start_paymaent_date", MAX(payment_date) as "end_paymaent_date" from incomeforecast_incomeforecastdata
		WHERE user_id = $1
		GROUP BY user_id;
    `, UserId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// rows.Scanがデータを変数に直接書き込むため
	for rows.Next() {
		var (
			userId            int
			stratPaymaentDate time.Time
			endPaymaentDate   time.Time
		)
		err := rows.Scan(
			&userId,
			&stratPaymaentDate,
			&endPaymaentDate,
		)

		if err != nil {
			return nil, err
		}

		// stratPaymaentDate および endPaymaentDate を文字列に変換
		var common common.CommonFetcher = common.NewCommonFetcher()
		startDateStr := common.TimeToStr(stratPaymaentDate)
		endDateStr := common.TimeToStr(endPaymaentDate)

		// 変換したデータをPaymentDate構造体にセットする
		replaceData := PaymentDate{
			UserID:            userId,
			StratPaymaentDate: startDateStr,
			EndPaymaentDate:   endDateStr,
		}

		paymentDate = append(paymentDate, replaceData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paymentDate, nil
}

// GetYearsIncomeAndDeduction は対象ユーザー情報の各年ごとの収入、差引額、手取を取得して返す。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) GetYearsIncomeAndDeduction(UserId string) ([]YearsIncomeData, error) {
	var yearsIncomeData []YearsIncomeData

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	rows, err := pf.db.Query(`
		SELECT 
			TO_CHAR(payment_date, 'YYYY') as "year" ,
			SUM(total_amount) as "sum_total_amount", 
			SUM(deduction_amount) as "sum_deduction_amount",  
			SUM(take_home_amount) as "sum_take_home_amount"
		FROM incomeforecast_incomeforecastdata
		WHERE user_id = $1
		GROUP BY TO_CHAR(payment_date, 'YYYY')
		ORDER BY TO_CHAR(payment_date, 'YYYY') asc;
    `, UserId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// rows.Scanがデータを変数に直接書き込むため
	for rows.Next() {
		var data YearsIncomeData
		err := rows.Scan(
			&data.Years,
			&data.TotalAmount,
			&data.DeductionAmount,
			&data.TakeHomeAmount,
		)

		if err != nil {
			return nil, err
		}

		yearsIncomeData = append(yearsIncomeData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return yearsIncomeData, nil
}

// InsertIncome は新規登録
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) InsertIncome(data []InsertIncomeData) error {

	defer pf.db.Close()
	var err error
	var common common.CommonFetcher = common.NewCommonFetcher()
	deleteFlag := 0
	createdAt := time.Now()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertStatement := `
        INSERT INTO public.incomeforecast_incomeforecastdata
        (income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, delete_flag, update_user, created_at, classification, user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	for _, insertData := range data {
		data := InsertIncomeData{
			PaymentDate:     insertData.PaymentDate,
			Age:             insertData.Age,
			Industry:        insertData.Industry,
			TotalAmount:     insertData.TotalAmount,
			DeductionAmount: insertData.DeductionAmount,
			TakeHomeAmount:  insertData.TakeHomeAmount,
			UpdateUser:      insertData.UpdateUser,
			Classification:  insertData.Classification,
			UserID:          insertData.UserID,
		}
		uuid := uuid.New().String()
		ageToInt, _ := common.StrToInt(data.Age)
		totalAmountToInt, _ := common.StrToInt(data.TotalAmount)
		deductionAmountToInt, _ := common.StrToInt(data.DeductionAmount)
		takeHomeAmountToInt, _ := common.StrToInt(data.TakeHomeAmount)
		userIdToInt, _ := common.StrToInt(data.UserID)
		if _, err = tx.Exec(insertStatement,
			uuid,
			data.PaymentDate,
			ageToInt,
			data.Industry,
			totalAmountToInt,
			deductionAmountToInt,
			takeHomeAmountToInt,
			deleteFlag,
			data.UpdateUser,
			createdAt,
			data.Classification,
			userIdToInt); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// トランザクションをコミット
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

// UpdateIncome は更新
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) UpdateIncome(data []UpdateIncomeData) error {

	defer pf.db.Close()
	var err error
	var common common.CommonFetcher = common.NewCommonFetcher()
	createdAt := time.Now()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	updateStatement := `
        UPDATE public.incomeforecast_incomeforecastdata
        SET 
			payment_date = $1, 
			age = $2, 
			industry = $3, 
			total_amount = $4, 
			deduction_amount = $5, 
			take_home_amount = $6, 
			created_at = $7, 
			classification = $8
        WHERE income_forecast_id = $9;
    `

	for _, insertData := range data {
		data := UpdateIncomeData{
			IncomeForecastID: insertData.IncomeForecastID,
			PaymentDate:      insertData.PaymentDate,
			Age:              insertData.Age,
			Industry:         insertData.Industry,
			TotalAmount:      insertData.TotalAmount,
			DeductionAmount:  insertData.DeductionAmount,
			TakeHomeAmount:   insertData.TakeHomeAmount,
			Classification:   insertData.Classification,
		}
		ageToInt, _ := common.StrToInt(data.Age)
		totalAmountToInt, _ := common.StrToInt(data.TotalAmount)
		deductionAmountToInt, _ := common.StrToInt(data.DeductionAmount)
		takeHomeAmountToInt, _ := common.StrToInt(data.TakeHomeAmount)
		if _, err = tx.Exec(updateStatement,
			data.PaymentDate,
			ageToInt,
			data.Industry,
			totalAmountToInt,
			deductionAmountToInt,
			takeHomeAmountToInt,
			createdAt,
			data.Classification,
			data.IncomeForecastID); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// トランザクションをコミット
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// DeleteIncome は削除
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) DeleteIncome(data []DeleteIncomeData) error {

	defer pf.db.Close()
	var err error

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	deleteStatement := `
        DELETE FROM public.incomeforecast_incomeforecastdata
        WHERE income_forecast_id = $1;
    `

	for _, deleteData := range data {
		if _, err = tx.Exec(deleteStatement, deleteData.IncomeForecastID); err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	// トランザクションをコミット
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
