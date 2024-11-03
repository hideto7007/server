// models/singin.go
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
	"golang.org/x/crypto/bcrypt"
)

type (
	SingInFetcher interface {
		GetSingIn(data RequestSingInData) (SingInData, error)
		PostSingUp(data RequestSingUpData) error
		PutSingInEdit(data RequestSingInEditData) error
		DeleteSingIn(data RequestSingInDeleteData) error
	}

	RequestSingInData struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestSingUpData struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
		NickName     string `json:"nick_name"`
	}

	RequestSingInEditData struct {
		UserId       interface{} `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
		UserName     string      `json:"user_name"`
		UserPassword string      `json:"user_password"`
	}

	RequestSingInDeleteData struct {
		UserId interface{} `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
	}

	SingInData struct {
		UserId       int
		UserName     string
		UserPassword string
	}

	SingUpData struct {
		UserName     string
		UserPassword string
	}

	SingInEditData struct {
		UserId       int
		UserName     string
		UserPassword string
	}

	SingInDeleteData struct {
		UserId interface{}
	}

	SingDataFetcher struct {
		db           *sql.DB
		UtilsFetcher utils.UtilsFetcher
	}
)

func NewSingDataFetcher(dataSourceName string, UtilsFetcher utils.UtilsFetcher) (*SingDataFetcher, sqlmock.Sqlmock, error) {
	if dataSourceName == "test" {
		db, mock, err := sqlmock.New()
		return &SingDataFetcher{db: db, UtilsFetcher: UtilsFetcher}, mock, err
	} else {
		// test実行時に以下のカバレッジは無視する
		db, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("sql.Open error %s", err)
		}
		return &SingDataFetcher{db: db, UtilsFetcher: UtilsFetcher}, nil, nil
	}
}

// SingIn サインイン情報を返す
//
// 引数:
//   - data: { user_id: int, user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) GetSingIn(data RequestSingInData) ([]SingInData, error) {

	var result []SingInData
	var err error

	// データベースクエリを実行
	rows, err := pf.db.Query(DB.GetSingInSyntax, data.UserName)

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
		var record SingInData
		err := rows.Scan(
			&record.UserId,
			&record.UserName,
			&record.UserPassword,
		)
		if err != nil {
			return result, err
		}

		// パスワードの整合性を確認
		err = bcrypt.CompareHashAndPassword([]byte(record.UserPassword), []byte(data.UserPassword))
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

// SingUp サインイン情報を新規登録API
//
// 引数:
//   - data: { user_name: string, user_password: string, nick_name: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) PostSingUp(data RequestSingUpData) error {

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

	singUp := DB.PostSingUpSyntax

	hashPassword, _ := pf.UtilsFetcher.EncryptPassword(data.UserPassword)

	if _, err = tx.Exec(singUp,
		data.UserName,
		hashPassword, // TBD:ここでハッシュ化して保存
		data.NickName,
		createdAt,
		data.NickName,
		createdAt,
		1); err != nil {
		return err
	}

	return nil
}

// SingUp サインイン情報を編集API
//
// 引数:
//   - data: { user_id: int, user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) PutSingInEdit(data RequestSingInEditData) error {

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

	singInEdit := DB.PutSingInEditSyntax

	if _, err = tx.Exec(singInEdit,
		userName,
		userPassword,
		createdAt,
		data.UserId); err != nil {
		return err
	}

	return nil
}

// DeleteSingIn サインイン情報を削除API
//
// 引数:
//   - data: { user_id: int}
//
// 戻り値:
//
//	戻り値1: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) DeleteSingIn(data RequestSingInDeleteData) error {

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

	singInDelete := DB.DeleteSingInSyntax

	if _, err = tx.Exec(singInDelete,
		data.UserId); err != nil {
		return err
	}

	return nil
}
