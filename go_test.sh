#!/bin/bash

echo "test start"

go clean -testcache

if [ $# -eq 0 ]; then
    cd test
    go test ./... -coverprofile=../coverage/coverage.out
    go tool cover -html=../coverage/coverage.out -o ../coverage/coverage.html
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
        go test . -coverprofile=${coverage_path}coverage/coverage.out
        go tool cover -html=${coverage_path}coverage/coverage.out -o ${coverage_path}coverage/coverage.html
    else
        echo "ディレクトリが存在しません: $path"
    fi
fi

echo "test end!"
