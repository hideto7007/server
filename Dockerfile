# ベースイメージを指定
FROM ubuntu:20.04

# パッケージの更新と必要なパッケージのインストール
RUN apt-get update && apt-get install -y \
    curl \
    git \
    wget \
    && apt-get clean

# open port 8080
EXPOSE 8080
