// models/singin.go
package models

import (
	"database/sql"
	"fmt"
	"log"
	"server/DB"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

type (
	SingInFetcher interface {
		GetSingIn(data RequestSingInData) (SingInData, error)
		PostSingUp(data RequestSingUpData) (SingUpData, error)
		PutSingInEdit(data RequestSingInEditData) (SingInEditData, error)
		DeleteSingInDelete(data RequestSingInDeleteData) (SingInDeleteData, error)
	}

	RequestSingInData struct {
		UserId       int    `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestSingUpData struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestSingInEditData struct {
		UserId       int    `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	RequestSingInDeleteData struct {
		UserId int `json:"user_id"` // stringにする理由、intだと内部で０に変換され本体の値の判定ができないためこのように指定する
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
		UserId int
	}

	SingDataFetcher struct{ db *sql.DB }
)

func NewSingDataFetcher(dataSourceName string) (*SingDataFetcher, sqlmock.Sqlmock, error) {
	if dataSourceName == "test" {
		db, mock, err := sqlmock.New()
		return &SingDataFetcher{db: db}, mock, err
	} else {
		// test実行時に以下のカバレッジは無視する
		db, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			log.Printf("sql.Open error %s", err)
		}
		return &SingDataFetcher{db: db}, nil, nil
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
	rows, err := pf.db.Query(DB.GetSingInSyntax, data.UserName, data.UserPassword)

	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data SingInData
		err := rows.Scan(
			&data.UserId,
			&data.UserName,
			&data.UserPassword,
		)

		if err != nil {
			return result, err
		}

		result = append(result, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// SingUp サインイン情報を新規登録API
//
// 引数:
//   - data: { user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) PostSingUp(data RequestSingUpData) (SingUpData, error) {

	var result SingUpData
	// var err error

	// データベースクエリを実行
	// rows, err := pf.db.Query(DB.GetSingInSyntax, data.UserName, data.UserPassword)

	// if err != nil {
	// 	fmt.Printf("Query failed: %v\n", err)
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var data SingUpData
	// 	err := rows.Scan(
	// 		&data.UserId,
	// 		&data.UserName,
	// 		&data.UserPassword,
	// 	)

	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	result = append(result, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return result, nil
}

// SingUp サインイン情報を編集API
//
// 引数:
//   - data: { user_id: int, user_name: string, user_password: string }
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) PutSingInEdit(data RequestSingInEditData) (SingInEditData, error) {

	var result SingInEditData
	// var err error

	// データベースクエリを実行
	// rows, err := pf.db.Query(DB.GetSingInSyntax, data.UserName, data.UserPassword)

	// if err != nil {
	// 	fmt.Printf("Query failed: %v\n", err)
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var data SingUpData
	// 	err := rows.Scan(
	// 		&data.UserId,
	// 		&data.UserName,
	// 		&data.UserPassword,
	// 	)

	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	result = append(result, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return result, nil
}

// SingUp サインイン情報を削除API
//
// 引数:
//   - data: { user_id: int}
//
// 戻り値:
//
//	戻り値1: サインイン情報を返し、ない場合は空のリスト
//	戻り値2: エラー内容(エラーがない場合はnil)
//

func (pf *SingDataFetcher) DeleteSingInDelete(data RequestSingInDeleteData) (SingInDeleteData, error) {

	var result SingInDeleteData
	// var err error

	// データベースクエリを実行
	// rows, err := pf.db.Query(DB.GetSingInSyntax, data.UserName, data.UserPassword)

	// if err != nil {
	// 	fmt.Printf("Query failed: %v\n", err)
	// 	return nil, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var data SingUpData
	// 	err := rows.Scan(
	// 		&data.UserId,
	// 		&data.UserName,
	// 		&data.UserPassword,
	// 	)

	// 	if err != nil {
	// 		return result, err
	// 	}

	// 	result = append(result, data)
	// }

	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return result, nil
}
