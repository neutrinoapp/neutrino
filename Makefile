default: get test

test:
	go test -v ./api-service/... ./queue-broker-service/... ./realtime-service/...

integration:
	go test -v ./integration-tests/

get:
	go get -t -v ./...

all:
	bash start-all.sh
