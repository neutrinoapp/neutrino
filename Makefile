#takes the first 8 letters of a commit, e.g.
#de1a3d902db01bf39c83020db1207ab1a3760750 -> de1a3d90

BUILD_NUMBER := $(shell git rev-parse HEAD|rev|cut -c33-|rev)

build-number:
	echo $(BUILD_NUMBER)

default: get integration

integration:
	cd ./src/tests
	go get -t -d ./...
	cd ../..
	DEBUG_N=true go test -v ./src/tests/

prep:
	DEBUG_N=true go run scripts/prepareRethinkDb/main.go

get:
	go get github.com/Masterminds/glide
	glide install

killapi:
	-fuser -k 5000/tcp

api: killapi
	DEBUG_N=true SERVICE_NAME=api go run src/services/api/main.go

killrealtime:
	-fuser -k 6000/tcp

realtime: killrealtime
	DEBUG_N=true SERVICE_NAME=realtime go run src/services/realtime/main.go

build-api:
	bash scripts/build.sh src/services/api/main.go build/api

build-realtime:
	bash scripts/build.sh src/services/realtime/main.go build/realtime

build-docker-api: build-api
	docker build -f api-dockerfile -t gcr.io/neutrino-1073/api-service:$(BUILD_NUMBER) .

build-docker-realtime: build-realtime
	docker build -f realtime-dockerfile -t gcr.io/neutrino-1073/realtime-service:$(BUILD_NUMBER) .

build-services: build-api build-realtime

build-docker: build-services build-docker-api build-docker-realtime
	gcloud docker push gcr.io/neutrino-1073/realtime-service:$(BUILD_NUMBER)
	gcloud docker push gcr.io/neutrino-1073/api-service:$(BUILD_NUMBER)

	docker tag gcr.io/neutrino-1073/realtime-service:$(BUILD_NUMBER) gcr.io/neutrino-1073/realtime-service:latest
	docker tag gcr.io/neutrino-1073/api-service:$(BUILD_NUMBER) gcr.io/neutrino-1073/api-service:latest
	gcloud docker push gcr.io/neutrino-1073/realtime-service:latest
	gcloud docker push gcr.io/neutrino-1073/api-service:latest

dev:
	bash scripts/dev-setup.sh

build-rethink-prepare:
	bash scripts/build.sh scripts/prepareRethinkDb/main.go build/prepare

build-rethink-prepare-docker: build-rethink-prepare
	docker build -f scripts/prepareRethinkDb/Dockerfile -t gcr.io/neutrino-1073/prepare-rethinkdb:v1 .

build-rethinkdb:
	docker build -f scripts/rethinkdb-next/Dockerfile -t gcr.io/neutrino-1073/realtime-service/rethinkdb-next:v1 .