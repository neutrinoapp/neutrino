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

docker inspect nats > /dev/null
inspectNatsExitCode=$?

if [ $inspectNatsExitCode != 0 ]; then
    echo "Creating new nats container"
    docker pull apcera/gnatsd:latest
    docker run --name nats -p 4222:4222 -p 8222:8222 -d nats:latest
else
    echo "Starting nats"
    docker start nats
fi

docker ps

echo "All done!"
