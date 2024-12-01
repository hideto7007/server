package templates

import (
	"bytes"
	"server/utils"
	"text/template"
	"time"
)

type (
	// データ構造体
	TemporayPostSingUpEmailData struct {
		Name        string
		ConfirmCode string
	}

	SingEmailData struct {
		Style    string
		Name     string
		UserName string
		DateTime string
		Footer   string
		Year     string
	}
)

var commonTemplate = template.Must(template.New("common").Parse(`
{{define "Style"}}
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
{{end}}

{{define "Footer"}}
<div class="footer">
	&copy; {{.Year}} たくわえる. All rights reserved.
</div>
{{end}}

{{define "Support"}}
<div style="padding: 20px; background-color: #e8f5e9; border-radius: 5px; margin-top: 20px; text-align: center;">
	<hr style="border: 1px solid #ddd; margin: 20px 0;">
	<p style="margin: 0; font-size: 14px; color: #666;">ご不明点やご質問がございましたら、以下のサポート窓口までお問い合わせください。</p>
	<div style="margin-top: 10px;">
		<h5 style="margin: 0; font-size: 16px; color: #333;">カスタマーサポート</h5>
		<p style="margin: 5px 0; font-size: 14px; color: #007bff;">
			Email: <a href="mailto:finance.1128.th@gmail.com" style="text-decoration: none; color: #007bff;">finance.1128.th@gmail.com</a>
		</p>
	</div>
	<hr style="border: 1px solid #ddd; margin: 20px 0;">
</div>
{{end}}
`))

var temporayPostSingUpTemplate = template.Must(template.New("auth_email").Parse(`
{{.Name}}さん！

たくわえるをご利用いただきありがとうございます。


確認コードは {{.ConfirmCode}} です。

この確認コードの有効期限は1時間です。
`))

var postSingUpTemplate = template.Must(template.Must(commonTemplate.Clone()).Parse(`
	{{template "Style"}}
	<!DOCTYPE html>
	<html>
		<head>
			<title>サインアップ完了のお知らせ</title>
			{{.Style}}
		</head>
		<body>
			<div class="container">
				<div class="header">
					たくわえる<br>
					ご登録完了のお知らせ
				</div>
				<div class="body">
					<h4>{{.Name}}さん！</h4>
					<p>たくわえるのご登録が完了しました。誠にありがとうございます。</p>
					<p>{{.Name}}さんにとって、より良い資産管理が出来ることを祈っております。</p>
					</br>

					<div class="info-section">
						<h4>登録ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>登録日時</h4>
						<p>{{.DateTime}}</p>
					</div>
					<p>今後ともよろしくお願いいたします。</p>
					{{template "Support"}}
				</div>
				{{template "Footer" .}}
			</div>
		</body>
	</html>
`))

var postSingInTemplate = template.Must(template.Must(commonTemplate.Clone()).Parse(`
	{{template "Style"}}
	<!DOCTYPE html>
	<html>
		<head>
			<title>サインイン通知</title>
		</head>
		<body>
			<div class="container">
				<div class="header">
					たくわえる<br>
					Webブラウザから新たなサインインがありました
				</div>
				<div class="body">
					<p>いつもたくわえるをご利用いただき、誠にありがとうございます。</p>
					<p>お客様がご利用中の登録ユーザーで、新たなサインインがありました。</p>

					<div class="info-section">
						<h4>登録ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>実行日時</h4>
						<p>{{.DateTime}}</p>
					</div>
					<p>こちらはご登録ユーザーでサインインした際に通知されますので、ご自身で実行された場合は無視してください。</p>
					{{template "Support"}}
				</div>
				{{template "Footer" .}}
			</div>
		</body>
	</html>
`))

var deleteSingInTemplate = template.Must(template.Must(commonTemplate.Clone()).Parse(`
	{{template "Style"}}
	<!DOCTYPE html>
	<html>
		<head>
			<title>削除完了通知</title>
		</head>
		<body>
			<div class="container">
				<div class="header">
					たくわえる<br>
					アカウント削除完了のお知らせ
				</div>
				<div class="body">
					<p>{{.UserName}}さん、この度はたくわえるをご利用いただき、誠にありがとうございました。</p>
					<p>以下の内容でアカウントの削除が完了しました。</p>

					<div class="info-section">
						<h4>削除ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>削除日時</h4>
						<p>{{.DateTime}}</p>
					</div>

					<p>アカウント削除に伴い、関連するすべてのデータが安全に削除されたことをお知らせいたします。</p>
					<p>また、いつでもたくわえるをご利用いただけるよう準備しておりますので、再度のご利用を心よりお待ちしております。</p>
					{{template "Support"}}
				</div>
				{{template "Footer" .}}
			</div>
		</body>
	</html>
`))

var deleteSignOutTemplate = template.Must(template.Must(commonTemplate.Clone()).Parse(`
	{{template "Style"}}
	<!DOCTYPE html>
	<html>
		<head>
			<title>サインアウト通知</title>
		</head>
		<body>
			<div class="container">
				<div class="header">
					たくわえる<br>
					Webブラウザからサインアウトされました
				</div>
				<div class="body">
					<p>いつもたくわえるをご利用いただき、誠にありがとうございます。</p>
					<p>お客様がご利用中の登録ユーザーで、サインアウトがありました。</p>

					<div class="info-section">
						<h4>登録ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>実行日時</h4>
						<p>{{.DateTime}}</p>
					</div>
					<p>こちらはご登録ユーザーでサインアウトした際に通知されますので、ご自身で実行された場合は無視してください。</p>
					{{template "Support"}}
				</div>
				{{template "Footer" .}}
			</div>
		</body>
	</html>
`))

func TemporayPostSingUpTemplate(Name, ConfirmCode string) (string, string, error) {
	subject := "【たくわえる】本登録を完了してください"
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

func PostSingUpTemplate(Name, UserName, DateTime string) (string, string, error) {
	subject := "【たくわえる】登録を完了致しました"
	// メールテンプレート定義

	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")

	// テンプレートに渡すデータを作成
	data := SingEmailData{
		Name:     Name,
		UserName: UserName,
		DateTime: DateTime,
		Year:     year,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := postSingUpTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}

func PostSingInTemplate(UserName, DateTime string) (string, string, error) {
	subject := "【たくわえる】サインイン致しました"
	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := SingEmailData{
		UserName: UserName,
		DateTime: DateTime,
		Year:     year,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := postSingInTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}

func DeleteSingInTemplate(UserName, DateTime string) (string, string, error) {
	subject := "【たくわえる】アカウント削除完了のお知らせ"
	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := SingEmailData{
		UserName: UserName,
		DateTime: DateTime,
		Year:     year,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := deleteSingInTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}

func SignOutTemplate(UserName, DateTime string) (string, string, error) {
	subject := "【たくわえる】サインアウトのお知らせ"
	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := SingEmailData{
		UserName: UserName,
		DateTime: DateTime,
		Year:     year,
	}

	// テンプレートの実行と結果の取得
	var body bytes.Buffer
	if err := deleteSignOutTemplate.Execute(&body, data); err != nil {
		return "", "", err // エラー時に空の件名と本文を返す
	}

	return subject, body.String(), nil
}
