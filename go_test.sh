#!/bin/bash

echo "Starting tests with coverage"

go clean -testcache

mkdir -p coverage

# すべてのテストファイルを含むディレクトリを検索し、リストに保存
test_dirs=$(find . -type f -name '*_test.go' -exec dirname {} \; | sort -u)

# 各ディレクトリでテストを実行し、カバレッジプロファイルを生成
for dir in $test_dirs; do
    echo "Running tests in $dir"
    go test -coverprofile=coverage/$(echo $dir | tr '/' '_').out ./$dir
done

# すべてのカバレッジプロファイルを結合
echo "mode: set" > coverage/coverage.out
for file in coverage/*.out; do
    if [ $file != "coverage/coverage.out" ]; then
        tail -n +2 $file >> coverage/coverage.out
    fi
done

# カバレッジレポートを生成
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

echo "Coverage report generated at coverage/coverage.html"
