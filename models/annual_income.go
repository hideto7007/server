// models/anuual_income.go
package models

import (
	"database/sql"
	"fmt"
	"log"
	"server/DB"
	"server/common"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type (
	AnuualIncomeFetcher interface {
		GetIncomeDataInRange(StartDate, EndDate, UserId string) ([]IncomeData, error)
		GetDateRange(UserId int) ([]PaymentDate, error)
		GetYearsIncomeAndDeduction(UserId int) ([]YearsIncomeData, error)
		InsertIncome(data []InsertIncomeData) error
		UpdateIncome(data []UpdateIncomeData) error
		DeleteIncome(data []DeleteIncomeData) error
	}

	IncomeData struct {
		IncomeForecastID uuid.UUID `json:"income_forecast_id"`
		PaymentDate      time.Time `json:"payment_date"`
		Age              string    `json:"age"`
		Industry         string    `json:"industry"`
		TotalAmount      int       `json:"total_amount"`
		DeductionAmount  int       `json:"deduction_amount"`
		TakeHomeAmount   int       `json:"take_home_amount"`
		Classification   string    `json:"classification"`
		UserID           int       `json:"user_id"`
	}

	PaymentDate struct {
		UserID           int    `json:"user_id"`
		StratPaymentDate string `json:"start_payment_date"`
		EndPaymentDate   string `json:"end_payment_date"`
	}

	YearsIncomeData struct {
		Years           string `json:"years"`
		TotalAmount     int    `json:"total_amount"`
		DeductionAmount int    `json:"deduction_amount"`
		TakeHomeAmount  int    `json:"take_home_amount"`
	}

	InsertIncomeData struct {
		PaymentDate     string      `json:"payment_date"`
		Age             int         `json:"age"`
		Industry        string      `json:"industry"`
		TotalAmount     interface{} `json:"total_amount"`
		DeductionAmount interface{} `json:"deduction_amount"`
		TakeHomeAmount  interface{} `json:"take_home_amount"`
		Classification  string      `json:"classification"`
		UserID          interface{} `json:"user_id"`
	}

	UpdateIncomeData struct {
		IncomeForecastID string      `json:"income_forecast_id"`
		PaymentDate      string      `json:"payment_date"`
		Age              int         `json:"age"`
		Industry         string      `json:"industry"`
		TotalAmount      interface{} `json:"total_amount"`
		DeductionAmount  interface{} `json:"deduction_amount"`
		TakeHomeAmount   interface{} `json:"take_home_amount"`
		UpdateUser       string      `json:"update_user"`
		Classification   string      `json:"classification"`
	}

	DeleteIncomeData struct {
		IncomeForecastID string `json:"income_forecast_id"`
	}

	AnnualIncomeDataFetcher struct{ db *sql.DB }
)

func NewAnnualIncomeDataFetcher(dataSourceName string) (*AnnualIncomeDataFetcher, sqlmock.Sqlmock, error) {
	if dataSourceName == "test" {
		db, mock, err := sqlmock.New()
		return &AnnualIncomeDataFetcher{db: db}, mock, err
	} else {
		// test実行時に以下のカバレッジは無視する
		db, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("sql.Open error %s", err)
		}
		return &AnnualIncomeDataFetcher{db: db}, nil, nil
	}
}

