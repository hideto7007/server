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

	PostSingEmailData struct {
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
	&copy; {{.Year}} ファイナンスアプリ. All rights reserved.
</div>
{{end}}
`))

var temporayPostSingUpTemplate = template.Must(template.New("auth_email").Parse(`
{{.Name}}さん！

ファイナンスアプリをご利用いただきありがとうございます。


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
					ファイナンスアプリ<br>
					ご登録完了のお知らせ
				</div>
				<div class="body">
					<h4>{{.Name}}さん！</h4>
					<p>ファイナンスアプリのご登録が完了しました。誠にありがとうございます。</p>
					<p>{{.Name}}さんにとって、より良い資産管理が出来ることを祈っております。</p>
					</br>

					<div class="info-section">
						<h4>登録ユーザ名</h4>
						<p>{{.UserName}}</p>
						<h4>登録日時</h4>
						<p>{{.DateTime}}</p>
					</div>
					<p>今後ともよろしくお願いいたします。</p>
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
						<p>{{.DateTime}}</p>
					</div>
					<p>こちらはご登録ユーザーでサインインした際に通知されますので、ご自身で実行された場合は無視してください。</p>
				</div>
				{{template "Footer" .}}
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

func PostSingUpTemplate(Name, UserName, DateTime string) (string, string, error) {
	subject := "【ファイナンスアプリ】登録を完了致しました"
	// メールテンプレート定義

	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")

	// テンプレートに渡すデータを作成
	data := PostSingEmailData{
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
	subject := "【ファイナンスアプリ】サインイン致しました"
	var year = utils.NewUtilsFetcher(utils.JwtSecret).DateTimeStr(time.Now(), "2006年")
	// メールテンプレート定義

	// テンプレートに渡すデータを作成
	data := PostSingEmailData{
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
