# Variables
BINARY_NAME=main
BINARY_PATH=main.go

ifeq ($(OS),Windows_NT)
    BINARY_OUT=$(BINARY_NAME).exe
    RM=if exist $(BINARY_OUT) del /Q $(BINARY_OUT)
    HELP_CMD=powershell -NoProfile -ExecutionPolicy Bypass -File ./help.ps1
else
    BINARY_OUT=$(BINARY_NAME)
    RM=rm -f $(BINARY_OUT)
    HELP_CMD=grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-15s %s\n", $$1, $$2}'
endif

.PHONY: run test clean build push up help

help: ## Show help messages
	@echo Available commands:
	@$(HELP_CMD)

run: ## Run the application
	@echo Running the application...
	go run $(BINARY_PATH)

test: ## Run tests
	@echo Testing the application...
	go test ./...

build: ## Build the application
	@echo Building the application...
	go build -o ${BINARY_OUT} ${BINARY_PATH}

clean: ## Delete the build file
	@echo Cleaning the application...
	go clean
	@$(RM)

migration: ## migrate db
	@echo Migrating db
	go run ./migration/cmd/main.go

docker-build: ## Build docker image
	@echo Building the application...
	docker build -t janghanul090801/spine-clean-architecture:latest .

docker-push: ## Push docker image
	@echo Pushing the docker image...
	docker push janghanul090801/spine-clean-architecture:latest

compose-up: ## Up docker-compose
	@echo Upping docker compose...
	docker-compose up -d

swag: ## generate docs
	swag init --parseDependency --parseInternal --parseDepth 1
