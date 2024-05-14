# ベースイメージとしてGoの公式イメージを使用
FROM golang:1.20

# パッケージの更新と必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    && apt-get clean

# 依存関係をダウンロード
RUN go mod download

# 依存関係初期化
RUN go mod init server
