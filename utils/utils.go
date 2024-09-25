// utils/utils.go
package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

// var tokenExpirationDate int = 60

// トークン生成関数
func GenerateJWT(UserId int) (string, error) {
	// トークンの有効期限を設定
	tokenExpirationTime := time.Now().Add(60 * time.Minute)

	// トークンのクレーム（データペイロード）を作成
	claims := jwt.MapClaims{
		"UserId": UserId,
		"exp":    tokenExpirationTime.Unix(),
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// トークンに署名
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
