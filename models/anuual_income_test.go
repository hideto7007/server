package models

import (
	"errors"
	"regexp"
	"server/DB"
	"server/common"
	"sort"
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
	MockPaymentData struct {
		PaymentDate string
		UserID      int
	}
	MockYearsIncomeData struct {
		Years           string
		TotalAmount     int
		DeductionAmount int
		TakeHomeAmount  int
	}
	MockUpdateIncomeData struct {
		IncomeForecastID string
		PaymentDate      string
		Age              int
		Industry         string
		TotalAmount      int
		DeductionAmount  int
		TakeHomeAmount   int
		UpdateUser       string
		Classification   string
	}
)

func TestNewPostgreSQLDataFetcher(t *testing.T) {
	// 通常のテストケース
	dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
	assert.NotNil(t, dbFetcher)
	assert.NotNil(t, mock)
	assert.NoError(t, err)

	// カバレッジを通すためのテストケース
	dbFetcher, _, err = NewPostgreSQLDataFetcher("")
	assert.NotNil(t, dbFetcher)
	assert.NoError(t, err)
}

// func NewMocks
func TestGetIncomeDataInRange(t *testing.T) {
	t.Run("success TestGetIncomeDataInRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		StartDate := "2022-11-01"
		EndDate := "2022-12-30"
		UserId := 1

		expectedData := []IncomeData{
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

		mockData := []IncomeData{
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
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		EndDate := "2022-12-30"

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange(StartDate, EndDate, 0)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error TestGetIncomeDataInRange1 log", err)
	})
	t.Run("error TestGetIncomeDataInRange2", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		UserId := 1

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange(StartDate, "invalidEndDate", UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error TestGetIncomeDataInRange2 log", err)
	})
	t.Run("error TestGetIncomeDataInRange3", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		EndDate := "2022-12-30"
		UserId := 1

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange("invalidStartDate", EndDate, UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error TestGetIncomeDataInRange3 log", err)
	})
	t.Run("error TestGetIncomeDataInRange4", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックに行データを設定
		mock.ExpectQuery(DB.GetIncomeDataInRangeSyntax).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetIncomeDataInRange("invalidStartDate", "invalidEndDate", 0)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error TestGetIncomeDataInRange4 log", err)
	})
	t.Run("error case rows.Scan TestGetIncomeDataInRange", func(t *testing.T) {
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		EndDate := "2022-12-30"
		UserId := 1

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"income_forecast_id", "payment_date", "age", "industry", "total_amount",
			"deduction_amount", "take_home_amount", "classification", "user_id",
		}).AddRow(
			"invalid-uuid", // 無効なUUIDを使用してScanのエラーを発生させる
			time.Date(2022, time.December, 23, 0, 0, 0, 0, time.UTC),
			"28",
			"システム開発",
			250000,
			78000,
			172000,
			"給料",
			1,
		)

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetIncomeDataInRangeSyntax)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), UserId).
			WillReturnRows(rows)

		_, err = dbFetcher.GetIncomeDataInRange(StartDate, EndDate, UserId)
		assert.Error(t, err)

		t.Log("error case rows.Scan TestGetIncomeDataInRange log", err)
	})

	t.Run("error case rows.Err TestGetIncomeDataInRange", func(t *testing.T) {
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		StartDate := "2022-11-01"
		EndDate := "2022-12-30"
		UserId := 1

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"income_forecast_id", "payment_date", "age", "industry", "total_amount",
			"deduction_amount", "take_home_amount", "classification", "user_id",
		}).AddRow(
			uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
			time.Date(2022, time.December, 23, 0, 0, 0, 0, time.UTC),
			"28",
			"システム開発",
			250000,
			78000,
			172000,
			"給料",
			1,
		)

		// 行エラーを設定
		rows.CloseError(errors.New("rows error"))

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetIncomeDataInRangeSyntax)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), UserId).
			WillReturnRows(rows)

		_, err = dbFetcher.GetIncomeDataInRange(StartDate, EndDate, UserId)
		assert.Error(t, err)

		t.Log("error case rows.Err TestGetIncomeDataInRange log", err)
	})
}

