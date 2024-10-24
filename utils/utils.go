// utils/utils.go
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenFetcher インターフェースの定義
type TokenFetcher interface {
	GenerateJWT(UserId int, ExpirationDate int) (string, error)
	NewToken(UserId int, ExpirationDate int) (string, error)
	RefreshToken(UserId int, ExpirationDate int) (string, error)
}

type TokenDataFetcher struct {
	JwtSecret []byte
}

type Response[T any] struct {
	RecodeRows int    `json:"recode_rows,omitempty"`
	Token      string `json:"token,omitempty"`
	Result     []T    `json:"result,omitempty"`
	ResultMsg  T      `json:"result_msg,omitempty"`
	ErrorMsg   string `json:"error_msg,omitempty"`
}

type Request struct {
	Data interface{} `json:"data"`
}

type ErrorStruct struct {
	Error    string `json:"error"`
	ErrorMsg string `json:"error_msg"`
}

type ErrorMessages struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func NewTokenFetcher(JwtSecret []byte) TokenFetcher {
	return &TokenDataFetcher{
		JwtSecret: JwtSecret,
	}
}

// トークン生成関数
func (tg *TokenDataFetcher) GenerateJWT(UserId int, ExpirationDate int) (string, error) {
	// トークンの有効期限を設定

	// トークンのクレーム（データペイロード）を作成
	// 検証時はtime.Now().Add(time.Duration(ExpirationDate) * time.Minute).Unix()で確認する
	claims := jwt.MapClaims{
		"UserId": UserId,
		"exp":    time.Now().Add(time.Duration(ExpirationDate) * time.Hour).Unix(),
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンのクレームを確認
	tokenClaims, _ := token.Claims.(jwt.MapClaims)
	fmt.Printf("生成されたトークンのクレーム: %+v\n", tokenClaims)

	fmt.Println("check ", tg.JwtSecret)

	// トークンに署名
	tokenString, err := token.SignedString(tg.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 新規有効期限付きのトークン発行
func (tg *TokenDataFetcher) NewToken(UserId int, ExpirationDate int) (string, error) {
	return tg.GenerateJWT(UserId, ExpirationDate)
}

// 新規有効期限付きのリフレッシュトークン発行
func (tg *TokenDataFetcher) RefreshToken(UserId int, ExpirationDate int) (string, error) {
	return tg.GenerateJWT(UserId, ExpirationDate)
}
