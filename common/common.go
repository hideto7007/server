// common/common.go
package common

import (
	// "fmt"

	"strconv"
	"time"

	// "net/http"
	"github.com/gin-gonic/gin"
)

type (
	CommonFetcher interface {
		IntgetPrameter(c *gin.Context, prams ...string) (map[string]int, error)
		TimeToStr(t time.Time) string
		StrToTime(dateStr string) (time.Time, error)
		StrToInt(str string) (int, error)
	}
	commonFetcherImpl struct{}
)

func NewCommonFetcher() CommonFetcher {
	return &commonFetcherImpl{}
}

// クエリーパラメータを整数値でまとめたマップで返す。
//
//
// 引数:
//
//	param1: コンテキスト
//	param2: 任意のクエリーパラメータ
//
// 戻り値:
//
//	戻り値1: nilと数値変換出来なかった際のエラー内容
//	戻り値2: 整数値が格納されたマップとnil

func (cf *commonFetcherImpl) IntgetPrameter(c *gin.Context, prams ...string) (map[string]int, error) {
	paramMap := map[string]int{}
	for _, keyParam := range prams {
		stringParam := c.DefaultQuery(keyParam, "0")
		intParam, err := strconv.Atoi(stringParam)
		if err != nil {
			return nil, err
		}
		paramMap[keyParam] = intParam
	}
	return paramMap, nil
}

// 日付を文字列に変換
//
// 引数:
//
//	param1: time.Time型
//
// 戻り値:
//
//	戻り値1: 日付を文字列変換

func (cf *commonFetcherImpl) TimeToStr(t time.Time) string {
	return t.Format("2006-01-02")
}

// 文字列日付を日付に変換
//
// 引数:
//
//	param1: string
//
// 戻り値:
//
//	戻り値1: 文字列を日付に変換

func (cf *commonFetcherImpl) StrToTime(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// 整数値型に変換
//
// 引数:
//
//	param1: string
//
// 戻り値:
//
//	戻り値1: 文字列を整数値型に変換して返す

func (cf *commonFetcherImpl) StrToInt(str string) (int, error) {
	var replaceInt int
	replaceInt, err := strconv.Atoi(str)
	if err != nil {
		return replaceInt, err
	}
	return replaceInt, nil
}

// 文字列型に変換
//
// 引数:
//
//	param1: interface{}
//
// 戻り値:
//
//	戻り値1: 整数値を文字列型に変換して返す

func AnyToStr[T interface{}](num T) string {
	// 型を直接判定する
	var replaceString string
	// デバッグ: 型と値を出力
	// fmt.Printf("Type: %T, Value: %v\n", num, num)
	switch v := any(num).(type) {
	case int:
		replaceString = strconv.Itoa(v)
	case float64:
		replaceString = strconv.FormatFloat(v, 'f', 0, 64)
	case string:
		replaceString = v
	}
	return replaceString
}

// // この関数はテストデータを削除するための独立関数
// func TestDataDelete() error {
// 	db, err := sql.Open("postgres", config.GetDataBaseSource())
// 	if err != nil {
// 		log.Printf("sql.Open error %s", err)
// 	}
// 	defer db.Close()

// 	// トランザクションを開始
// 	tx, err := db.Begin()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	deleteStatement := `
// 		DELETE FROM public.incomeforecast_incomeforecastdata
// 		WHERE update_user = 'user123';
// 	`

// 	if _, err = tx.Exec(deleteStatement); err != nil {
// 		tx.Rollback()
// 		fmt.Println(err)
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	return nil
// }

// クエリーパラメータを整数値でまとめたマップで返す。
//
//
// 引数:
//
//	param1: コンテキスト
//	param2: 任意のクエリーパラメータ
//
// 戻り値:
//
//	戻り値1: nilと数値変換出来なかった際のエラー内容
//	戻り値2: 整数値が格納されたマップとnil

// func IntgetPrameter(req *http.Request, prams ...string) (map[string]int, error) {
// 	paramMap := map[string]int{}
// 	for _, keyParam := range prams {
// 		stringParam := req.URL.Query().Get(keyParam)
// 		intParam, err := strconv.Atoi(stringParam)
// 		if err != nil {
// 			return nil, err
// 		}
// 		paramMap[keyParam] = intParam
// 	}
// 	return paramMap, nil
// }
