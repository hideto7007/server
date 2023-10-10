package common

import (
	"server/config"
	"server/models"
	"testing"
	"time"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type MockIncomeData struct {
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

func TestGetIncomeDataInRange(t *testing.T) {
	t.Run("success TestGetIncomeDataInRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating DB mock: %v", err)
		}
		defer db.Close()

		// テスト対象のデータ
		startDate := "2022-11-01"
		endDate := "2023-03-30"
		expectedData := []models.IncomeData{
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
				IncomeForecastID: uuid.MustParse("f35a8725-16b8-406b-8521-186020a7fff6"),
				PaymentDate:      time.Date(2023, time.January, 10, 0, 0, 0, 0, time.FixedZone("", 0)),
				Age:              "28",
				Industry:         "システム開発",
				TotalAmount:      300000,
				DeductionAmount:  93000,
				TakeHomeAmount:   207000,
				Classification:   "給料",
				UserID:           1,
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"income_forecast_id", "payment_date", "age", "industry", "total_amount",
			"deduction_amount", "take_home_amount", "classification", "user_id",
		}).AddRow(
			expectedData[0].IncomeForecastID.String(),
			expectedData[0].PaymentDate,
			expectedData[0].Age,
			expectedData[0].Industry,
			expectedData[0].TotalAmount,
			expectedData[0].DeductionAmount,
			expectedData[0].TakeHomeAmount,
			expectedData[0].Classification,
			expectedData[0].UserID,
		)

		// モックに行データを設定
		mock.ExpectQuery(`
			SELECT income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, classification, user_id 
			FROM incomeforecast_incomeforecastdata 
			WHERE payment_date BETWEEN $1 AND $2`).
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		// テスト対象のPostgreSQLDataFetcherを作成
		dataFetcher := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストを実行
		result, err := dataFetcher.GetIncomeDataInRange(startDate, endDate)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)
	})
	t.Run("error TestGetIncomeDataInRange", func(t *testing.T) {
		// テスト用のDBモックを作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating DB mock: %v", err)
		}
		defer db.Close()

		// モックに行データを設定
		mock.ExpectQuery(`
			SELECT income_forecast_id, payment_date, age, industry, total_amount, deduction_amount, take_home_amount, classification, user_id 
			FROM incomeforecast_incomeforecastdata 
			WHERE payment_date BETWEEN $1 AND $2`).
			WillReturnError(sql.ErrNoRows)

		// テスト対象のPostgreSQLDataFetcherを作成
		dataFetcher := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// エラーケースをテスト
		_, err = dataFetcher.GetIncomeDataInRange("invalidStartDate", "invalidEndDate")

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
}

func TestGetStartDataAndEndDate(t *testing.T) {
	t.Run("success GetStartDataAndEndDate", func(t *testing.T) {
		// テスト用のDBモックを作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating DB mock: %v", err)
		}
		defer db.Close()

		// テスト対象のデータ
		UserId := "1"
		expectedData := []models.PaymentDate{
			{
				UserID:            1,
				StratPaymaentDate: "2018-04-27",
				EndPaymaentDate:   "2023-01-10",
			},
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "start_paymaent_date", "end_paymaent_date",
		}).AddRow(
			expectedData[0].UserID,
			expectedData[0].StratPaymaentDate,
			expectedData[0].EndPaymaentDate,
		)

		// モックに行データを設定
		mock.ExpectQuery(`
			SELECT user_id, MIN(payment_date) as "start_paymaent_date", MAX(payment_date) as "end_paymaent_date" from incomeforecast_incomeforecastdata
			WHERE user_id = $1
			GROUP BY user_id;`).
			WithArgs(UserId).
			WillReturnRows(rows)

		// テスト対象のPostgreSQLDataFetcherを作成
		dataFetcher := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// テストを実行
		result, err := dataFetcher.GetStartDataAndEndDate(UserId)

		// エラーがないことを検証
		assert.NoError(t, err)

		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)
	})
	t.Run("error GetStartDataAndEndDate", func(t *testing.T) {
		// テスト用のDBモックを作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating DB mock: %v", err)
		}
		defer db.Close()

		UserId := "hoge"

		// モックに行データを設定
		mock.ExpectQuery(`
			SELECT user_id, MIN(payment_date) as "start_paymaent_date", MAX(payment_date) as "end_paymaent_date" from incomeforecast_incomeforecastdata
			WHERE user_id = $1
			GROUP BY user_id;`).
			WillReturnError(sql.ErrNoRows)

		// テスト対象のPostgreSQLDataFetcherを作成
		dataFetcher := models.NewPostgreSQLDataFetcher(config.DataSourceName)

		// エラーケースをテスト
		_, err = dataFetcher.GetStartDataAndEndDate(UserId)

		// エラーが期待通りに発生することを検証
		assert.Error(t, err)

		t.Log(err)
	})
}
