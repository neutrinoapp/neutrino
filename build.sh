#!/bin/bash

echo "Cleaning up"
mkdir -p build
rm -rfv build/*

echo "Building api-service"
go build -o build/api-service api-service/main.go

echo "Building realtime-service"
go build -o build/realtime-service realtime-service/main.go

echo "Building queue-broker-service"
go build -o build/queue-broker-service queue-broker-service/main.go
