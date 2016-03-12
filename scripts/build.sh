#!/bin/bash

echo "Cleaning up"
mkdir -p build
rm -rfv build/*

export GOOS=linux
export CGO_ENABLED=0

echo "Building api service"
go build -a -installsuffix cgo -o build/api src/services/api/main.go

echo "Building realtime service"
go build -a -installsuffix cgo -o build/realtime src/services/realtime/main.go