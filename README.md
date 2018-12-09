# Ingestor and storer services

This project contains two services: svc-ingestor and svc-storer.

## Build the project
```
make build
```

## Build the Docker images
```
make docker_images
```

## Run the project
```
make docker_run
```

## Information
The information and configurations for services are located in the respective
directories.

# Note
This project does not use `vgo` introduced in Go 1.11.0. Package dependencies
are managed using `glide`.
