#!/bin/bash

input="$1"
output="$2"

echo "Cleaning up"
mkdir -p build
rm -fv build/${output}

export GOOS=linux
export CGO_ENABLED=0

echo "Building ${input} to ${output}"
go build -a -installsuffix cgo -o ${output} ${input}