build:
	echo "Building the project..."
	go build -o build/meilisearch-mcp cmd/mcp/main.go

build-docker:
	echo "Building the Docker image..."
	docker build -t meilisearch-mcp:latest .

dev-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0
	go install mvdan.cc/gofumpt@latest

tests:
	@echo "Running unit tests..."
	go test -v ./...

check-all:
	@echo "running all checks..."
	$(MAKE) fmt
	$(MAKE) check

fmt:
	@echo "formatting code..."
	gofumpt -l -w .

check:
	@echo "linting..."
	golangci-lint run --timeout=20m0s --tests=false -v
