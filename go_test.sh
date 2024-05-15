#!/bin/bash

echo "test start"

go clean -testcache

# 共通変数定義
out="coverage/coverage.out"
html="coverage/coverage.html"

if [ $# -eq 0 ]; then
    cd test
    go test ./... -coverprofile="../${out}"
    go tool cover -html="../${out}" -o "../${html}"
else
    input="$1"
    path="./${input//./\/}/"
    dot_count=$(grep -o "\." <<< "$input" | wc -l)
    coverage_path=""
    if [ -d "$path" ]; then
        cd $path
        for ((i=0; i<=dot_count; i++)); do
            coverage_path+="../"
        done
        go test . -coverprofile="${coverage_path}${out}"
        go tool cover -html="${coverage_path}${out}" -o "${coverage_path}${html}"
    else
        echo "そのようなディレクトリは存在しません: $path"
    fi
fi

echo "test end!"
