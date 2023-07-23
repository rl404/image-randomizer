# Base Go commands.
GO_CMD     := go
GO_FMT     := $(GO_CMD) fmt
GO_GET     := $(GO_CMD) get
GO_INSTALL := $(GO_CMD) install
GO_MOD     := $(GO_CMD) mod
GO_CLEAN   := $(GO_CMD) clean
GO_BUILD   := $(GO_CMD) build

# Base swagger commands.
SWAG     := swag
SWAG_GEN := $(SWAG) init

# Base golangci-lint commands.
GCL_CMD := golangci-lint
GCL_RUN := $(GCL_CMD) run

# Project executable file, and its binary.
CMD_PATH    := ./cmd/image-randomizer
BINARY_NAME := image-randomizer

# Default makefile target.
.DEFAULT_GOAL := run

# Standarize go coding style for the whole project.
.PHONY: fmt
fmt:
	@$(GO_FMT) ./...

# Lint go source code.
.PHONY: lint
lint: fmt
	@$(GCL_RUN) -D errcheck --timeout 5m

# Clean project binary, test, and coverage file.
.PHONY: clean
clean:
	@$(GO_CLEAN) ./...

# Generate swagger docs.
.PHONY: swagger
swagger:
	@$(SWAG_GEN) -g cmd/image-randomizer/main.go -o ./docs

# Install library.
.PHONY: install
install:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.49.0
	@$(GCL_CMD) version
	@$(GO_INSTALL) github.com/swaggo/swag/cmd/swag@latest
	@$(SWAG) -v

# Build the project executable binary.
.PHONY: build
build: clean fmt
	@cd $(CMD_PATH); \
	$(GO_BUILD) -o $(BINARY_NAME) -v .

# Build and migrate database.
.PHONY: migrate
migrate: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) migrate

# Build and run the binary.
.PHONY: run
run: build
	@cd $(CMD_PATH); \
	./$(BINARY_NAME) server


# Docker base command.
DOCKER_CMD   := docker
DOCKER_IMAGE := $(DOCKER_CMD) image

# Docker-compose base command and docker-compose.yml path.
COMPOSE_CMD     := docker-compose
COMPOSE_BUILD   := deployment/build.yml
COMPOSE_API     := deployment/api.yml
COMPOSE_MIGRATE := deployment/migrate.yml
COMPOSE_LINT    := deployment/lint.yml

# Build docker images and container for the project
# then delete builder image.
.PHONY: docker-build
docker-build:
	@$(COMPOSE_CMD) -f $(COMPOSE_BUILD) build
	@$(DOCKER_IMAGE) prune -f --filter label=stage=image_randomizer_be_builder

# Start built docker containers for api.
.PHONY: docker-api
docker-api:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p image-randomizer-be-api up -d
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p image-randomizer-be-api logs --follow --tail 20

# Start built docker containers for migrate.
.PHONY: docker-migrate
docker-migrate:
	@$(COMPOSE_CMD) -f $(COMPOSE_MIGRATE) -p image-randomizer-be-migrate up

# Start docker to run lint check.
.PHONY: docker-lint
docker-lint:
	@$(COMPOSE_CMD) -f $(COMPOSE_LINT) -p image-randomizer-be-lint run --rm image-randomizer-be-lint $(GCL_RUN) -D errcheck --timeout 5m

# Update docker containers.
.PHONY: docker-update
docker-update:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p image-randomizer-be-api up -d
	@$(DOCKER_IMAGE) prune -f --filter label=stage=image-randomizer-be_binary

# Stop running docker containers.
.PHONY: docker-stop
docker-stop:
	@$(COMPOSE_CMD) -f $(COMPOSE_API) -p image-randomizer-be-api stop