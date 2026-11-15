# Go command.
GO ?= go

# Tool versions.
GCL_VERSION  := v2.15.0
SWAG_VERSION := v1.16.4

# Project paths.
CMD_PATH    := ./cmd/image-randomizer
BINARY_NAME := image-randomizer
BINARY_PATH := $(CMD_PATH)/$(BINARY_NAME)

.PHONY: fmt clean lint swagger install-linter install-swagger build api migrate

##@ Tools

fmt: ## Standarize go coding style for the whole project.
	@$(GO) fmt ./...

clean: ## Clean project binary, test, and coverage file.
	@$(GO) clean ./...
	@rm -f $(COVER_OUT)

lint: fmt ## Lint go source code.
	@golangci-lint run -D errcheck --timeout 5m

swagger: ## Generate swagger docs.
	@swag init -g $(CMD_PATH)/main.go -o ./docs

install-linter: ## Install golangci-lint.
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GCL_VERSION)

install-swagger: ## Install swag.
	@$(GO) install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

##@ Build & Run

build: clean fmt ## Build the project executable binary. 
	@$(GO) build -o $(BINARY_PATH) -v $(CMD_PATH)

migrate: build ## Build and migrate database.
	@$(BINARY_PATH) migrate

api: build ## Build and run API server.
	@$(BINARY_PATH) server
