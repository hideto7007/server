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
