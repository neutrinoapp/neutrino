#!/usr/bin/env bash

docker -v

docker inspect mongodb > /dev/null
inspectMongoDbExitCode=$?

if [ $inspectMongoDbExitCode != 0 ]; then
    echo "Starting new mongodb container"
    docker pull mongo:latest
    docker run --name mongodb -p 27017:27017 -v /data/mongo:/data/db -d mongo:latest --storageEngine wiredTiger
fi

docker inspect redis-cache > /dev/null
inspectRedisExitCode=$?

if [ $inspectRedisExitCode != 0 ]; then
    echo "Starting new redis container"
    docker pull redis:latest
    docker run --name redis-cache -p 6379:6379 -d redis redis-server --appendonly yes
fi

docker ps

echo "All done!"
