// utils/utils.go
package utils

import (
	"fmt"
	"server/common"
	"server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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
	SendMail(toEmail, subject, body string, isHTML bool) error
	DateTimeStr(t time.Time, format string) string
}

type MailDialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type SMTPMailDialer struct {
	dialer *gomail.Dialer
}

func NewSMTPMailDialer(host string, port int, userEmail, password string) MailDialer {
	return &SMTPMailDialer{
		dialer: gomail.NewDialer(host, port, userEmail, password),
	}
}

func (s *SMTPMailDialer) DialAndSend(m ...*gomail.Message) error {
	return s.dialer.DialAndSend(m...)
}

type UtilsDataFetcher struct {
	JwtSecret  []byte
	MailDialer MailDialer
}

type ErrorMessages struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ResponseData[T any] struct {
	RecodeRows int    `json:"recode_rows,omitempty"`
	Token      string `json:"token,omitempty"`
	Result     T      `json:"result,omitempty"`
}

type Request struct {
	Data interface{} `json:"data"`
}

type ErrorStruct struct {
	Error  string `json:"error"`
	Result string `json:"error_msg"`
}

var JwtSecret = []byte(config.GlobalEnv.JwtSecret)

var AuthToken = "auth_token"
var GoogleToken = "google_token"
var LineToken = "line_token"
var RefreshAuthToken = "refresh_auth_token"
var UserId = "user_id"
var OauthState = "oauth_state"
var AuthTokenHour = 1

var Uuid = 36

type ErrorValidationResponse = ResponseData[[]ErrorMessages]
type ErrorMessageResponse = ResponseData[string]

// 推奨：90日間だが、一旦12時間で設定
var RefreshAuthTokenHour = 12
var SecondsInHour = 3600

func NewUtilsFetcher(JwtSecret []byte) UtilsFetcher {
	var common common.CommonFetcher = common.NewCommonFetcher()
	smtpHost := config.GlobalEnv.SmtpHost                     // SMTPサーバー
	smtpPort, _ := common.StrToInt(config.GlobalEnv.SmtpPort) // SMTPポート
	fromEmail := config.GlobalEnv.FromEmail                   // 送信元メールアドレス
	password := config.GlobalEnv.EmailPassword                // 送信元メールのパスワード（またはアプリパスワード）
	mailDialer := NewSMTPMailDialer(smtpHost, smtpPort, fromEmail, password)
	return &UtilsDataFetcher{
		JwtSecret:  JwtSecret,
		MailDialer: mailDialer,
	}
}

// クライアントベースURL
func GetBaseURL() string {
	return fmt.Sprintf(
		"%s://%s%s",
		config.GlobalEnv.Protocol,
		config.GlobalEnv.ClinetDomain,
		"/money_management/",
	)
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

// SendMail はメールを送信する関数
func (ud *UtilsDataFetcher) SendMail(toEmail, subject, body string, isHTML bool) error {
	fromEmail := config.GlobalEnv.FromEmail // 送信元メールアドレス
	fromName := "たくわえる"

	// メール設定
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName)) // 送信元
	m.SetHeader("To", toEmail)                                // 送信先
	m.SetHeader("Subject", subject)                           // 件名

	// メール本文を設定 (HTMLまたはプレーンテキスト)
	if isHTML {
		m.SetBody("text/html", body) // HTML形式の本文
	} else {
		m.SetBody("text/plain", body) // プレーンテキスト形式の本文
	}

	// メール送信
	return ud.MailDialer.DialAndSend(m)
}

func (ud *UtilsDataFetcher) DateTimeStr(t time.Time, format string) string {
	return t.Format(format)
}
