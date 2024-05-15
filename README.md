- 新規環境構築
    - 各OSに合わせてGoのインストーラをダウンロード URL: https://golang.org/dl/
	- コマンドプロンプト or ターミナルで 'go version' と入力して実行しGoのバージョンが表示されるか確認
	- 表示されたら git clone でソース取得して以下のコマンドでサーバー立ち上げる
	
```bash
mkdir server
cd ./server
go mod init server
go run main.go
```

- git clone 後の環境構築
   - dockerコンテナ実行
   ```bash
   # new docker images command
   docker build -t server .
   docker container run -it -v ${home directory}/money_management/server/:/server --name server server

   # windows
   docker container run -it -v ${home directory}\\money_management\\server:/server --name server server

   # again docker container
   docker start server
   docker container exec -it server /bin/bash

   ```


