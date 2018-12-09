# Makefile used for running the top-level Docker Compose file.

.DEFAULT_GOAL := all

# NOTE:
# This Makefile load a list of .env files from the respective services.
# The above approach is not considered to be maintainable nor scalable.
# The author chose this implementation to speed up the development.
# A more suitable and dynamic implementation would be to create the same list of
# environment variables per service, by looping through all the .env files
# and prepending the service name to each one, e.g.:
# SERVICE_GRPC_PORT=10080 --> INGESTOR_SERVICE_GRPC_PORT=10081


# load the environment variables for svc-ingestor
include svc-ingestor/.env
export

ifndef VERSION
	VERSION := undefined
endif

ifndef INGESTOR_DATA_FILE
	INGESTOR_DATA_FILE=/data/data.csv
endif

INGESTOR_SERVICE_NAME := ${SERVICE_NAME}
INGESTOR_IMAGE_NAME := svc-${SERVICE_NAME}:latest
INGESTOR_SERVICE_VERSION := ${VERSION}
INGESTOR_SERVICE_GRPC_PORT := ${SERVICE_GRPC_PORT}

# load the environment variables for svc-storer
include svc-storer/.env
export
ifndef VERSION
	VERSION := undefined
endif
STORER_SERVICE_NAME := ${SERVICE_NAME}
STORER_IMAGE_NAME := svc-${SERVICE_NAME}:latest
STORER_SERVICE_VERSION := ${VERSION}
STORER_SERVICE_GRPC_PORT := ${SERVICE_GRPC_PORT}
STORER_SERVICE_GRPC_HOST := ${SERVICE_NAME}


## Build and run the project
all: build docker_images docker_run


## Build all the services binaries
build:
	$(MAKE) -C svc-ingestor deps build
	$(MAKE) -C svc-storer deps build


## Create the Docker images
docker_images:
	$(info Creating Docker image)
	$(MAKE) -C svc-ingestor docker_image
	$(MAKE) -C svc-storer docker_image


## Run Docker images
docker_run:
	$(info Starting Docker image)
	@docker-compose up --remove-orphans --abort-on-container-exit


## Display this help screen
help:
	@gawk 'match($$0, /^## (.*)/, a) \
		{ getline x; x = gensub(/(.+:) .+/, "\\1", "g", x) ; \
		printf "\033[36m%-30s\033[0m %s\n", x, a[1]; }' $(MAKEFILE_LIST) | sort


.PHONY: all docker_images docker_run help
