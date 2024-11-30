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

	PostSingUpEmailData struct {
		Name             string
		UserName         string
		RegisterDateTime string
	}

	PostSingInEmailData struct {
		UserName       string
		SignInDateTime string
	}
)

var temporayPostSingUpTemplate = template.Must(template.New("auth_email").Parse(`
{{.Name}}さん！

ファイナンスアプリをご利用いただきありがとうございます。


確認コードは {{.ConfirmCode}} です。

この確認コードの有効期限は1時間です。
`))

var postSingUpTemplate = template.Must(template.New("sing_up").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<title>サインアップ完了のお知らせ</title>
		</head>
		<body>
			<h3>{{.Name}}さん！</h3>
			</br>
			<p>ファイナンスアプリのご登録が完了しました。誠にありがとうございます。</p>
			<p>{{.Name}}さんにとって、より良い資産管理が出来ることを祈っております。</p>
			</br>
			<h4>登録ユーザ名</h4>
			<p>{{.UserName}}</p>
			</br>
			<h4>登録日時</h4>
			<p>{{.RegisterDateTime}}</p>
			</br>
			<p>今後ともよろしくお願いいたします。</p>
		</body>
	</html>
`))

var postSingInTemplate = template.Must(template.New("sing_in").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<title>サインイン通知</title>
			<style>
				/* 全体のスタイル */
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f9f9f9;
				}

				/* メインコンテナ */
				.container {
					width: 100%;
					max-width: 600px;
					margin: 20px auto;
					background: #ffffff;
					border: 1px solid #ddd;
					border-radius: 5px;
					overflow: hidden;
				}

				/* ヘッダー部分 */
				.header {
					background-color: #007bff; /* 青色 */
					color: white;
					text-align: center;
					padding: 20px;
					font-size: 20px;
					font-weight: bold;
				}

				/* 本文部分 */
				.body {
					padding: 20px;
					color: #333;
				}

				/* 情報セクション */
				.info-section {
					background-color: #f1f1f1; /* グレー */
					padding: 15px;
					margin: 10px 0;
					border-radius: 5px;
				}

				/* フッター部分 */
				.footer {
					padding: 10px;
					text-align: center;
					font-size: 12px;
					color: #666;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					ファイナンスアプリ<br>
					Webブラウザから新たなサインインがありました
				</div>
				<div class="body">
					<p>いつもファイナンスアプリをご利用いただき、誠にありがとうございます。</p>
					<p>お客様がご利用中の登録ユーザーで、新たなサインインがありました。</p>

					<div class="info-section">
						<h4>登録ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>登録日時</h4>
						<p>{{.SignInDateTime}}</p>
					</div>

					<p>こちらはご登録ユーザーでサインインした際に通知されますので、ご自身で実行された場合は無視してください。</p>
				</div>
				<div class="footer">
					&copy; 2024 ファイナンスアプリ. All rights reserved.
				</div>
			</div>
		</body>
	</html>
`))

func TemporayPostSingUpTemplate(Name, ConfirmCode string) (string, string, error) {
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

func PostSingUpTemplate(Name, UserName, RegisterDateTime string) (string, string, error) {
	subject := "【ファイナンスアプリ】登録を完了致しました"
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := PostSingUpEmailData{
		Name:             Name,
		UserName:         UserName,
		RegisterDateTime: RegisterDateTime,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := postSingUpTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}

func PostSingInTemplate(UserName, SignInDateTime string) (string, string, error) {
	subject := "【ファイナンスアプリ】サインイン致しました"
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := PostSingInEmailData{
		UserName:       UserName,
		SignInDateTime: SignInDateTime,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := postSingInTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}
