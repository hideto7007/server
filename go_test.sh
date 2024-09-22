#!/bin/bash

echo "Starting tests with coverage"

mkdir -p coverage

go clean -testcache

# すべてのテストファイルを含むディレクトリを検索し、リストに保存
TEST_DIRS=$(find . -type f -name '*_test.go' -exec dirname {} \; | sort -u)

echo "mode: set" > coverage/coverage.out

# 各ディレクトリでテストを実行し、カバレッジプロファイルを生成
for dir in $TEST_DIRS; do
    echo "Running tests in $dir"
    ENV=test go test -coverprofile=coverage/$(echo $dir | tr '/' '_').out ./$dir -count=1

    # カバレッジファイルの結合
    OUT_FILE="coverage/$(echo $dir | tr '/' '_').out"
    if [ -f "$OUT_FILE" ]; then
        tail -n +2 "$OUT_FILE" >> coverage/coverage.out
        rm -f "$OUT_FILE" # ファイルを削除
    else
        echo "$OUT_FILE not found."
    fi
done

# カバレッジレポートを生成
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# coverage.outを削除する場合
rm -f coverage/coverage.out

echo "Coverage report generated at coverage/coverage.html"