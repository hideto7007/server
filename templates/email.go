package templates

import (
	"bytes"
	"text/template"
)

type (
	// データ構造体
	TemporayPostSingUpEmailData struct {
		Name        string
		ConfirmCode string
	}
)

var temporayPostSingUpTemplate = template.Must(template.New("auth_email").Parse(`
{{.Name}}さん！

ファイナンスアプリをご利用いただきありがとうございます。


確認コードは {{.ConfirmCode}} です。

この確認コードの有効期限は1時間です。
`))

func TemporayPostSingUpTemplate(Name string, ConfirmCode string) (string, string, error) {
	subject := "【ファイナンスアプリ】本登録を完了してください"
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := TemporayPostSingUpEmailData{
		Name:        Name,
		ConfirmCode: ConfirmCode,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := temporayPostSingUpTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}
