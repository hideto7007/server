// utils/utils.go
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Response struct {
	Token    string      `json:"token,omitempty"`
	Result   interface{} `json:"result,omitempty"`
	ErrorMsg string      `json:"error_msg,omitempty"`
}

type ErrorStruct struct {
	Error string `json:"error"`
	ErrorMsg string `json:"error_msg"` 
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

// トークン生成関数
func GenerateJWT(UserId int, ExpirationDate int) (string, error) {
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

	// トークンに署名
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 新規有効期限付きのトークン発行
func NewToken(UserId int, ExpirationDate int) (string, error) {
	return GenerateJWT(UserId, ExpirationDate)
}

// 新規有効期限付きのリフレッシュトークン発行
func RefreshToken(UserId int, ExpirationDate int) (string, error) {
	return GenerateJWT(UserId, ExpirationDate)
}
