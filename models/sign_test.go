package models

import (
	"errors"
	"fmt"
	"regexp"
	"server/DB"
	"server/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetSignIn(t *testing.T) {
	t.Run("GetSignIn クエリー実行時エラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnError(fmt.Errorf("クエリの実行に失敗しました"))

		// テストを実行
		_, err = dbFetcher.GetSignIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "クエリの実行に失敗しました")
	})

	t.Run("GetSignIn rows.Scan時にエラー発生 UserIdで検証", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
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
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.GetSignIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
	})

	t.Run("GetSignIn rows.Err()にエラー発生", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
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
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.GetSignIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forced row error")
	})

	t.Run("GetSignIn 成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// 平文のパスワードをハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("Error hashing password: %v", err)
		}

		mockData := []SignInData{
			{
				UserId:       1,
				UserName:     "test@exmple.com",
				UserPassword: string(hashedPassword),
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

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetSignIn(requestData)

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

	t.Run("GetSignIn 存在しないユーザー名です", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		expectedData := "存在しないユーザー名です。"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_name", "user_password",
		})

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetSignIn(requestData)

		// 取得したデータが期待値と一致することを検証
		assert.Empty(t, result)
		assert.Equal(t, expectedData, err.Error())

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock check: %s", err)
		}
	})

	t.Run("GetSignIn パスワードが一致しません。", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
		}

		mockData := []SignInData{
			{
				UserId:       1,
				UserName:     "test@exmple.com",
				UserPassword: "Test12345!",
			},
		}

		expectedData := "パスワードが一致しませんでした。"

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

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserName).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetSignIn(requestData)

		// 取得したデータが期待値と一致することを検証
		assert.Empty(t, result)
		assert.Equal(t, expectedData, err.Error())

		// モックが期待通りのクエリを受け取ったか確認
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock check: %s", err)
		}
	})
}

func TestPostSignUp(t *testing.T) {
	t.Run("PostSignUp 登録成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSignUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PostSignUpSyntax)).
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
		err = dbFetcher.PostSignUp(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PostSignUp 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSignUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PostSignUpSyntax)).
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
		err = dbFetcher.PostSignUp(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insert failed")

		t.Log("error PostSignUp log", err)
	})
	t.Run("PostSignUp トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSignUpData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			NickName:     "test",
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PostSignUp(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error PostSignUp log", err)
	})
}

func TestPutSignInEdit(t *testing.T) {
	t.Run("PutSignInEdit 登録成功 1", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSignInEditData{
			UserName:     "",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSignInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSignInEdit(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PutSignInEdit 登録成功 2", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSignInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSignInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSignInEdit(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("PutSignInEdit 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSignInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutSignInEditSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("update failed")) // Execの結果にエラーを返す
		mock.ExpectRollback() // エラー発生時にはロールバックを期待

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSignInEdit(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")

		t.Log("error PutSignInEdit log", err)
	})
	t.Run("PutSignInEdit トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSignInEditData{
			UserName:     "test@exmple.com",
			UserPassword: "Test12345!",
			UserId:       1,
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.PutSignInEdit(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error PutSignInEdit log", err)
	})
}

func TestDeleteSignIn(t *testing.T) {
	t.Run("DeleteSignIn 登録成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSignInDeleteData{
			UserId: 1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSignInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSignIn(testData)

		// エラーがないことを検証
		assert.NoError(t, err)
	})
	t.Run("DeleteSignIn 失敗", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テストデータを作成
		testData := RequestSignInDeleteData{
			UserId: 1,
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSignInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("delete failed")) // Execの結果にエラーを返す
		mock.ExpectRollback() // エラー発生時にはロールバックを期待

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSignIn(testData)

		// エラーが発生すること
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")

		t.Log("error DeleteSignIn log", err)
	})
	t.Run("DeleteSignIn トランザクションエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テストデータを作成
		testData := RequestSignInDeleteData{
			UserId: 1,
		}

		// InsertIncomeメソッドを呼び出し
		err = dbFetcher.DeleteSignIn(testData)

		// エラーが発生することを検証
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")

		t.Log("transaction begin error DeleteSignIn log", err)
	})
}
