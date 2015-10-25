test: get
	go test -v ./api-service/... ./queue-broker-service/... ./realtime-service/...

integration: get
	go test -v ./integration-tests/

get:
	go get -t -v ./...