func TestGetDateRange(t *testing.T) {
	t.Run("success GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		var startTS int64 = 1525157431
		var endTS int64 = 1696488631
		var StratPaymaentDate time.Time
		var EndPaymaentDate time.Time
		var UserId int = 0
		var common common.CommonFetcher = common.NewCommonFetcher()
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserID := 1
		expectedData := []PaymentDate{
			{
				UserID:            1,
				StratPaymaentDate: "2018-04-27",
				EndPaymaentDate:   "2023-10-10",
			},
		}

		mockData := []MockPaymentData{
			{
				PaymentDate: "2018-04-20",
				UserID:      2,
			},
			{
				PaymentDate: "2018-04-27",
				UserID:      1,
			},
			{
				PaymentDate: "2019-05-27",
				UserID:      1,
			},
			{
				PaymentDate: "2020-06-27",
				UserID:      1,
			},
			{
				PaymentDate: "2021-07-27",
				UserID:      1,
			},
			{
				PaymentDate: "2022-08-27",
				UserID:      1,
			},
			{
				PaymentDate: "2022-09-27",
				UserID:      1,
			},
			{
				PaymentDate: "2023-10-10",
				UserID:      1,
			},
			{
				PaymentDate: "2023-10-15",
				UserID:      2,
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "start_paymaent_date", "end_paymaent_date",
		})

		for _, data := range mockData {
			if data.UserID == 1 {
				UserId = data.UserID
				dt, err := time.Parse("2006-01-02", data.PaymentDate)
				if err != nil {
					return
				}
				unix := dt.Unix()

				if unix <= startTS {
					startTS = unix
					StratPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}

				if unix >= endTS {
					endTS = unix
					EndPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}
			}
		}
		rows.AddRow(
			UserId,
			StratPaymaentDate,
			EndPaymaentDate,
		)

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetDateRangeSyntax)).
			WithArgs(UserID).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetDateRange(UserID)

		// エラーがないことを検証
		assert.NoError(t, err)

		t.Log("check = ", result)

		assert.Equal(t, expectedData[0].UserID, result[0].UserID)
		assert.Equal(t, expectedData[0].StratPaymaentDate, result[0].StratPaymaentDate)
		assert.Equal(t, expectedData[0].EndPaymaentDate, result[0].EndPaymaentDate)

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
	t.Run("falied GetDateRange data empty", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserId := 999

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetDateRangeSyntax)).
			WillReturnError(sql.ErrNoRows)

		// テストを実行
		result, err := dbFetcher.GetDateRange(UserId)

		// エラーが期待通りに発生することを検証
		assert.Empty(t, result)

		t.Log("falied GetDateRange data empty log", err)
	})
	t.Run("error GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		UserId := 0

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetDateRangeSyntax)).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetDateRange(UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error GetDateRange log", err)
	})
	t.Run("error case rows.Scan GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		var startTS int64 = 1525157431
		var endTS int64 = 1696488631
		var StratPaymaentDate time.Time
		var EndPaymaentDate time.Time
		var UserId string = "hoge"
		var common common.CommonFetcher = common.NewCommonFetcher()
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserID := 1

		mockData := []MockPaymentData{
			{
				PaymentDate: "2018-04-20",
				UserID:      2,
			},
			{
				PaymentDate: "2018-04-27",
				UserID:      1,
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "start_paymaent_date", "end_paymaent_date",
		})

		for _, data := range mockData {
			if data.UserID == 1 {
				dt, err := time.Parse("2006-01-02", data.PaymentDate)
				if err != nil {
					return
				}
				unix := dt.Unix()

				if unix <= startTS {
					startTS = unix
					StratPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}

				if unix >= endTS {
					endTS = unix
					EndPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}
			}
		}
		rows.AddRow(
			UserId, // 無効な値
			StratPaymaentDate,
			EndPaymaentDate,
		)

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetDateRangeSyntax)).
			WithArgs(UserID).
			WillReturnRows(rows)

		_, err = dbFetcher.GetDateRange(UserID)
		assert.Error(t, err)

		t.Log("error case rows.Scan GetDateRange log", err)
	})

	t.Run("error case rows.Err GetDateRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		var startTS int64 = 1525157431
		var endTS int64 = 1696488631
		var StratPaymaentDate time.Time
		var EndPaymaentDate time.Time
		var UserId int = 0
		var common common.CommonFetcher = common.NewCommonFetcher()
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		ParamUserId := 1

		mockData := []MockPaymentData{
			{
				PaymentDate: "2018-04-20",
				UserID:      2,
			},
			{
				PaymentDate: "2018-04-27",
				UserID:      1,
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "start_paymaent_date", "end_paymaent_date",
		})

		for _, data := range mockData {
			if data.UserID == 1 {
				UserId = data.UserID
				dt, err := time.Parse("2006-01-02", data.PaymentDate)
				if err != nil {
					return
				}
				unix := dt.Unix()

				if unix <= startTS {
					startTS = unix
					StratPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}

				if unix >= endTS {
					endTS = unix
					EndPaymaentDate, _ = common.StrToTime(data.PaymentDate)
				}
			}
		}
		rows.AddRow(
			UserId,
			StratPaymaentDate,
			EndPaymaentDate,
		)

		// 行エラーを設定
		rows.CloseError(errors.New("rows error"))

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetDateRangeSyntax)).
			WithArgs(ParamUserId).
			WillReturnRows(rows)

		_, err = dbFetcher.GetDateRange(ParamUserId)
		assert.Error(t, err)

		t.Log("error case rows.Err GetDateRange log", err)
	})
}

