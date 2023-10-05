// models/price_management_controllers.go
package models

import (
	"database/sql"
	"log"
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
	TakeHomeAmount   sql.NullInt64
	Classification   string
	UserID           int
}

// GetIncomeDataInRange はDBに登録された給料及び賞与の金額を指定期間で返す。
//
// 引数:
//   - startDate: 始まりの期間
//   - endDate  : 終わりの期間
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
