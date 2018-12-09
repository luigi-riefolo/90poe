# Ingestor service

The Ingestor reads a CSV file containing a list of user entries and sends them
to the Storer service.

All entries in the file are atomically processed and normalised (phone number
and email) before being sent to the downstream service Storer.

The service definition is declared in `pb/ingestor.proto`.

NOTE
The service starts sending the user entries as soon as it starts.
This implementation requires a synchronisation mechanism between svc-ingestor
and svc-storer, that checks whether svc-storer is able to accept requests.
This mechanism is currently not implemented.

A different approach is represented by an initiator endpoint, such as `IngestFileHandler` that can be exposed to the user in order to start the parsing process.
This endpoint is currently present but not implemented (NOP handler).

# Data file

The services can parse a list of user entries in CSV file.
All data files are stored in the `data` directory.

To specify a data file run:

```
INGESTOR_DATA_FILE=./data/data.csv go run cmd/main.go
```

or

```
INGESTOR_DATA_FILE=/data/data.csv make docker_run
```

## Run
```go run cmd/main.go```

or

```SERVICE_GRPC_PORT=10081 go run cmd/main.go```

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
