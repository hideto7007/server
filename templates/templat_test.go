package templates

import (
	"bytes"
	"server/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailTemplateService(t *testing.T) {
	var Name string = "test"
	var UserEmail string = "test@example.com"
	var ConfirmCode string = "1234"
	var Update string = "ユーザーパスワード"
	var UpdateValue string = "Update"
	var DateTime string = "2024年12月07日 20:00"
	var Year string = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")
	var Link string = "http://exmaple/test"
	var NewPassword string = "Test12345!"

	t.Run("TemporayPostSignUpTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.TemporayPostSignUpTemplate(Name, ConfirmCode)

		data := TemporayPostSignUpEmailData{
			Name:        Name,
			ConfirmCode: ConfirmCode,
		}

		var expectedBody bytes.Buffer
		temporayPostSignUpTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】本登録を完了してください")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("PostSignUpTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.PostSignUpTemplate(Name, UserEmail, DateTime)

		data := GenericEmailData{
			Name:      Name,
			UserEmail: UserEmail,
			DateTime:  DateTime,
			Year:      Year,
		}

		var expectedBody bytes.Buffer
		postSignUpTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】登録を完了致しました")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("PostSignInEditTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.PostSignInEditTemplate(Update, UpdateValue, DateTime)

		data := GenericEmailData{
			Update:      Update,
			UpdateValue: UpdateValue,
			DateTime:    DateTime,
			Year:        Year,
		}

		var expectedBody bytes.Buffer
		postSignInEditTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】登録情報編集致しました")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("PostSignInTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.PostSignInTemplate(UserEmail, DateTime)

		data := GenericEmailData{
			UserEmail: UserEmail,
			DateTime:  DateTime,
			Year:      Year,
		}

		var expectedBody bytes.Buffer
		postSignInTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】サインイン致しました")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("DeleteSignInTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.DeleteSignInTemplate(Name, UserEmail, DateTime)

		data := GenericEmailData{
			Name:      Name,
			UserEmail: UserEmail,
			DateTime:  DateTime,
			Year:      Year,
		}

		var expectedBody bytes.Buffer
		deleteSignInTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】アカウント削除完了のお知らせ")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("SignOutTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.SignOutTemplate(UserEmail, DateTime)

		data := GenericEmailData{
			UserEmail: UserEmail,
			DateTime:  DateTime,
			Year:      Year,
		}

		var expectedBody bytes.Buffer
		deleteSignOutTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】サインアウトのお知らせ")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("RegisterEmailCheckNoticeTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.RegisterEmailCheckNoticeTemplate(Link, DateTime)

		data := GenericEmailData{
			Link:     Link,
			DateTime: DateTime,
			Year:     Year,
		}

		var expectedBody bytes.Buffer
		registerEmailCheckNoticeTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】パスワード再設定通知のお知らせ")
		assert.Equal(t, body, expectedBody.String())
	})

	t.Run("NewPasswordUpdateTemplate テンプレート", func(t *testing.T) {
		emailTemplateService := NewEmailTemplateManager()

		subject, body, err := emailTemplateService.NewPasswordUpdateTemplate(NewPassword, DateTime)

		data := GenericEmailData{
			NewPassword: NewPassword,
			DateTime:    DateTime,
			Year:        Year,
		}

		var expectedBody bytes.Buffer
		newPasswordUpdateTemplate.Execute(&expectedBody, data)

		assert.NoError(t, err)

		assert.Equal(t, subject, "【たくわえる】パスワード再発行成功のお知らせ")
		assert.Equal(t, body, expectedBody.String())
	})
}
