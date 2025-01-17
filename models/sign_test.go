package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"server/DB"
	"server/utils"
	"testing"

	mock_utils "server/mock/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnError(fmt.Errorf("クエリの実行に失敗しました"))

		// テストを実行
		_, err = dbFetcher.GetSignIn(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "クエリー実行エラー： クエリの実行に失敗しました")
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		}).AddRow(
			"test",
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// 正常な行データを用意
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		}).AddRow(
			1,
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ結果が返った後にエラーを発生させる
		rows.RowError(0, fmt.Errorf("forced row error"))

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
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
			UserEmail:    "test@exmple.com",
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
				UserEmail:    "test@exmple.com",
				UserPassword: string(hashedPassword),
			},
		}

		expectedData := mockData

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		})

		for _, data := range mockData {
			rows.AddRow(
				data.UserId,
				data.UserEmail,
				data.UserPassword,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		expectedData := "存在しないユーザー名です。"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		})

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		mockData := []SignInData{
			{
				UserId:       1,
				UserEmail:    "test@exmple.com",
				UserPassword: "Test12345!",
			},
		}

		expectedData := "パスワードが一致しませんでした。"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		})

		for _, data := range mockData {
			rows.AddRow(
				data.UserId,
				data.UserEmail,
				data.UserPassword,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
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

func TestGetExternalAuth(t *testing.T) {
	t.Run("GetExternalAuth クエリー実行時エラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserEmail := "test@exmple.com"

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetExternalAuthSyntax)).
			WithArgs(UserEmail).
			WillReturnError(fmt.Errorf("クエリの実行に失敗しました"))

		// テストを実行
		_, err = dbFetcher.GetExternalAuth(UserEmail)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "クエリー実行エラー： クエリの実行に失敗しました")
	})

	t.Run("GetExternalAuth rows.Scan時にエラー発生 UserIdで検証", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserEmail := "test@exmple.com"

		rows := sqlmock.NewRows([]string{"user_id", "user_email"}).
			AddRow("test", UserEmail)

		// モッククエリの期待値を設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetExternalAuthSyntax)).
			WithArgs(UserEmail). // ここで引数が一致しないとマッチしない
			WillReturnRows(rows)

		// テスト対象関数を実行
		var result []ExternalAuthData
		result, err = dbFetcher.GetExternalAuth(UserEmail)

		// エラーが発生したことを確認
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("GetExternalAuth rows.Err()にエラー発生", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserEmail := "test@exmple.com"

		rows := sqlmock.NewRows([]string{"user_id", "user_email"}).
			AddRow(1, UserEmail)
		rows.RowError(0, fmt.Errorf("forced row error"))

		// モッククエリの期待値を設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetExternalAuthSyntax)).
			WithArgs(UserEmail). // ここで引数が一致しないとマッチしない
			WillReturnRows(rows)

		// テスト対象関数を実行
		var result []ExternalAuthData
		result, err = dbFetcher.GetExternalAuth(UserEmail)

		// エラーが発生したことを確認
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "forced row error")
	})

	t.Run("GetExternalAuth 成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserEmail := "test@exmple.com"

		mockData := []ExternalAuthData{
			{
				UserId:    1,
				UserEmail: "test@exmple.com",
			},
		}

		expectedData := mockData

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email",
		})

		for _, data := range mockData {
			rows.AddRow(
				data.UserId,
				data.UserEmail,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetExternalAuthSyntax)).
			WithArgs(UserEmail).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetExternalAuth(UserEmail)

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

	t.Run("GetExternalAuth 存在しないユーザー名です", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		UserEmail := "test@exmple.com"

		expectedData := "存在しないユーザー名です。"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email",
		})

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetExternalAuthSyntax)).
			WithArgs(UserEmail).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.GetExternalAuth(UserEmail)

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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
			UserName:     "test",
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
			UserName:     "test",
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
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
			UserName:     "test",
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
			UserEmail:    "",
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
			UserEmail:    "test@exmple.com",
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
			UserEmail:    "test@exmple.com",
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
			UserEmail:    "test@exmple.com",
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

