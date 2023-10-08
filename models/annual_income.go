// models/price_management_controllers.go
package models

import (
	"database/sql"
	"log"
	"server/common"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// IncomeDataの構造体
type IncomeData struct {
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

// PaymentDateの構造体
type PaymentDate struct {
	UserID            int
	StratPaymaentDate string
	EndPaymaentDate   string
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

func GetIncomeDataInRange(startDate, endDate string) ([]IncomeData, error) {
	var incomeData []IncomeData

	dataSourceName := "user=postgres dbname=postgres password=pedev7007 host=localhost port=5432  sslmode=disable"

	// startDate と endDate を日付型に変換
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("sql.Open error %s", err)
	}

	// データベースクエリを実行
	rows, err := db.Query(`
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

func GetStartDataAndEndDate(UserId string) ([]PaymentDate, error) {
	var paymentDate []PaymentDate

	dataSourceName := "user=postgres dbname=postgres password=pedev7007 host=localhost port=5432  sslmode=disable"
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("sql.Open error %s", err)
	}

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	rows, err := db.Query(`
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
