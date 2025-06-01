.PHONY: all build run test clean tidy help

all: build

build: ## Build the Go application
	@echo "Building Go application..."
	go build -o feed-generator .

run: build ## Run the feed generator application
	@echo "Running feed generator..."
	@echo "NOTE: For full functionality, ensure API keys are set as environment variables (e.g., export LINKEDIN_API_KEY='your_key')."
	./feed-generator

test: ## Run all Go tests
	@echo "Running tests..."
	go test ./...

clean: ## Clean up compiled binaries and output files
	@echo "Cleaning up..."
	rm -f feed-generator
	rm -rf output/

tidy: ## Run go mod tidy to synchronize module dependencies
	@echo "Tidying Go modules..."
	go mod tidy

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
