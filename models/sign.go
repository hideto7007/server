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
		GetSignIn(data RequestSignInData) (SignInData, error)
		PostSignUp(data RequestSignUpData) error
		PutSignInEdit(data RequestSignInEditData) error
		PutCheck(data RequestSignInEditData) (string, error)
		DeleteSignIn(data RequestSignInDeleteData) error
	}

	RequestSignInData struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestSignUpData struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
		NickName     string `json:"nick_name"`
	}

	RequestSignInEditData struct {
		UserId       interface{} `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
		UserName     string      `json:"user_name"`
		UserPassword string      `json:"user_password"`
	}

	RequestSignInDeleteData struct {
		UserId   interface{} `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
		UserName string      `json:"user_name"`
	}

	SignInData struct {
		UserId       int
		UserName     string
		UserPassword string
	}

	SignUpData struct {
		UserName     string
		UserPassword string
	}

	SignInEditData struct {
		UserId       int
		UserName     string
		UserPassword string
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
//   - data: { user_id: int, user_name: string, user_password: string }
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
	rows, err := pf.db.Query(DB.GetSignInSyntax, data.UserName)

	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	// if flg := rows.Next(); !flg {
	// 	fmt.Println(result, "存在しないユーザー名です。")
	// 	// return result, errors.New("存在しないユーザー名です。")
	// }

	// fmt.Println(rows.Next())

	for rows.Next() {
		var record SignInData
		err := rows.Scan(
			&record.UserId,
			&record.UserName,
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
		return result, errors.New("存在しないユーザー名です。")
	}

	return result, nil
}

// SignUp サインアップ情報を新規登録API
//
// 引数:
//   - data: { user_name: string, user_password: string, nick_name: string }
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
		fmt.Println(err)
		return err
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
		data.UserName,
		data.UserPassword,
		data.NickName,
		createdAt,
		data.NickName,
		createdAt,
		1); err != nil {
		return err
	}

	return nil
}

// PutSignInEdit サイン情報を編集API
//
// 引数:
//   - data: { user_id: int, user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SignDataFetcher) PutSignInEdit(data RequestSignInEditData) error {

	var err error
	createdAt := time.Now()
	// 初期値nullにするためポインター型で定義
	var userName *string
	var userPassword *string

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// deferでロールバックまたはコミットを管理
	defer func() {
		if p := recover(); p != nil || err != nil {
			tx.Rollback() // パニックまたはエラー発生時にロールバック
		} else {
			err = tx.Commit() // エラーがなければコミット
		}
	}()

	// ユーザー名及びユーザーパスワードが存在する場合、ポインターに変数代入
	if data.UserName != "" {
		userName = &data.UserName
	}

	if data.UserPassword != "" {
		// TBD:変更する値もハッシュ化する
		hashPassword, _ := pf.UtilsFetcher.EncryptPassword(data.UserPassword)
		userPassword = &hashPassword
	}

	signInEdit := DB.PutSignInEditSyntax

	if _, err = tx.Exec(signInEdit,
		userName,
		userPassword,
		createdAt,
		data.UserId); err != nil {
		return err
	}

	return nil
}

// PutCheck サイン情報修正した際に、ユーザー名かパスワードどちらを更新したかチェックする
//
// 引数:
//   - data: { user_id: int, user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: 文字列, nil(errorの場合error)
//

func (pf *SignDataFetcher) PutCheck(data RequestSignInEditData) (string, error) {

	var err error
	var result string

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetSignInSyntax, data.UserName)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return "", err
	}
	defer rows.Close()

	// `rows.Next()`で結果があるかを確認
	if rows.Next() {
		// 最初の行をスキャン（ユーザーネームが存在する）
		var record SignInData
		err := rows.Scan(
			&record.UserId,
			&record.UserName,
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
		result = "ユーザー名更新"
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

func (pf *SignDataFetcher) DeleteSignIn(data RequestSignInDeleteData) error {

	var err error

	// データベースのクローズをdeferで最初に宣言
	defer pf.db.Close()

	// トランザクションを開始
	tx, err := pf.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
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

	if _, err = tx.Exec(signInDelete,
		data.UserId); err != nil {
		return err
	}

	return nil
}