func TestPutCheck(t *testing.T) {
	t.Run("PutCheck クエリー実行時エラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInEditData{
			UserId:       "1",
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnError(fmt.Errorf("クエリの実行に失敗しました"))

		// テストを実行
		_, err = dbFetcher.PutCheck(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "クエリー実行エラー： クエリの実行に失敗しました")
	})

	t.Run("PutCheck rows.Scan時にエラー発生 UserIdで検証", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInEditData{
			UserId:       "test",
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		}).AddRow(
			"test",
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.PutCheck(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
	})

	t.Run("PutCheck rows.Err()にエラー発生", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInEditData{
			UserId:       "1",
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		// 正常な行データを用意
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		}).AddRow(
			1,
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ結果が返った後にエラーを発生させる
		rows.RowError(0, fmt.Errorf("forced row error"))

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnRows(rows)

		// テストを実行
		_, err = dbFetcher.PutCheck(requestData)

		// クエリエラーが発生したことを確認
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "forced row error")
	})

	t.Run("PutCheck 成功 パスワード更新", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// UtilsFetcher のモックを作成
		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("不一致"))

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInEditData{
			UserId:       "1",
			UserEmail:    "test@exmple.com",
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
				UserEmail:    "test@exmple.com",
				UserPassword: string(hashedPassword),
			},
		}

		expectedData := "パスワード更新"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		})

		for _, data := range mockData {
			rows.AddRow(
				data.UserId,
				data.UserEmail,
				data.UserPassword,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.PutCheck(requestData)

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

	t.Run("PutCheck 成功 ユーザー名更新", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト対象のデータ
		requestData := RequestSignInEditData{
			UserId:       "1",
			UserEmail:    "test@exmple.com",
			UserPassword: "Test12345!",
		}

		expectedData := "ユーザー名更新"

		// テスト用の行データを設定
		rows := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		})

		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(requestData.UserEmail).
			WillReturnRows(rows)

		// テストを実行
		result, err := dbFetcher.PutCheck(requestData)

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

func TestDeleteSignIn(t *testing.T) {
	t.Run("DeleteSignIn 削除成功", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		testData := RequestSignInDeleteData{
			UserId:    1,
			UserEmail: "text@example.com",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSignInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
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
			UserId:    1,
			UserEmail: "text@example.com",
		}

		// モックの準備
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.DeleteSignInSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
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

func TestGetUserId(t *testing.T) {
	UserEmail := "text@example.com"
	t.Run("GetUserId 登録ユーザーが存在しない", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(UserEmail).
			WillReturnError(sql.ErrNoRows)

		// テスト実行
		userId, err := dbFetcher.GetUserId(UserEmail)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, -1, userId)
		assert.Equal(t, err.Error(), "登録ユーザーが存在しません")
	})
	t.Run("GetUserId dbエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(UserEmail).
			WillReturnError(fmt.Errorf("dbエラー"))

		// テスト実行
		userId, err := dbFetcher.GetUserId(UserEmail)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, -1, userId)
		assert.Equal(t, err.Error(), "dbエラー")
	})
	t.Run("GetUserId 登録ユーザー存在", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// テスト用の行データを設定
		row := sqlmock.NewRows([]string{
			"user_id", "user_email", "user_password",
		}).AddRow(
			"1",
			"test@exmple.com",
			"Test12345!",
		)

		// クエリ実行時にエラーを返すようにモックを設定
		mock.ExpectQuery(regexp.QuoteMeta(DB.GetSignInSyntax)).
			WithArgs(UserEmail).
			WillReturnRows(row)

		// テスト実行
		userId, err := dbFetcher.GetUserId(UserEmail)

		// 検証
		assert.Equal(t, 1, userId)
		assert.Nil(t, err)
	})
}

