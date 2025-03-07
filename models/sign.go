// models/sign.go
package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/DB"
	"server/utils"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

type (
	SignInFetcher interface {
		GetSignIn(data RequestSignInData) ([]SignInData, error)
		GetExternalAuth(UserEmail string) ([]ExternalAuthData, error)
		PostSignUp(data RequestSignUpData) error
		PutSignInEdit(UserId int, data RequestSignInEditData) error
		PutCheck(data RequestSignInEditData) (string, error)
		DeleteSignIn(userId int, data RequestSignInDeleteData) error
		GetUserId(UserEmail string) (int, error)
		NewPasswordUpdate(data RequestNewPasswordUpdateData) (string, error)
	}

	RequestSignInData struct {
		UserEmail    string `json:"user_email"`
		UserPassword string `json:"user_password"`
	}

	RequestSignUpData struct {
		UserEmail    string `json:"user_email"`
		UserPassword string `json:"user_password"`
		UserName     string `json:"user_name"`
	}

	RequestSignInEditData struct {
		UserEmail    string `json:"user_email"`
		UserPassword string `json:"user_password"`
	}

	RequestSignInDeleteData struct {
		UserEmail  string `json:"user_email"`
		DeleteName string `json:"delete_name"`
	}

	RequestNewPasswordUpdateData struct {
		TokenId         string `json:"token_id"`
		NewUserPassword string `json:"new_user_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	SignInData struct {
		UserId       int
		UserEmail    string
		UserPassword string
	}

	ExternalAuthData struct {
		UserId    int
		UserEmail string
	}

	SignUpData struct {
		UserEmail    string
		UserPassword string
	}

	SignInEditData struct {
		UserId       int
		UserEmail    string
		UserPassword string
	}

	NewPasswordUpdateData struct {
		UserEmail string
	}

	SignInDeleteData struct {
		UserId interface{}
	}

	SignDataFetcher struct {
		db           *sql.DB
		UtilsFetcher utils.UtilsFetcher
	}
)

func NewSignDataFetcher(dataSourceName string, UtilsFetcher utils.UtilsFetcher) (*SignDataFetcher, sqlmock.Sqlmock, error) {
	if dataSourceName == "test" {
		db, mock, err := sqlmock.New()
		return &SignDataFetcher{db: db, UtilsFetcher: UtilsFetcher}, mock, err
	} else {
		// test実行時に以下のカバレッジは無視する
		db, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("sql.Open error %s", err)
		}
		return &SignDataFetcher{db: db, UtilsFetcher: UtilsFetcher}, nil, nil
	}
}

// SignIn サインイン情報を返す
//
// 引数:
//   - data: { user_id: int, user_email: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) GetSignIn(data RequestSignInData) ([]SignInData, error) {

	var result []SignInData
	var err error

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetSignInSyntax, data.UserEmail)

	if err != nil {
		return nil, fmt.Errorf("クエリー実行エラー： %v", err)
	}
	defer rows.Close()

	// if flg := rows.Next(); !flg {
	// 	fmt.Println(result, "存在しないメールアドレスです。")
	// 	// return result, errors.New("存在しないメールアドレスです。")
	// }

	// fmt.Println(rows.Next())

	for rows.Next() {
		var record SignInData
		err := rows.Scan(
			&record.UserId,
			&record.UserEmail,
			&record.UserPassword,
		)
		if err != nil {
			return result, err
		}

		// パスワードの整合性を確認
		err = pf.UtilsFetcher.CompareHashPassword(record.UserPassword, data.UserPassword)
		if err == nil {
			// パスワードが一致する場合のみ結果に追加
			result = append(result, record)
		} else {
			return result, errors.New("パスワードが一致しませんでした。")
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return result, errors.New("存在しないメールアドレスです。")
	}

	return result, nil
}

// GetExternalAuth 外部認証のサインイン情報を返す
//
// 引数:
//   - data: { user_id: int, user_email: string}
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) GetExternalAuth(UserEmail string) ([]ExternalAuthData, error) {

	var result []ExternalAuthData
	var err error

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetExternalAuthSyntax, UserEmail)

	if err != nil {
		return nil, fmt.Errorf("クエリー実行エラー： %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record ExternalAuthData
		err := rows.Scan(
			&record.UserId,
			&record.UserEmail,
		)
		if err != nil {
			return result, err
		}
		result = append(result, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return result, errors.New("存在しないメールアドレスです。")
	}

	return result, nil
}

// SignUp サインアップ情報を新規登録API
//
// 引数:
//   - data: { user_email: string, user_password: string, user_name: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) PostSignUp(data RequestSignUpData) error {

	var err error
	createdAt := time.Now()

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// deferでロールバックまたはコミットを管理
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback() // パニックまたはエラー発生時にロールバック
		} else {
			err = tx.Commit() // エラーがなければコミット
		}
	}()

	signUp := DB.PostSignUpSyntax

	if _, err = tx.Exec(signUp,
		data.UserEmail,
		data.UserPassword,
		data.UserName,
		createdAt,
		data.UserName,
		createdAt,
		1); err != nil {
		return err
	}

	return nil
}

// PutSignInEdit サイン情報を編集API
//
// 引数:
//   - data: { user_email: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) PutSignInEdit(UserId int, data RequestSignInEditData) error {

	var err error
	updateAt := time.Now()
	// 初期値nullにするためポインター型で定義
	var userEmail *string
	var userPassword *string

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// deferでロールバックまたはコミットを管理
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback() // パニックまたはエラー発生時にロールバック
		} else {
			err = tx.Commit() // エラーがなければコミット
		}
	}()

	// メールアドレス及びユーザーパスワードが存在する場合、ポインターに変数代入
	if data.UserEmail != "" {
		userEmail = &data.UserEmail
	}

	if data.UserPassword != "" {
		// TBD:変更する値もハッシュ化する
		hashPassword, _ := pf.UtilsFetcher.EncryptPassword(data.UserPassword)
		userPassword = &hashPassword
	}

	signInEdit := DB.PutSignInEditSyntax

	if _, err = tx.Exec(signInEdit,
		userEmail,
		userPassword,
		updateAt,
		UserId); err != nil {
		return err
	}

	return nil
}

// PutCheck サイン情報修正した際に、メールアドレスかパスワードどちらを更新したかチェックする
//
// 引数:
//   - data: { user_id: int, user_email: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: 文字列, nil(errorの場合error)
//

func (pf *SignDataFetcher) PutCheck(data RequestSignInEditData) (string, error) {

	var err error
	var result string

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetSignInSyntax, data.UserEmail)
	if err != nil {
		return "", fmt.Errorf("クエリー実行エラー： %v", err)
	}
	defer rows.Close()

	// `rows.Next()`で結果があるかを確認
	if rows.Next() {
		// 最初の行をスキャン（ユーザーネームが存在する）
		var record SignInData
		err := rows.Scan(
			&record.UserId,
			&record.UserEmail,
			&record.UserPassword,
		)
		if err != nil {
			return "", err
		}

		// パスワードの整合性を確認
		err = pf.UtilsFetcher.CompareHashPassword(record.UserPassword, data.UserPassword)
		if err != nil {
			// パスワードが異なる場合（更新の必要あり）
			result = "パスワード更新"
		}
	} else {
		// ユーザーネームが存在しない場合
		result = "メールアドレス更新"
	}

	// `rows.Err()`でカーソル操作中のエラーを確認
	if err := rows.Err(); err != nil {
		return "", err
	}

	return result, nil
}

// DeleteSignIn サインイン情報を削除API
//
// 引数:
//   - data: { user_id: int}
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) DeleteSignIn(userId int, data RequestSignInDeleteData) error {

	var err error

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// deferでロールバックまたはコミットを管理
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback() // パニックまたはエラー発生時にロールバック
		} else {
			err = tx.Commit() // エラーがなければコミット
		}
	}()

	signInDelete := DB.DeleteSignInSyntax

	if _, err = tx.Exec(
		signInDelete,
		userId,
		data.UserEmail,
	); err != nil {
		return err
	}

	return nil
}

// GetUserId user_idを返す
//
// 引数:
//   - data: { user_email: string }
//
// 戻り値:
//
//	戻り値1: user_idを返し、ない場合は-1
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) GetUserId(UserEmail string) (int, error) {

	// データベースクエリを実行
	row := pf.db.QueryRow(DB.GetSignInSyntax, UserEmail)
	var record SignInData
	if err := row.Scan(&record.UserId, &record.UserEmail, &record.UserPassword); err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("登録ユーザーが存在しません")
		}
		return -1, err
	}
	return record.UserId, nil
}

// NewPasswordUpdate 新しいパスワードに更新
//
// 引数:
//   - data: { current_password: string, new_user_password: string, confirm_password: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) NewPasswordUpdate(data RequestNewPasswordUpdateData) (string, error) {

	var hashPassword string

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// 1. 登録済みのユーザーパスワードの整合性チェック

	userId := data.TokenId[utils.Uuid:]

	// データベースクエリを実行
	row := pf.db.QueryRow(DB.PasswordCheckSyntax, userId)
	var record NewPasswordUpdateData
	if err := row.Scan(&record.UserEmail); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("登録ユーザーが存在しません")
		}
		return "", err
	}

	// 2. 新しいパスワードへ更新
	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		return "", fmt.Errorf("トランザクションの開始に失敗しました: %v", err)
	}

	// deferでロールバックまたはコミットを管理
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback() // パニックまたはエラー発生時にロールバック
		} else {
			err = tx.Commit() // エラーがなければコミット
		}
	}()

	if data.NewUserPassword != data.ConfirmPassword {
		return "", fmt.Errorf("新しいパスワードと確認用のパスワードが一致しませんでした。")
	}
	hashPassword, _ = pf.UtilsFetcher.EncryptPassword(data.NewUserPassword)

	if _, err = tx.Exec(
		DB.PutPasswordSyntax,
		hashPassword,
		time.Now(),
		userId,
	); err != nil {
		return "", fmt.Errorf("パスワード更新クエリの実行に失敗しました: %v", err)
	}

	return record.UserEmail, nil
}
