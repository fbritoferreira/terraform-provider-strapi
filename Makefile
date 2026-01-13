.PHONY: build clean test lint generate help

# Variables
BINARY_NAME=terraform-provider-strapi
MAIN_FILE=main.go

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the provider binary
	go build -o $(BINARY_NAME) $(MAIN_FILE)

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf bin/

test: ## Run tests
	go test -v ./...

testacc: ## Run acceptance tests (requires STRAPI_ENDPOINT and STRAPI_API_TOKEN to be set)
	TF_ACC=1 go test -v -cover ./internal/provider/

sandbox-start: ## Start Strapi sandbox and create admin user + API token
	@./scripts/test-setup.sh start

sandbox-stop: ## Stop Strapi sandbox
	@./scripts/test-setup.sh stop

sandbox-test: sandbox-start ## Start sandbox, run acceptance tests, then stop sandbox
	@echo "Running acceptance tests..."
	@STRAPI_ENDPOINT=http://localhost:1337 \
		STRAPI_API_TOKEN=$$(cat .strapi-test-token) \
		TF_ACC=1 go test -v -cover ./internal/provider/ || (./scripts/test-setup.sh stop && exit 1)
	@./scripts/test-setup.sh stop
	@echo "Tests completed!"

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run

generate: ## Generate documentation
	@echo "Generating documentation..."
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name strapi

fmt: ## Format code
	go fmt ./...

go-mod-tidy: ## Run go mod tidy
	go mod tidy

go-mod-download: ## Download dependencies
	go mod download

install: ## Install the provider locally
	@echo "Installing provider locally..."
	@mkdir -p $$HOME/.terraform.d/plugins/registry.terraform.io/fbritoferreira/strapi/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	@cp $(BINARY_NAME) $$HOME/.terraform.d/plugins/registry.terraform.io/fbritoferreira/strapi/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/$(BINARY_NAME)
	@echo "Provider installed to $$HOME/.terraform.d/plugins/registry.terraform.io/fbritoferreira/strapi/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/$(BINARY_NAME)"

dev: build ## Build and install provider for development
	@echo "Building and installing for development..."
	@make build
	@make install

.DEFAULT_GOAL := help
