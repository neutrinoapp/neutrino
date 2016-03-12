default: get integration

integration:
	go test -v ./src/tests/

prep:
	go run scripts/prepareRethinkDb/main.go

get:
	go get -t -v -d ./src/...

killapi:
	-fuser -k 5000/tcp

api: killapi
	go run src/services/api/main.go

killrealtime:
	-fuser -k 6000/tcp

realtime: killrealtime
	go run src/services/realtime/main.go

build-docker-api:
	docker build -f api-dockerfile -t api-service .

build-docker-realtime:
	docker build -f realtime-dockerfile -t realtime-service .

build-docker: build-services build-docker-api build-docker-realtime

build-services:
	bash scripts/build.sh

dev:
	bash scripts/dev-setup.sh
