#!/usr/bin/env bash

docker -v

docker inspect rethinkdb > /dev/null
inspectRethinkdbExitCode=$?

if [ $inspectRethinkdbExitCode != 0 ]; then
    echo "Creating new rethinkdb container"
    docker pull rethinkdb:latest
	docker run --name rethinkdb -v "/data/rethink" -p 8080:8080 -p 28015:28015 -d rethinkdb
else
    echo "Starting rethinkdb"
    docker start rethinkdb
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

docker ps

echo "All done!"