// GetIncomeDataInRange はDBに登録された給料及び賞与の金額を指定期間で返す。
//
// 引数:
//   - StratPaymentDate: 始まりの期間
//   - EndPaymentDate: 終わりの期間
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *AnnualIncomeDataFetcher) GetIncomeDataInRange(StartDate string, EndDate string, UserId int) ([]IncomeData, error) {
	var incomeData []IncomeData

	// startDate と endDate を日付型に変換
	start, err := time.Parse("2006-01-02", StartDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", EndDate)
	if err != nil {
		return nil, err
	}

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetIncomeDataInRangeSyntax, start, end, UserId)

	if err != nil {
		return nil, fmt.Errorf("クエリー実行エラー： %v", err)
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

// GetDateRange は対象ユーザーの情報で最も古い日付と最も新しい日付を取得して返す。
//
// 引数:
//   - UserId: ユーザーID
//
// 戻り値:
//
//	戻り値1: 取得したDBの構造体
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *AnnualIncomeDataFetcher) GetDateRange(UserId int) ([]PaymentDate, error) {
	var paymentDate []PaymentDate

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	rows, err := pf.db.Query(DB.GetDateRangeSyntax, UserId)

	if err != nil {
		return nil, fmt.Errorf("クエリー実行エラー： %v", err)
	}
	defer rows.Close()

	// ユーザーidと日付は別々の型で受け取り、各変数のポインターに渡す
	// rows.Scanがデータを変数に直接書き込むため
	for rows.Next() {
		var (
			userId           int
			stratPaymentDate time.Time
			endPaymentDate   time.Time
		)
		err := rows.Scan(
			&userId,
			&stratPaymentDate,
			&endPaymentDate,
		)

		if err != nil {
			return nil, err
		}

		// stratPaymentDate および endPaymentDate を文字列に変換
		var common common.CommonFetcher = common.NewCommonFetcher()
		startDateStr := common.TimeToStr(stratPaymentDate)
		endDateStr := common.TimeToStr(endPaymentDate)

		// 変換したデータをPaymentDate構造体にセットする
		replaceData := PaymentDate{
			UserID:           userId,
			StratPaymentDate: startDateStr,
			EndPaymentDate:   endDateStr,
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

func (pf *AnnualIncomeDataFetcher) GetYearsIncomeAndDeduction(UserId int) ([]YearsIncomeData, error) {
	var yearsIncomeData []YearsIncomeData

	// データベースクエリを実行
	// 集計関数で値を取得する際は、必ずカラム名を指定する
	rows, err := pf.db.Query(DB.GetYearsIncomeAndDeductionSyntax, UserId)

	if err != nil {
		return nil, fmt.Errorf("クエリー実行エラー： %v", err)
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

func (pf *AnnualIncomeDataFetcher) InsertIncome(data []InsertIncomeData) error {

	var err error
	createdAt := time.Now()

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// ロールバックをデフォルトに設定
	rollback := true
	defer func() {
		if rollback {
			tx.Rollback()
		}
	}()

	insertStatement := DB.InsertIncomeSyntax

	for _, insertData := range data {
		data := InsertIncomeData{
			PaymentDate:     insertData.PaymentDate,
			Age:             insertData.Age,
			Industry:        insertData.Industry,
			TotalAmount:     insertData.TotalAmount,
			DeductionAmount: insertData.DeductionAmount,
			TakeHomeAmount:  insertData.TakeHomeAmount,
			Classification:  insertData.Classification,
			UserID:          insertData.UserID,
		}
		uuid := uuid.New().String()
		if _, err = tx.Exec(insertStatement,
			uuid,
			data.PaymentDate,
			data.Age,
			data.Industry,
			data.TotalAmount,
			data.DeductionAmount,
			data.TakeHomeAmount,
			createdAt,
			data.Classification,
			data.UserID); err != nil {
			return err
		}
	}

	// コミット処理
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %v", err)
	}

	rollback = false

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

func (pf *AnnualIncomeDataFetcher) UpdateIncome(data []UpdateIncomeData) error {

	var err error
	createdAt := time.Now()

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// ロールバックをデフォルトに設定
	rollback := true
	defer func() {
		if rollback {
			tx.Rollback()
		}
	}()

	updateStatement := DB.UpdateIncomeSyntax

	for _, updateData := range data {
		data := UpdateIncomeData{
			IncomeForecastID: updateData.IncomeForecastID,
			PaymentDate:      updateData.PaymentDate,
			Age:              updateData.Age,
			Industry:         updateData.Industry,
			TotalAmount:      updateData.TotalAmount,
			DeductionAmount:  updateData.DeductionAmount,
			TakeHomeAmount:   updateData.TakeHomeAmount,
			UpdateUser:       updateData.UpdateUser,
			Classification:   updateData.Classification,
		}
		if _, err = tx.Exec(updateStatement,
			data.PaymentDate,
			data.Age,
			data.Industry,
			data.TotalAmount,
			data.DeductionAmount,
			data.TakeHomeAmount,
			createdAt,
			data.UpdateUser,
			data.Classification,
			data.IncomeForecastID); err != nil {
			return err
		}
	}

	// コミット処理
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %v", err)
	}

	rollback = false

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

func (pf *AnnualIncomeDataFetcher) DeleteIncome(data []DeleteIncomeData) error {

	var err error

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// ロールバックをデフォルトに設定
	rollback := true
	defer func() {
		if rollback {
			tx.Rollback()
		}
	}()

	deleteStatement := DB.DeleteIncomeSyntax

	for _, deleteData := range data {
		if _, err = tx.Exec(deleteStatement, deleteData.IncomeForecastID); err != nil {
			return err
		}
	}

	// コミット処理
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %v", err)
	}

	rollback = false

	return nil
}
