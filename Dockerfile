# ベースイメージを指定
FROM ubuntu:22.04

# パッケージの更新と必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    iproute2 \
    && apt-get clean

# アーキテクチャを判別してGo言語のインストール
ARG TARGETARCH
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      wget https://go.dev/dl/go1.20.linux-amd64.tar.gz; \
    else \
      wget https://go.dev/dl/go1.20.linux-arm64.tar.gz; \
    fi \
    && tar -xvf go1.20.linux-*.tar.gz \
    && mv go /usr/local

# 環境変数の設定
ENV PATH="/usr/local/go/bin:${PATH}"

# 必要なパッケージをインストールし、依存関係を取得
RUN go get github.com/agiledragon/gomonkey/v2

# open port 8080
EXPOSE 8080