func TestGetYearsIncomeAndDeduction(t *testing.T) {
	t.Run("success GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}
		var Years string
		var TotalAmount int
		var DeductionAmount int
		var TakeHomeAmount int

		// テスト対象のデータ
		UserId := 1
		expectedData := []YearsIncomeData{
			{
				Years:           "2017",
				TotalAmount:     250000,
				DeductionAmount: 78000,
				TakeHomeAmount:  172000,
			},
			{
				Years:           "2018",
				TotalAmount:     500000,
				DeductionAmount: 156000,
				TakeHomeAmount:  344000,
			},
			{
				Years:           "2019",
				TotalAmount:     250000,
				DeductionAmount: 78000,
				TakeHomeAmount:  172000,
			},
		}

		mockData := []IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2018, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2019, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2017, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
		}
		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"year", "sum_total_amount", "sum_deduction_amount", "sum_take_home_amount",
		})

		// グループ化するためのデータ構造
		groupedData := make(map[int][]IncomeData)

		// データをグループ化
		for _, record := range mockData {
			key := record.PaymentDate.Year()
			groupedData[key] = append(groupedData[key], record)
		}

		// キーをソート
		var keys []int
		for key := range groupedData {
			keys = append(keys, key)
		}

		sort.Ints(keys)

		for _, key := range keys {
			records := groupedData[key]

			TotalAmount = 0
			DeductionAmount = 0
			TakeHomeAmount = 0
			Years = common.AnyToStr(key)
			for _, record := range records {
				if record.UserID == UserId {
					TotalAmount += record.TotalAmount
					DeductionAmount += record.DeductionAmount
					TakeHomeAmount += record.TakeHomeAmount
				}

			}
			if TotalAmount != 0 && DeductionAmount != 0 && TakeHomeAmount != 0 {
				rows.AddRow(
					Years,
					TotalAmount,
					DeductionAmount,
					TakeHomeAmount,
				)
			} else {
				rows.AddRow(nil, nil, nil, nil)
			}
		}

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetYearsIncomeAndDeductionSyntax)).
			WithArgs(UserId).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetYearsIncomeAndDeduction(UserId)

		t.Log("debug ", result)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("success GetYearsIncomeAndDeduction data empty", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserId := 999

		var Years string
		var TotalAmount int
		var DeductionAmount int
		var TakeHomeAmount int

		expectedData := []YearsIncomeData{
			{
				Years:           "",
				TotalAmount:     0,
				DeductionAmount: 0,
				TakeHomeAmount:  0,
			},
			{
				Years:           "",
				TotalAmount:     0,
				DeductionAmount: 0,
				TakeHomeAmount:  0,
			},
			{
				Years:           "",
				TotalAmount:     0,
				DeductionAmount: 0,
				TakeHomeAmount:  0,
			},
		}

		mockData := []IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2018, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2019, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2017, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
		}
		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"year", "sum_total_amount", "sum_deduction_amount", "sum_take_home_amount",
		})

		// グループ化するためのデータ構造
		groupedData := make(map[int][]IncomeData)

		// データをグループ化
		for _, record := range mockData {
			key := record.PaymentDate.Year()
			groupedData[key] = append(groupedData[key], record)
		}

		// キーをソート
		var keys []int
		for key := range groupedData {
			keys = append(keys, key)
		}

		sort.Ints(keys)

		for _, key := range keys {
			records := groupedData[key]

			TotalAmount = 0
			DeductionAmount = 0
			TakeHomeAmount = 0
			Years = common.AnyToStr(key)
			for _, record := range records {
				if record.UserID == UserId {
					TotalAmount += record.TotalAmount
					DeductionAmount += record.DeductionAmount
					TakeHomeAmount += record.TakeHomeAmount
				}

			}
			if TotalAmount != 0 && DeductionAmount != 0 && TakeHomeAmount != 0 {
				rows.AddRow(
					Years,
					TotalAmount,
					DeductionAmount,
					TakeHomeAmount,
				)
			} else {
				rows.AddRow("", 0, 0, 0)
			}
		}

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetYearsIncomeAndDeductionSyntax)).
			WithArgs(UserId).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetYearsIncomeAndDeduction(UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("error GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		UserId := 0

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetYearsIncomeAndDeductionSyntax)).
			WillReturnError(sql.ErrNoRows)

		// エラーケースをテスト
		_, err = dbFetcher.GetYearsIncomeAndDeduction(UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log("error GetYearsIncomeAndDeduction log", err)
	})
	t.Run("error case rows.Scan GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}
		var Years string
		var TotalAmount int
		var DeductionAmount int
		var TakeHomeAmount int

		// テスト対象のデータ
		UserId := 1

		mockData := []IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2018, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2019, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2017, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
		}
		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"year", "sum_total_amount", "sum_deduction_amount", "sum_take_home_amount",
		})

		// グループ化するためのデータ構造
		groupedData := make(map[int][]IncomeData)

		// データをグループ化
		for _, record := range mockData {
			key := record.PaymentDate.Year()
			groupedData[key] = append(groupedData[key], record)
		}

		// キーをソート
		var keys []int
		for key := range groupedData {
			keys = append(keys, key)
		}

		sort.Ints(keys)

		for _, key := range keys {
			records := groupedData[key]

			TotalAmount = 0
			DeductionAmount = 0
			TakeHomeAmount = 0
			Years = common.AnyToStr(key)
			for _, record := range records {
				if record.UserID == UserId {
					TotalAmount += record.TotalAmount
					DeductionAmount += record.DeductionAmount
					TakeHomeAmount += record.TakeHomeAmount
				}

			}
			rows.AddRow(
				Years,
				"TotalAmount", // 無効な値
				DeductionAmount,
				TakeHomeAmount,
			)
		}

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetYearsIncomeAndDeductionSyntax)).
			WithArgs(UserId).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.GetYearsIncomeAndDeduction(UserId)

		assert.Error(t, err)

		t.Log("error case rows.Scan GetYearsIncomeAndDeduction log", err)
	})

	t.Run("error case rows.Err GetYearsIncomeAndDeduction", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}
		var Years string
		var TotalAmount int
		var DeductionAmount int
		var TakeHomeAmount int

		// テスト対象のデータ
		UserId := 1

		mockData := []IncomeData{
			{
				IncomeForecastID: uuid.MustParse("8df939de-5a97-4f20-b41b-9ac355c16e36"),
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
			{
				IncomeForecastID: uuid.MustParse("92fa978b-876a-4693-b5af-a8d4010b4bfe"),
				PaymentDate:      time.Date(2018, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2019, time.November, 25, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2017, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
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
				PaymentDate:      time.Date(2018, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      250000,
				DeductionAmount:  78000,
				TakeHomeAmount:   172000,
				Classification:   "給料",
				UserID:           2,
			},
		}
		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"year", "sum_total_amount", "sum_deduction_amount", "sum_take_home_amount",
		})

		// グループ化するためのデータ構造
		groupedData := make(map[int][]IncomeData)

		// データをグループ化
		for _, record := range mockData {
			key := record.PaymentDate.Year()
			groupedData[key] = append(groupedData[key], record)
		}

		// キーをソート
		var keys []int
		for key := range groupedData {
			keys = append(keys, key)
		}

		sort.Ints(keys)

		for _, key := range keys {
			records := groupedData[key]

			TotalAmount = 0
			DeductionAmount = 0
			TakeHomeAmount = 0
			Years = common.AnyToStr(key)
			for _, record := range records {
				if record.UserID == UserId {
					TotalAmount += record.TotalAmount
					DeductionAmount += record.DeductionAmount
					TakeHomeAmount += record.TakeHomeAmount
				}

			}
			if TotalAmount != 0 && DeductionAmount != 0 && TakeHomeAmount != 0 {
				rows.AddRow(
					Years,
					TotalAmount,
					DeductionAmount,
					TakeHomeAmount,
				)
			} else {
				rows.AddRow(nil, nil, nil, nil)
			}
		}

		// 行エラーを設定
		rows.CloseError(errors.New("rows error"))

		// モックに行データを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetYearsIncomeAndDeductionSyntax)).
			WithArgs(UserId).
			WillReturnRows(rows)

		_, err = dbFetcher.GetYearsIncomeAndDeduction(UserId)
		assert.Error(t, err)

		t.Log("error case rows.Err GetYearsIncomeAndDeduction log", err)
	})
}

