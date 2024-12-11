- 新規環境構築
    - 各OSに合わせてGoのインストーラをダウンロード URL: https://golang.org/dl/
	- コマンドプロンプト or ターミナルで 'go version' と入力して実行しGoのバージョンが表示されるか確認
	- 表示されたら git clone でソース取得して以下のコマンドでサーバー立ち上げる
	
```bash
mkdir server
cd ./server
go mod init server
go run main.go
go build -o server main.go
```

- git clone 後の環境構築
   - dockerコンテナ実行
   ```bash
   # new docker images command
   docker build -t server .
   docker container run -it -v ${home directory}/money_management/server:/server --name server server

   # windows
   docker container run -it -v ${home directory}\\money_management\\server:/server --name server server

   # again docker container
   docker start server
   docker container exec -it server /bin/bash

   ```

- test exec command
  ```bash
  cd <test dir>
  go test -coverprofile="../coverage/coverage.out"
  go tool cover -html=../coverage/coverage.out -o ../coverage/coverage.html
  ```

- mock作成
  - 対象ファイルにインターフェイスを定義
  - 以下のコマンドを実行
  ```bash
  // <>の中は適宜変える
  mockgen -source=<./controllers/sing_controllers.go> -destination=<./mock_func/mock_controllers/sing_controllers_mock.go>
  ```

  {
  "data": [
    {
      "user_name": "ma_kux@icloud.com",
      "user_password": "Test12345!",
      "nick_name": "ひでと"
    }
  ]
}

{
  "data": [
    {
      "user_name": "ma_kux@icloud.com",
      "user_password": "Test12345!"
    }
  ]
}