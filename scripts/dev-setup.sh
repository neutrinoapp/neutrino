#!/usr/bin/env bash

docker -v

docker inspect mongodb > /dev/null
inspectMongoDbExitCode=$?

if [ $inspectMongoDbExitCode != 0 ]; then
    echo "Creating new mongodb container"
    docker pull mongo:latest
    docker run --name mongodb -p 27017:27017 -v /data/mongo:/data/db -d mongo:latest --storageEngine wiredTiger
else
    echo "Starting mongodb"
    docker start mongodb
fi

docker inspect redis-cache > /dev/null
inspectRedisExitCode=$?

if [ $inspectRedisExitCode != 0 ]; then
    echo "Creating new redis container"
    docker pull redis:latest
    docker run --name redis-cache -p 6379:6379 -d redis redis-server --appendonly yes
else
    echo "Starting redis-cache"
    docker start redis-cache
fi

docker inspect nats > /dev/null
inspectNatsExitCode=$?

if [ $inspectNatsExitCode != 0 ]; then
    echo "Creating new nats container"
    docker pull apcera/gnatsd:latest
    docker run --name nats -p 4222:4222 -p 8333:8333 -d apcera/gnatsd:latest -m 8333
else
    echo "Starting nats"
    docker start nats
fi

docker ps

echo "All done!"