func TestInsertIncome(t *testing.T) {
	t.Run("success TestInsertIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := []InsertIncomeData{
			{
				PaymentDate:     "9999-01-01",
				Age:             30,
				Industry:        "IT",
				TotalAmount:     1000,
				DeductionAmount: 200,
				TakeHomeAmount:  800,
				Classification:  "A",
				UserID:          "1",
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.InsertIncomeSyntax)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.InsertIncome(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("error TestInsertIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := []InsertIncomeData{
			{
				PaymentDate:     "9999-01-01",
				Age:             30,
				Industry:        "IT",
				TotalAmount:     3333,
				DeductionAmount: 2222,
				TakeHomeAmount:  1111,
				Classification:  "A",
				UserID:          "999", // pkの値を違反させてエラー確認する
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.InsertIncomeSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.InsertIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insert failed")

		t.Log("error TestInsertIncome log", err)
	})
	t.Run("transaction begin error TestInsertIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := []InsertIncomeData{
			{
				PaymentDate:     "9999-01-01",
				Age:             30,
				Industry:        "IT",
				TotalAmount:     1000,
				DeductionAmount: 200,
				TakeHomeAmount:  800,
				Classification:  "A",
				UserID:          "1",
			},
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.InsertIncome(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error TestInsertIncome log", err)
	})
}

func TestUpdateIncome(t *testing.T) {
	t.Run("success TestUpdateIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := []UpdateIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d90a63bfd0", // 既存のレコードの ID
				PaymentDate:      "9999-01-01",
				Age:              30,
				Industry:         "IT",
				TotalAmount:      1200,
				DeductionAmount:  250,
				TakeHomeAmount:   950,
				UpdateUser:       "test_user",
				Classification:   "B",
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.UpdateIncomeSyntax)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// UpdateIncome メソッドを呼び出し
		err = dbFetcher.UpdateIncome(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})

	t.Run("error TestUpdateIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := []UpdateIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d90a63bfrr656gtgtgtfd0", // エラー用のuuid
				PaymentDate:      "9999-01-01",
				Age:              30,
				Industry:         "IT",
				TotalAmount:      1200,
				DeductionAmount:  250,
				TakeHomeAmount:   950,
				UpdateUser:       "test_user",
				Classification:   "B",
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.UpdateIncomeSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("update failed"))
		mock.ExpectCommit()

		// UpdateIncome メソッドを呼び出し
		err = dbFetcher.UpdateIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")

		t.Log("error", err)
	})
	t.Run("transaction begin error TestUpdateIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := []UpdateIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d90a63bfrr656gtgtgtfd0", // エラー用のuuid
				PaymentDate:      "9999-01-01",
				Age:              30,
				Industry:         "IT",
				TotalAmount:      1200,
				DeductionAmount:  250,
				TakeHomeAmount:   950,
				UpdateUser:       "test_user",
				Classification:   "B",
			},
		}

		// UpdateIncome メソッドを呼び出し
		err = dbFetcher.UpdateIncome(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error TestUpdateIncome log", err)
	})
}

func TestDeleteIncome(t *testing.T) {
	t.Run("success TestDeleteIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := []DeleteIncomeData{
			{
				IncomeForecastID: "57cbdd21-3cce-42f2-ad3c-2f727d7edae7", // 既存のレコードの ID
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteIncomeSyntax)).
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
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := []DeleteIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d9frfrde450a63bfd0", // エラー用のuuid
			},
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteIncomeSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("delete failed")) // Execの結果にエラーを返す
		mock.ExpectCommit()

		// DeleteIncome メソッドを呼び出し
		err = dbFetcher.DeleteIncome(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")

		t.Log("error", err)
	})
	t.Run("transaction begin error TestDeleteIncome", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewPostgreSQLDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := []DeleteIncomeData{
			{
				IncomeForecastID: "ecdb3762-9417-419d-c458-42d9frfrde450a63bfd0", // エラー用のuuid
			},
		}

		// DeleteIncome メソッドを呼び出し
		err = dbFetcher.DeleteIncome(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error TestDeleteIncome log", err)
	})
}

// テストデータを削除するための関数
// func TestTestDataDelete(t *testing.T) {
// 	err := common.TestDataDelete()
// 	t.Log(err)
// }
