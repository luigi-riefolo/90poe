include .env
export

.DEFAULT_GOAL := all

COMMIT_HASH := $(shell git rev-parse --short HEAD)
VERSION := $(COMMIT_HASH)

ifndef VERSION
	VERSION := undefined
endif

ifndef STORER_DATA_FILE
	STORER_DATA_FILE=/data/data.csv
endif

IMAGE_NAME := svc-${SERVICE_NAME}:${VERSION}
SERVICE_VERSION := ${VERSION}


## Install dependencies and build Docker image
all: deps unit_tests integration_tests build docker_image


## Install dependencies
deps:
	$(info Installing dependencies)
	@glide --no-color --debug install


## Build the project binary
build:
	$(info Building project)
	@CGO_ENABLED=0 GOOS=linux go build -a \
		-ldflags "-X main.version=${VERSION}" \
		-installsuffix cgo -o ${SERVICE_NAME} \
		./cmd/main.go


## Generate the gRPC code 
protobuf:
	$(info Creating Go protobuf)
	@protoc -I=pb -I=${GOPATH}/src \
		-I=${GOPATH}/src/github.com/google/protobuf/src/ \
		--go_out=plugins=grpc:pb/ \
		pb/*.proto
	@protoc-go-inject-tag -input=./pb/storer.pb.go
	@mockery -all -dir pb -output pb/mocks -case underscore


## Create the Docker image
docker_image:
	$(info Creating Docker image)
	@docker-compose build --force-rm --no-cache
	@docker tag ${IMAGE_NAME} svc-${SERVICE_NAME}:latest


## Run Docker image
docker_run:
	$(info Starting Docker image)
	@docker-compose up


## Unit test the Go code
unit_tests:
	$(info Running unit tests)
	@go test -v -race -cover -failfast ./api


## Display this help screen
help:
	@gawk 'match($$0, /^## (.*)/, a) \
		{ getline x; x = gensub(/(.+:) .+/, "\\1", "g", x) ; \
		printf "\033[36m%-30s\033[0m %s\n", x, a[1]; }' $(MAKEFILE_LIST) | sort


.PHONY: all build deps docker_image \
	docker_run help protobuf unit_tests
