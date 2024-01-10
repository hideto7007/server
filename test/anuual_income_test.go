package common

import (
	"errors"
	"regexp"
	"server/DB"
	"server/config"
	"server/models"
	"testing"
	"time"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type (
	MockIncomeData struct {
		IncomeForecastID uuid.UUID
		PaymentDate      string
		Age              string
		Industry         string
		TotalAmount      int
		DeductionAmount  int
		TakeHomeAmount   int
		Classification   string
		UserID           int
	}
)

// func NewMocks
func TestGetIncomeDataInRange(t *testing.T) {
	t.Run("success TestGetIncomeDataInRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		StartDate := "2022-11-01"
		EndDate := "2022-12-30"
		UserId := "1"

		expectedData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2022, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
		}

		mockData := []models.IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2021, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2021, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2022, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2022, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("3d9752bd-0e2b-9994-7b90-55ecfd2105b5"),
				PaymentDate:      time.Date(2023, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
			{
				IncomeForecastID: uuid.MustParse("5ce422da-8989-bb7f-1e56-32db74aaa4ac"),
				PaymentDate:      time.Date(2023, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           1,
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"income_forecast_id", "payment_date", "age", "industry", "total_amount",
			"deduction_amount", "take_home_amount", "classification", "user_id",
		})

		start, err := time.Parse("2006-01-02", StartDate)
		if err != nil {
			return
		}

		end, err := time.Parse("2006-01-02", EndDate)
		if err != nil {
			return
		}

		for _, data := range mockData {

			// 指定期間内のデータのみを rows に追加
			if data.PaymentDate.After(start) && data.PaymentDate.Before(end) {
				rows.AddRow(
					data.IncomeForecastID.String(),
					data.PaymentDate,
					data.Age,
					data.Industry,
					data.TotalAmount,
					data.DeductionAmount,
					data.TakeHomeAmount,
					data.Classification,
					data.UserID,
				)
			}
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetIncomeDataInRangeSyntax)).
			WithArgs(start, end, UserId).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetIncomeDataInRange(StartDate, EndDate, UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		t.Log("result : ", result)
		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
	t.Run("error TestGetIncomeDataInRange1", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")

		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		EndDate := "2022-12-30"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange(StartDate, EndDate, "invalid")

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
	t.Run("error TestGetIncomeDataInRange2", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")

		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		UserId := "1"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange(StartDate, "invalidEndDate", UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
	t.Run("error TestGetIncomeDataInRange3", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")

		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		EndDate := "2022-12-30"
		UserId := "1"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange("invalidStartDate", EndDate, UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
	t.Run("error TestGetIncomeDataInRange4", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")

		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange("invalidStartDate", "invalidEndDate", "invalid")

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
}

func TestGetDateRange(t *testing.T) {
	t.Run("success GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, _, err := models.NewPostgreSQLDataFetcher("test")

		// テスト対象のデータ
		UserId := "1"
		// expectedData := []models.PaymentDate{
		// 	{
		// 		UserID:            1,
		// 		StratPaymaentDate: "2018-04-27",
		// 		EndPaymaentDate:   "2023-10-10",
		// 	},
		// }

		// テスト用の行データを設定
		// rows := sqlmock.NewRows([]string{
		// 	"user_id", "start_paymaent_date", "end_paymaent_date",
		// }).AddRow(
		// 	expectedData[0].UserID,
		// 	expectedData[0].StratPaymaentDate,
		// 	expectedData[0].EndPaymaentDate,
		// )

		// // モックに行データを設定
		// mock.ExpectQuery(`
		// 	SELECT user_id, MIN(payment_date) as "start_paymaent_date", MAX(payment_date) as "end_paymaent_date" from incomeforecast_incomeforecastdata
		// 	WHERE user_id = $1
		// 	GROUP BY user_id;`).
		// 	WithArgs(UserId).
		// 	WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetDateRange(UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 問題なく値が取得出来ていること
		assert.NotEmpty(t, result[0].UserID)
		assert.NotEmpty(t, result[0].StratPaymaentDate)
		assert.NotEmpty(t, result[0].EndPaymaentDate)
	})
	t.Run("success GetDateRange data empty", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, _, _ := models.NewPostgreSQLDataFetcher("test")

		// テスト対象のデータ
		UserId := "999"

		// テストを実行
		result, err := dbFetcher.GetDateRange(UserId)

		// エラーが期待通りに発生することを検証
		assert.Empty(t, result)

		t.Log(err)
	})
	t.Run("error GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher("test")

		UserId := "hoge"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetDateRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetDateRange(UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
}

func TestGetYearsIncomeAndDeduction(t *testing.T) {
	t.Run("success GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テスト対象のデータ
		UserId := "1"
		expectedData := models.YearsIncomeData{
			Years:           "2018",
			TotalAmount:     2904246,
			DeductionAmount: 450036,
			TakeHomeAmount:  2454210,
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"year", "sum_total_amount", "sum_deduction_amount", "sum_take_home_amount",
		}).AddRow(
			expectedData.Years,
			expectedData.TotalAmount,
			expectedData.DeductionAmount,
			expectedData.TakeHomeAmount,
		)

		// モックに行データを設定
		mock.ExpectQuery(DB.GetYearsIncomeAndDeductionSyntax).
			WithArgs(UserId).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetYearsIncomeAndDeduction(UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result[0])
	})
	t.Run("success GetYearsIncomeAndDeduction data empty", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, _, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テスト対象のデータ
		UserId := "999"

		// テストを実行
		result, err := dbFetcher.GetYearsIncomeAndDeduction(UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Empty(t, result)
	})
	t.Run("error GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		UserId := "hoge"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetYearsIncomeAndDeductionSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetYearsIncomeAndDeduction(UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
}

func TestInsertIncome(t *testing.T) {
	t.Run("success TestInsertIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		testData := []models.InsertIncomeData{
			{
				PaymentDate:     "9999-01-01",
				Age:             30,
				Industry:        "IT",
				TotalAmount:     1000,
				DeductionAmount: 200,
				TakeHomeAmount:  800,
				UpdateUser:      "user123",
				Classification:  "A",
				UserID:          1,
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.InsertIncomeSyntax).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.InsertIncome(testData)

		// エラーがないことを検証
		assert.NoError(t, err)

		t.Log(err)
	})
	t.Run("error TestInsertIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストデータを作成
		testData := []models.InsertIncomeData{
			{
				PaymentDate:     "9999-01-01",
				Age:             30,
				Industry:        "IT",
				TotalAmount:     3333,
				DeductionAmount: 2222,
				TakeHomeAmount:  1111,
				UpdateUser:      "user123",
				Classification:  "A",
				UserID:          999, // pkの値を違反させてエラー確認する
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.InsertIncomeSyntax).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("ERROR"))).
			WillReturnError(errors.New("INSERT FAILED"))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.InsertIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)

		t.Log("error", err)
	})
}

func TestUpdateIncome(t *testing.T) {
	t.Run("success TestUpdateIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストデータを作成
		testData := []models.UpdateIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d90a63bfd0", // 既存のレコードの ID
				PaymentDate:      "9999-01-01",
				Age:              30,
				Industry:         "IT",
				TotalAmount:      1200,
				DeductionAmount:  250,
				TakeHomeAmount:   950,
				Classification:   "B",
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.UpdateIncomeSyntax).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// UpdateIncome メソッドを呼び出し
		err = dbFetcher.UpdateIncome(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})

	t.Run("error TestUpdateIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストデータを作成
		testData := []models.UpdateIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d90a63bfrr656gtgtgtfd0", // エラー用のuuid
				PaymentDate:      "9999-01-01",
				Age:              30,
				Industry:         "IT",
				TotalAmount:      1200,
				DeductionAmount:  250,
				TakeHomeAmount:   950,
				Classification:   "B",
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.UpdateIncomeSyntax).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("ERROR"))).
			WillReturnError(errors.New("UPDATE FAILED"))
		mock.ExpectCommit()

		// UpdateIncome メソッドを呼び出し
		err = dbFetcher.UpdateIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)

		t.Log("error", err)
	})
}

func TestDeleteIncome(t *testing.T) {
	t.Run("success TestDeleteIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストデータを作成
		testData := []models.DeleteIncomeData{
			{
				IncomeForecastID: "57cbdd21-3cce-42f2-ad3c-2f727d7edae7", // 既存のレコードの ID
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.DeleteIncomeSyntax).
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// DeleteIncome メソッドを呼び出し
		err = dbFetcher.DeleteIncome(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("error TestDeleteIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストデータを作成
		testData := []models.DeleteIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d9frfrde450a63bfd0", // エラー用のuuid
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(DB.DeleteIncomeSyntax).
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("ERROR"))).
			WillReturnError(errors.New("DELETE FAILED"))
		mock.ExpectCommit()

		// DeleteIncome メソッドを呼び出し
		err = dbFetcher.DeleteIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)

		t.Log("error", err)
	})
}

// テストデータを削除するための関数
// func TestTestDataDelete(t *testing.T) {
// 	err := common.TestDataDelete()
// 	t.Log(err)
// }
