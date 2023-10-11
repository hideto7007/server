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
		CreateIncomeData(UserID string) ([]IncomeData, error)
		UpdateIncomeData(UserID string) ([]IncomeData, error)
		DeleteIncomeData(UserID string) ([]IncomeData, error)
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

// CreateIncomeData は新規作成。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) CreateIncomeData(UserId string) ([]IncomeData, error) {
	var incomeData []IncomeData

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	// rows, err := pf.db.Query(`
	// 	SELECT
	// 		TO_CHAR(payment_date, 'YYYY') as "year" ,
	// 		SUM(total_amount) as "sum_total_amount",
	// 		SUM(deduction_amount) as "sum_deduction_amount",
	// 		SUM(take_home_amount) as "sum_take_home_amount"
	// 	FROM incomeforecast_incomeforecastdata
	// 	WHERE user_id = $1
	// 	GROUP BY TO_CHAR(payment_date, 'YYYY')
	// 	ORDER BY TO_CHAR(payment_date, 'YYYY') asc;
	// `, UserId)

	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// // ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// // rows.Scanがデータを変数に直接書き込むため
	// for rows.Next() {
	// 	var data YearsIncomeData
	// 	err := rows.Scan(
	// 		&data.labels,
	// 		&data.totalAmount,
	// 		&data.deductionAmount,
	// 		&data.takeHomeAmount,
	// 	)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	yearsIncomeData = append(yearsIncomeData, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return incomeData, nil
}

// UpdateIncomeData は更新。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) UpdateIncomeData(UserId string) ([]IncomeData, error) {
	var incomeData []IncomeData

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	// rows, err := pf.db.Query(`
	// 	SELECT
	// 		TO_CHAR(payment_date, 'YYYY') as "year" ,
	// 		SUM(total_amount) as "sum_total_amount",
	// 		SUM(deduction_amount) as "sum_deduction_amount",
	// 		SUM(take_home_amount) as "sum_take_home_amount"
	// 	FROM incomeforecast_incomeforecastdata
	// 	WHERE user_id = $1
	// 	GROUP BY TO_CHAR(payment_date, 'YYYY')
	// 	ORDER BY TO_CHAR(payment_date, 'YYYY') asc;
	// `, UserId)

	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// // ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// // rows.Scanがデータを変数に直接書き込むため
	// for rows.Next() {
	// 	var data YearsIncomeData
	// 	err := rows.Scan(
	// 		&data.labels,
	// 		&data.totalAmount,
	// 		&data.deductionAmount,
	// 		&data.takeHomeAmount,
	// 	)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	yearsIncomeData = append(yearsIncomeData, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return incomeData, nil
}

// DeleteIncomeData は削除。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *PostgreSQLDataFetcher) DeleteIncomeData(UserId string) ([]IncomeData, error) {
	var incomeData []IncomeData

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	// rows, err := pf.db.Query(`
	// 	SELECT
	// 		TO_CHAR(payment_date, 'YYYY') as "year" ,
	// 		SUM(total_amount) as "sum_total_amount",
	// 		SUM(deduction_amount) as "sum_deduction_amount",
	// 		SUM(take_home_amount) as "sum_take_home_amount"
	// 	FROM incomeforecast_incomeforecastdata
	// 	WHERE user_id = $1
	// 	GROUP BY TO_CHAR(payment_date, 'YYYY')
	// 	ORDER BY TO_CHAR(payment_date, 'YYYY') asc;
	// `, UserId)

	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// // ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// // rows.Scanがデータを変数に直接書き込むため
	// for rows.Next() {
	// 	var data YearsIncomeData
	// 	err := rows.Scan(
	// 		&data.labels,
	// 		&data.totalAmount,
	// 		&data.deductionAmount,
	// 		&data.takeHomeAmount,
	// 	)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	yearsIncomeData = append(yearsIncomeData, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return incomeData, nil
}
