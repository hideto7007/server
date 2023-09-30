// common/common.go
package common

import (
	// "fmt"
	"strconv"
	// "time"
	// "net/http"
	"github.com/gin-gonic/gin"
)

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

func IntgetPrameter(c *gin.Context, prams ...string) (map[string]int, error) {
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
