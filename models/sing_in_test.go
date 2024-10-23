package models

import (
	"errors"
	"fmt"
	"regexp"
	"server/DB"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// func NewMocks
func TestGetSingIn(t *testing.T) {
	t.Run("GetSingIn クエリー実行時エラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSingInData{
			UserId:       1,
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSingInSyntax)).
			WithArgs(requestData.UserName, requestData.UserPassword).
			WillReturnError(fmt.Errorf("クエリの実行に失敗しました"))

		// テストを実行
		_, err = dbFetcher.GetSingIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "クエリの実行に失敗しました")
	})

	t.Run("GetSingIn rows.Scan時にエラー発生 UserIdで検証", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSingInData{
			UserId:       1,
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_name", "user_password",
		}).AddRow(
			"test",
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSingInSyntax)).
			WithArgs(requestData.UserName, requestData.UserPassword).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.GetSingIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
	})

	t.Run("GetSingIn rows.Err()にエラー発生", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSingInData{
			UserId:       1,
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// 正常な行データを用意
		rows := sqlmock.NewRows([]string{
			"user_id", "user_name", "user_password",
		}).AddRow(
			1,
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ結果が返った後にエラーを発生させる
		rows.RowError(0, fmt.Errorf("forced row error"))

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSingInSyntax)).
			WithArgs(requestData.UserName, requestData.UserPassword).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.GetSingIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forced row error")
	})

	t.Run("GetSingIn 成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSingInData{
			UserId:       1,
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		mockData := []SingInData{
			{
				UserId:       1,
				UserName:     "test@exmple.com",
				UserPassword: "Test12345!",
			},
		}

		expectedData := mockData

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_name", "user_password",
		})

		for _, data := range mockData {
			rows.AddRow(
				data.UserId,
				data.UserName,
				data.UserPassword,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSingInSyntax)).
			WithArgs(requestData.UserName, requestData.UserPassword).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetSingIn(requestData)

		// エラーがないことを検証
		assert.NoError(t, err)

		t.Log("result : ", result)
		// 取得したデータが期待値と一致することを検証
		assert.Equal(t, expectedData, result)

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock check: %s", err)
		}

	})
}

func TestPostSingUp(t *testing.T) {
	t.Run("PostSingUp 登録成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSingUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PostSingUpSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PostSingUp(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PostSingUp 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSingUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PostSingUpSyntax)).
			WithArgs(
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
		err = dbFetcher.PostSingUp(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insert failed")

		t.Log("error PostSingUp log", err)
	})
	t.Run("PostSingUp トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSingUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PostSingUp(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error PostSingUp log", err)
	})
}

func TestPutSingInEdit(t *testing.T) {
	t.Run("PutSingInEdit 登録成功 1", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSingInEditData{
			UserName:     "",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSingInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSingInEdit(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PutSingInEdit 登録成功 2", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSingInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSingInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSingInEdit(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PutSingInEdit 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSingInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSingInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("update failed")) // Execの結果にエラーを返す
		mock.ExpectRollback() // エラー発生時にはロールバックを期待

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSingInEdit(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")

		t.Log("error PutSingInEdit log", err)
	})
	t.Run("PutSingInEdit トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSingInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSingInEdit(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error PutSingInEdit log", err)
	})
}

func TestDeleteSingIn(t *testing.T) {
	t.Run("DeleteSingIn 登録成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSingInDeleteData{
			UserId: 1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSingInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSingIn(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("DeleteSingIn 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSingInDeleteData{
			UserId: 1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSingInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("delete failed")) // Execの結果にエラーを返す
		mock.ExpectRollback() // エラー発生時にはロールバックを期待

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSingIn(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")

		t.Log("error DeleteSingIn log", err)
	})
	t.Run("DeleteSingIn トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSingDataFetcher("test")
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSingInDeleteData{
			UserId: 1,
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSingIn(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error DeleteSingIn log", err)
	})
}
