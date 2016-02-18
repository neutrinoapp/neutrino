#!/usr/bin/env bash

echo {'api-service/main.go','realtime-service/main.go','queue-broker-service/main.go'} | xargs -n 1 -P 3 go run