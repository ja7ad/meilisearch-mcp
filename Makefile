build:
	echo "Building the project..."
	go build -o build/meilisearch-mcp cmd/mcp/main.go

build-docker:
	echo "Building the Docker image..."
	docker build -t meilisearch-mcp:latest .

tests:
	@echo "Running unit tests..."
	go test -v ./...

check-all:
	@echo "running all checks..."
	$(MAKE) fmt
	$(MAKE) check

fmt:
	@echo "formatting code..."
	go tool gofumpt -l -w .

check:
	@echo "linting..."
	go tool golangci-lint run --timeout=20m0s --tests=false -v
