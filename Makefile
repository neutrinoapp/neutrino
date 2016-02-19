default: get test

test:
	go test -v ./api-service/... ./queue-broker-service/... ./realtime-service/...

integration:
	go test -v ./integration-tests/

get:
	go get -t -v ./...

kill:
	-fuser -k 4000/tcp 5000/tcp 6000/tcp

all: kill
	bash start-all.sh

build:
	bash build.sh