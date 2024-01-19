- 環境構築
    - 各OSに合わせてGoのインストーラをダウンロード URL: https://golang.org/dl/
	- コマンドプロンプト or ターミナルで 'go version' と入力して実行しGoのバージョンが表示されるか確認
	- 表示されたら git clone でソース取得して以下のコマンドでサーバー立ち上げる
	
```bash
mkdir server
cd ./server
go mod init server
go run main.go
```


