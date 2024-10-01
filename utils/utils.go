// utils/utils.go
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

// トークン生成関数
func GenerateJWT(UserId int) (string, error) {
	// トークンの有効期限を設定

	// トークンのクレーム（データペイロード）を作成
	claims := jwt.MapClaims{
		"UserId": UserId,
		"exp":    time.Now().Add(2 * time.Hour).Unix(),
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
