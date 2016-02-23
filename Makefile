default: get test

test:
	go test -v services/api/... services/realtime/...

integration:
	go test -v test/integration/

get:
	go get -t -v ./...

killapi:
	-fuser -k 4000/tcp

api: killapi
	go run src/services/api/main.go

killrealtime:
	-fuser -k 6000/tcp

realtime: killrealtime
	go run src/services/realtime/main.go

build:
	bash scripts/build.sh
