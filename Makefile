default: get integration

integration:
	go test -v ./test/

get:
	go get -t -v ./...

killapi:
	-fuser -k 5000/tcp

api: killapi
	go run src/services/api/main.go

killrealtime:
	-fuser -k 6000/tcp

realtime: killrealtime
	go run src/services/realtime/main.go

build:
	bash scripts/build.sh

dev:
	bash scripts/dev-setup.sh
