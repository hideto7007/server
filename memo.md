```
map宣言
// マップの宣言
var personMap map[string]string

// マップを初期化
personMap = make(map[string]string)

personMap := map[string]string{}

func IntgetPrameter(prams ...string) ([]int, error) {
	var paramList []int
	for _, param := range prams {
		intParam, err := strconv.Atoi(param)
		if err != nil {
			return nil, err
		}
		paramList = append(paramList, intParam)
	}
	return paramList, nil
}
```

```
package main

import (
	"fmt"
	"log"
	"net/http"
	"server/controllers"

	"github.com/rs/cors"
)

func main() {
	// CORSミドルウェアを設定
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           10, // 10秒間有効
	})

	// HTTPハンドラーを作成
	handler := http.NewServeMux()
	handler.HandleFunc("/delftstack", DelftstackHandler)
	handler.HandleFunc("/api/price", controllers.GetPriceInfoApi)

	// CORSミドルウェアを適用したサーバーを作成
	server := &http.Server{
		Addr:    ":8080", // ポート番号を指定
		Handler: corsMiddleware.Handler(handler),
	}

	log.Println("Listening for requests on :8080...")
	log.Fatal(server.ListenAndServe())
}

func DelftstackHandler(Response_Writer http.ResponseWriter, _ *http.Request) {
	Response_Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(Response_Writer, "Hello, this is delftstack.com!")
}


// controllers/price_management_controllers.go
package controllers

import (
	"encoding/json"
	"net/http"
	"server/common"
	// "github.com/gin-gonic/gin"
)

type PriceInfo struct {
	LeftAmount  int `json:"left_amount"`
	TotalAmount int `json:"total_amount"`
}

type Response struct {
	PriceInfo PriceInfo `json:"result"`
	Error     string    `json:"error,omitempty"`
}

func priceCalc(moneyReceived, bouns, fixedCost, loan, private int) PriceInfo {

	var priceinfo PriceInfo
	priceinfo.LeftAmount = moneyReceived - fixedCost - loan - private
	priceinfo.TotalAmount = (priceinfo.LeftAmount * 12) + bouns

	return priceinfo
}

// func GetPriceInfoApi(c *gin.Context) {
// 	// CORSヘッダーを設定
// 	// c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
// 	// c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	// c.Header("Access-Control-Allow-Headers", "Content-Type")

// 	data, err := common.IntgetPrameter(c, "money_received", "bouns", "fixed_cost", "loan", "private")

// 	if err == nil {
// 		res := priceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])

// 		response := Response{PriceInfo: res}

// 		c.JSON(http.StatusOK, gin.H{"message": response})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"message": err})
// 	}
//
// }

func GetPriceInfoApi(Response_Writer http.ResponseWriter, req *http.Request) {

	// レスポンスのContent-Typeを設定
	Response_Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	// HTTPリクエストからクエリパラメータを取得してパラメーラ値を整数値に変換
	data, err := common.IntgetPrameter(req, "money_received", "bouns", "fixed_cost", "loan", "private")

	var response Response
	if err == nil {
		res := priceCalc(data["money_received"], data["bouns"], data["fixed_cost"], data["loan"], data["private"])
		response.PriceInfo = res
	} else {
		// エラーメッセージをErrorフィールドに設定
		response.Error = err.Error()
	}

	// Response型のデータをJSONに変換
	jsonResponse, _ := json.Marshal(response)

	// JSONレスポンスを書き込む
	Response_Writer.WriteHeader(http.StatusOK)
	Response_Writer.Write(jsonResponse)
}


```

- DB 取得時のエラーについて
- 構造体と db 出力カラムが一致してないと以下のエラーが発生

```
expected 9 destination arguments in Scan, not 12

```

- 定数や関数定義する際
- 初めのイニシャルは大文字にする

- テストコマンド

```bash
go test -coverprofile="../coverage/coverage.out"
go tool cover -html=../coverage/coverage.out -o ../coverage/coverage.html
```
