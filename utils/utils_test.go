package utils

import (
	"fmt"
	"testing"

	mock_utils "server/mock/utils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	t.Run("GenerateJWT token発行できる", func(t *testing.T) {
		tokenFetcher := NewUtilsFetcher(JwtSecret)

		token, err := tokenFetcher.GenerateJWT(1, 3)

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
		tokenFetcher := NewUtilsFetcher(JwtSecret)

		token, err := tokenFetcher.NewToken(1, 3)

		// クエリエラーが発生したことを確認
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("RefreshToken token発行できる", func(t *testing.T) {
		tokenFetcher := NewUtilsFetcher(JwtSecret)

		token, err := tokenFetcher.RefreshToken(1, 3)

		// クエリエラーが発生したことを確認
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
