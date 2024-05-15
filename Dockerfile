# ベースイメージとしてGoの公式イメージを使用
FROM golang:1.20

# パッケージの更新と必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    && apt-get clean
