build:
	echo "Building the project..."
	go build -o build/meilisearch-mcp cmd/mcp/main.go

build-docker:
	echo "Building the Docker image..."
	docker build -t meilisearch-mcp:latest .