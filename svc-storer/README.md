# Storeer service

The Storer service receives user entries and saves them in a store.

The service definition is declared in `pb/storer.proto`.


## Run
```go run cmd/main.go```

## Test
```make unit_tests```

Or

```make integration_test```

## Metrics

All service metrices are served at: `http://localhost:10060/metrics`


## Make targets
```
all:                           Install dependencies and build Docker image
build:                         Build the project binary
coverage:                      Show code Go code coverage
deps:                          Install dependencies
docker_image:                  Create the Docker image
docker_run:                    Run Docker image
help:                          Display this help screen
protobuf:                      generate the gRPC code 
unit_tests:                    Unit test the Go code
```
