package utils

import (
	"fmt"
	"testing"
	"time"

	mock_utils "server/mock/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetBaseURL(t *testing.T) {
	t.Run("GetBaseURL 文字列取得できること", func(t *testing.T) {

		url := GetBaseURL()

		// クエリエラーが発生したことを確認
		assert.Equal(t, url, ":///money_management/")
	})
}

func TestGenerateJWT(t *testing.T) {
	t.Run("GenerateJWT token発行できる", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		token, err := utilsFetcher.GenerateJWT(1, 3)

		// クエリエラーが発生したことを確認
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
	t.Run("GenerateJWT token発行時にエラー", func(t *testing.T) {
		// TODO:現状テストは通るけどカバレッジは反映されない
		// 理由は実際のテスト対象の関数をテストしてるわけではなくただのmockをテストしてるため
		// 実装は後々行う
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUtilsFetcher := mock_utils.NewMockUtilsFetcher(ctrl)

		// 署名エラーを返すモックの挙動を定義
		mockUtilsFetcher.EXPECT().
			GenerateJWT(1, 3).
			Return("", fmt.Errorf("署名エラー"))

		// JWT トークンの生成をテスト（エラーが発生する）
		token, err := mockUtilsFetcher.GenerateJWT(1, 3)

		// エラーが発生することを確認
		assert.Error(t, err)
		assert.EqualError(t, err, "署名エラー")

		// トークンは空であることを確認
		assert.Empty(t, token)
	})
}

func TestNewToken(t *testing.T) {
	t.Run("NewToken token発行できる", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		token, err := utilsFetcher.NewToken(1, 3)

		// クエリエラーが発生したことを確認
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("RefreshToken token発行できる", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		token, err := utilsFetcher.RefreshToken(1, 3)

		// クエリエラーが発生したことを確認
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestEncryptPassword(t *testing.T) {
	t.Run("EncryptPassword ハッシュ化できること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)
		val := "test"

		// パスワードをハッシュ化
		result, err := utilsFetcher.EncryptPassword(val)
		assert.NoError(t, err, "ハッシュ化時にエラーが発生しました")

		// ハッシュが平文のパスワードと一致するかを確認
		err = utilsFetcher.CompareHashPassword(result, val)
		assert.NoError(t, err, "ハッシュが平文パスワードと一致しませんでした")
	})
}

func TestCompareHashPassword(t *testing.T) {
	t.Run("CompareHashPassword nilが返されること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)
		val := "test"

		// パスワードをハッシュ化
		hashedPassword, err := utilsFetcher.EncryptPassword(val)
		assert.NoError(t, err, "ハッシュ化時にエラーが発生しました")

		// ハッシュ化されたパスワードと元の平文パスワードを比較
		err = utilsFetcher.CompareHashPassword(hashedPassword, val)
		assert.NoError(t, err, "ハッシュが平文パスワードと一致しませんでした")
	})
	t.Run("CompareHashPassword errが発生すること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)
		val := "test"

		// ハッシュ化されたパスワードと元の平文パスワードを比較
		err := utilsFetcher.CompareHashPassword(val, val)
		assert.NotNil(t, err)
	})
}

func TestParseWithClaims(t *testing.T) {
	t.Run("ParseWithClaims トークンが返されること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)
		token, _ := utilsFetcher.NewToken(1, 3)

		_, err := utilsFetcher.ParseWithClaims(token)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
	t.Run("ParseWithClaims エラーが発生されること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		token, err := utilsFetcher.ParseWithClaims("token")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestMapClaims(t *testing.T) {
	t.Run("MapClaims クレームが返されて、trueが返ってくること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)
		token, _ := utilsFetcher.NewToken(1, 3)

		token1, _ := utilsFetcher.ParseWithClaims(token)

		claims, ok := utilsFetcher.MapClaims(token1.(*jwt.Token))

		assert.Equal(t, ok, true)
		assert.NotEmpty(t, claims)
	})
	t.Run("MapClaims クレームが空で、falseが返ってくること", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		var token *jwt.Token // nil トークンを渡す

		// クレームを取得
		claims, ok := utilsFetcher.MapClaims(token)

		// ok が false であることを確認
		assert.Equal(t, false, ok)
		// claims が空であることを確認
		assert.Empty(t, claims)
	})
}

func TestSendMail(t *testing.T) {
	t.Run("SendMail エラーが起きないこと 2", func(t *testing.T) {
		// gomail のモックを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		// MockMailDialer の作成
		mockMailDialer := mock_utils.NewMockMailDialer(ctrl)

		// MockMailDialer の挙動を設定
		mockMailDialer.EXPECT().
			DialAndSend(gomock.Any()).
			Return(nil)

		// UtilsDataFetcher のモック設定
		utilsFetcher := &UtilsDataFetcher{
			MailDialer: mockMailDialer,
		}

		// テスト用の引数
		toEmail := "recipient@example.com"
		subject := "テスト件名"
		body := "テスト本文"

		// テスト対象の関数を呼び出し
		err := utilsFetcher.SendMail(toEmail, subject, body, true) // HTMLメール

		// エラーがないことを確認
		assert.NoError(t, err)
		// assert.NoError(t, err2)
	})

	t.Run("SendMail エラーが起きないこと 2", func(t *testing.T) {
		// gomail のモックを作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		// MockMailDialer の作成
		mockMailDialer := mock_utils.NewMockMailDialer(ctrl)

		// MockMailDialer の挙動を設定
		mockMailDialer.EXPECT().
			DialAndSend(gomock.Any()).
			Return(nil)

		// UtilsDataFetcher のモック設定
		utilsFetcher := &UtilsDataFetcher{
			MailDialer: mockMailDialer,
		}

		// テスト用の引数
		toEmail := "recipient@example.com"
		subject := "テスト件名"
		body := "テスト本文"

		// テスト対象の関数を呼び出し
		err := utilsFetcher.SendMail(toEmail, subject, body, false) // プレーンテキストメール

		// エラーがないことを確認
		assert.NoError(t, err)
	})

	t.Run("SendMail エラー発生", func(t *testing.T) {
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		// テスト用の引数
		toEmail := "recipient@example.com"
		subject := "テスト件名"
		body := "テスト本文"

		// テスト対象の関数を呼び出し
		err := utilsFetcher.SendMail(toEmail, subject, body, false) // プレーンテキストメール

		// エラーがないことを確認
		assert.Equal(t, err.Error(), "dial tcp :0: connect: connection refused")
	})
}

func TestDateTimeStr(t *testing.T) {
	t.Run("DateTimeStr エラーが起きないこと 2", func(t *testing.T) {
		// 日本標準時（JST: UTC+9）を明示的に設定
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)

		// 指定した日時を生成
		specifiedTime := time.Date(
			2024,          // 年
			time.December, // 月
			7,             // 日
			14,            // 時
			30,            // 分
			0,             // 秒
			0,             // ナノ秒
			jst,           // タイムゾーン
		)
		utilsFetcher := NewUtilsFetcher(JwtSecret)

		result := utilsFetcher.DateTimeStr(specifiedTime, "2006年01月02日 15:04")

		// エラーがないことを確認
		assert.Equal(t, result, "2024年12月07日 14:30")
	})
}
