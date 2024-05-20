# ベースイメージを指定
FROM ubuntu:22.04

# パッケージの更新と必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    && apt-get clean

# Go言語のインストール
RUN wget https://go.dev/dl/go1.20.linux-amd64.tar.gz \
    && tar -xvf go1.20.linux-amd64.tar.gz \
    && mv go /usr/local

# 環境変数の設定
ENV PATH="/usr/local/go/bin:${PATH}"

# open port 8080
EXPOSE 8080