func TestNewPasswordUpdate(t *testing.T) {
	Data := RequestNewPasswordUpdateData{
		TokenId:         "b2781af7-794a-1871-9865-bdc3c19291ff1",
		CurrentPassword: "Test12345!",
		NewUserPassword: "Test12345!",
		ConfirmPassword: "Test12345!",
	}
	UserId := Data.TokenId[utils.Uuid:]
	t.Run("NewPasswordUpdate 登録ユーザーが存在しない", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnError(sql.ErrNoRows)

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, err.Error(), "登録ユーザーが存在しません")
	})
	t.Run("NewPasswordUpdate dbエラー", func(t *testing.T) {
		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			utils.NewUtilsFetcher(utils.JwtSecret),
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnError(fmt.Errorf("dbエラー"))

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, err.Error(), "dbエラー")
	})
	t.Run("NewPasswordUpdate パスワードの不整合", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("エラー"))

		row := sqlmock.NewRows([]string{
			"user_email", "user_password",
		}).AddRow(
			"test@exmaple.com",
			"Test12345!",
		)
		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnRows(row)

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, "現在のパスワードと一致しませんでした。", err.Error())
	})
	t.Run("NewPasswordUpdate トランザクション失敗", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		row := sqlmock.NewRows([]string{
			"user_email", "user_password",
		}).AddRow(
			"test@exmaple.com",
			"Test12345!",
		)
		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnRows(row)

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(nil)

		// トランザクションの開始に失敗させる
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, "トランザクションの開始に失敗しました: transaction begin error", err.Error())
	})
	t.Run("NewPasswordUpdate 新しいパスワードと確認用のパスワードが一致しませんでした。", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		row := sqlmock.NewRows([]string{
			"user_email", "user_password",
		}).AddRow(
			"test@exmaple.com",
			"Test12345!",
		)
		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnRows(row)

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(nil)

		// トランザクション
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutPasswordSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 変更のパスワード一時保存
		tmpPs := Data.ConfirmPassword

		Data.ConfirmPassword = "Test1234567!"

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, "新しいパスワードと確認用のパスワードが一致しませんでした。", err.Error())
		Data.ConfirmPassword = tmpPs
	})
	t.Run("NewPasswordUpdate 更新失敗", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		row := sqlmock.NewRows([]string{
			"user_email", "user_password",
		}).AddRow(
			"test@exmaple.com",
			"Test12345!",
		)
		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnRows(row)

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(nil)

		mockUtilsFetcher.EXPECT().
			EncryptPassword(gomock.Any()).
			Return("hashPassword", nil)

		// トランザクション
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutPasswordSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(errors.New("update failed")) // Execの結果にエラーを返す
		mock.ExpectRollback() // エラー発生時にはロールバックを期待

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.Error(t, err)
		assert.Equal(t, "", userEmail)
		assert.Equal(t, "パスワード更新クエリの実行に失敗しました: update failed", err.Error())
	})
	t.Run("NewPasswordUpdate 更新成功", func(t *testing.T) {
		// gomock のコントローラを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// テスト用のDBモックを作成
		dbFetcher, mock, err := NewSignDataFetcher(
			"test",
			mockUtilsFetcher,
		)
		if err != nil {
			t.Fatalf("Error creating DB mock: %v", err)
		}

		row := sqlmock.NewRows([]string{
			"user_email", "user_password",
		}).AddRow(
			"test@exmaple.com",
			"Test12345!",
		)
		// モックの準備
		mock.ExpectQuery(regexp.QuoteMeta(DB.PasswordCheckSyntax)).
			WithArgs(UserId).
			WillReturnRows(row)

		mockUtilsFetcher.EXPECT().
			CompareHashPassword(gomock.Any(), gomock.Any()).
			Return(nil)

		mockUtilsFetcher.EXPECT().
			EncryptPassword(gomock.Any()).
			Return("hashPassword", nil)

		// トランザクション
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(DB.PutPasswordSyntax)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		// テスト実行
		userEmail, err := dbFetcher.NewPasswordUpdate(Data)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, "test@exmaple.com", userEmail)
		assert.Nil(t, err)
	})
}
