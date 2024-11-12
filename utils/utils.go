// utils/utils.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UtilsFetcher インターフェースの定義
type UtilsFetcher interface {
	GenerateJWT(UserId int, ExpirationDate int) (string, error)
	NewToken(UserId int, ExpirationDate int) (string, error)
	RefreshToken(UserId int, ExpirationDate int) (string, error)
	EncryptPassword(password string) (string, error)
	CompareHashPassword(hashedPassword, requestPassword string) error
	ParseWithClaims(validationToken string) (interface{}, error)
	MapClaims(token *jwt.Token) (interface{}, bool)
}

type UtilsDataFetcher struct {
	JwtSecret []byte
}

// ResponseWithSlice with slice Result
type ResponseWithSlice[T any] struct {
	RecodeRows int    `json:"recode_rows,omitempty"`
	Token      string `json:"token,omitempty"`
	Result     []T    `json:"result,omitempty"`
	ErrorMsg   string `json:"error_msg,omitempty"`
}

// ResponseWithSlice with single Result
type ResponseWithSingle[T any] struct {
	RecodeRows int    `json:"recode_rows,omitempty"`
	Token      string `json:"token,omitempty"`
	Result     T      `json:"result,omitempty"`
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

var AuthToken = "AuthToken"
var RefreshAuthToken = "RefreshAuthToken"
var AuthTokenHour = 1

// 推奨：90日間だが、一旦12時間で設定
var RefreshAuthTokenHour = 12
var SecondsInHour = 3600

func NewUtilsFetcher(JwtSecret []byte) UtilsFetcher {
	return &UtilsDataFetcher{
		JwtSecret: JwtSecret,
	}
}

// トークン生成関数
func (ud *UtilsDataFetcher) GenerateJWT(UserId int, ExpirationDate int) (string, error) {
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
	// tokenClaims, _ := token.Claims.(jwt.MapClaims)
	// fmt.Printf("生成されたトークンのクレーム: %+v\n", tokenClaims)

	// トークンに署名
	tokenString, err := token.SignedString(ud.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 新規有効期限付きのトークン発行
func (ud *UtilsDataFetcher) NewToken(UserId int, ExpirationDate int) (string, error) {
	return ud.GenerateJWT(UserId, ExpirationDate)
}

// 新規有効期限付きのリフレッシュトークン発行
func (ud *UtilsDataFetcher) RefreshToken(UserId int, ExpirationDate int) (string, error) {
	return ud.GenerateJWT(UserId, ExpirationDate)
}

// パスワードの平文をハッシュ化
func (ud *UtilsDataFetcher) EncryptPassword(password string) (string, error) {
	// パスワードの文字列をハッシュ化する
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ハッシュを比較
func (ud *UtilsDataFetcher) CompareHashPassword(hashedPassword, requestPassword string) error {
	// パスワードの文字列をハッシュ化して、既に登録されているハッシュ化したパスワードと比較します
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword)); err != nil {
		return err
	}
	return nil
}

// トークンの検証
// テストの都合上、*jwt.Tokenだと厳密チェックができないためinterfaceで対応
func (ud *UtilsDataFetcher) ParseWithClaims(validationToken string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(validationToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// クレームからユーザー情報を取得
// テストの都合上、jwt.MapClaimsだと厳密チェックができないためinterfaceで対応
func (ud *UtilsDataFetcher) MapClaims(token *jwt.Token) (interface{}, bool) {
	if token == nil {
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}
