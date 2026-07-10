# Meilisearch MCP

A high-performance Go implementation of the Model Context Protocol (MCP) for Meilisearch. This server enables AI assistants and desktop environments (like Claude Desktop, Cursor, Zed, etc.) to connect directly with your Meilisearch database to query, manage, and inspect search schemas, documents, and settings.

<p align="center">
  <img src="docs/mcp.gif" alt="Meilisearch MCP demo" />
</p>


## Table of Contents
- [Features](#features)
- [Quick Start](#quick-start)
    - [Connect to Remote MCP Server](#1-connect-to-remote-mcp-server)
    - [Download Latest Release](#2-download-latest-release)
    - [Build from Source](#3-build-from-source)
        - [Prerequisites](#prerequisites)
        - [Build](#build)
        - [Run (HTTP Transport)](#run-http-transport)
        - [Run (STDIO Transport)](#run-stdio-transport)
- [Available MCP Capabilities](#available-mcp-capabilities)
    - [1. Tools](#1-tools)
    - [2. Resources](#2-resources)
    - [3. Prompts](#3-prompts)
- [Docker](#docker)
    - [Build Image](#build-image)
    - [Run Container](#run-container)
    - [Environment Variables](#environment-variables)
    - [docker-compose Example](#docker-compose-example)
- [Configuration Examples](#configuration-examples)
    - [Generic HTTP (mcp-remote)](#generic-http-mcp-remote)
    - [Local STDIO](#local-stdio)
- [Supported Clients](#supported-clients)
- [Security Notes](#security-notes)
- [Contributing](#contributing)
- [License](#license)


## Features
- **Comprehensive Meilisearch Integration**: Supports all features of the `meilisearch-go` SDK (v0.36.3) including index administration, settings, keys, tasks, custom webhooks, experimental controls, chat workspace integration, and network topology.
- **MCP over HTTP and stdio**: Flexible transport configurations. HTTP transport supports multi-client and multi-instance pool connections.
- **Embedded API Documentation**: Every tool description automatically appends a reference link to the official Meilisearch markdown documentation (e.g. `.md`), allowing MCP-compatible LLMs to fetch and inspect the precise API specification on demand.
- **Resource Templates**: Exposes server state, configuration, and document previews directly via the `meilisearch://` URI scheme.
- **Context Prompts**: Predefined prompt templates to guide LLMs through searching, bulk data import, and settings setup.
- **Secure & Robust**: Rate limiting, CORS, reverse-proxy headers forwarding, and connection pool management built-in.

---

## Quick Start

### 1. Connect to Remote MCP Server
You can use the public MCP server at:
```
https://meilisearch.javad.dev/mcp
```
Example configuration for `mcp-remote`:
```json
{
  "command": "npx",
  "args": ["-y","mcp-remote@latest","https://meilisearch.javad.dev/mcp", "--header","X-Meili-Instance: ${MEILISEARCH_INSTANCE}", "--header","X-Meili-APIKey: ${MEILISEARCH_API_KEY}"],
  "env": {"MEILISEARCH_INSTANCE": "http://localhost:7700", "MEILISEARCH_API_KEY": "masterKey"},
  "active": true
}
```

### 2. Download Latest Release
Download the latest binary from the [releases page](https://github.com/ja7ad/meilisearch-mcp/releases) and make it executable:
```sh
wget https://github.com/ja7ad/meilisearch-mcp/releases/latest/download/mcp
chmod +x mcp
```

### 3. Build from Source

#### Prerequisites
- Go 1.26+
- Running [Meilisearch](https://www.meilisearch.com/) instance

#### Build
```sh
git clone https://github.com/ja7ad/meilisearch-mcp.git
cd meilisearch-mcp
make build
```

#### Run (HTTP Transport)
```sh
./build/meilisearch-mcp serve http --addr :8080 --meili-host http://localhost:7700 --meili-api-key masterKey
```

#### Run (STDIO Transport)
```sh
./build/meilisearch-mcp serve stdio --meili-host http://localhost:7700 --meili-api-key masterKey
```

---

## Available MCP Capabilities

### 1. Tools
A breakdown of all registered tools by domain. Each tool returns JSON or raw text.

#### 🗂️ Indexes
- `create_index`: Create a new index (supports `primary_key`).
- `update_index`: Update primary key or rename index UID.
- `get_index`: Retrieve index definition by name.
- `list_indexes`: Paginated list of all indexes.
- `delete_index`: Delete an existing index.
- `swap_index`: Atomically swap pairs of indexes.

#### 📝 Documents
- `get_document`: Retrieve a single document by its ID (supports selective `fields`).
- `get_documents`: Retrieve multiple documents (supports paginating, filtering, sorting, and IDs filtering).
- `add_documents`: Add or replace a JSON array of documents (asynchronous).
- `update_documents`: Partially update a JSON array of documents (asynchronous).
- `delete_document`: Delete a single document by ID.
- `delete_documents`: Bulk delete documents by a list of IDs.
- `delete_all_documents`: Truncate all documents in an index.
- `delete_documents_by_filter`: Delete documents matching a filter expression.

#### 🔍 Search
- `search`: Full search query interface (supports attributes, crops, highlights, filters, sort, scoring, and hybrid vector params).
- `facet_search`: Search for facet values inside an index.
- `search_similar_documents`: Find documents similar to a given document ID using embeddings.
- `multi_search`: Execute multiple search queries across different indexes in a single request.

#### ⚙️ Settings
- `get_settings`: Retrieve full index configuration.
- `update_settings`: Update setting fields (e.g. stopWords, synonyms, filterable/sortable attributes).
- `reset_settings`: Revert all settings back to default.

#### ⚡ Tasks & Keys
- `get_task`: Inspect details of a background operation.
- `list_tasks`: List all tasks (filterable by status, type, index, etc.).
- `cancel_tasks`: Cancel enqueued or processing tasks.
- `delete_tasks`: Clear finished task history.
- `list_keys`: List active API keys.
- `get_key`: Inspect details of a specific key.
- `create_key`: Create a new API key with actions & index scope.
- `update_key`: Update API key metadata.
- `delete_key`: Delete an API key.

#### 🔗 Webhooks & Network
- `add_webhook`: Create a new webhook trigger.
- `update_webhook`: Modify webhook URL and custom headers.
- `delete_webhook`: Delete a webhook.
- `get_webhook`: Retrieve webhook details.
- `list_webhooks`: List all configured webhooks.
- `get_network`: Read network topology (Experimental).
- `update_network`: Update network cluster topology (Experimental).

#### 💬 Chat Workspaces
- `chat_completion`: Chat with a Meilisearch chat workspace (collects stream chunks and returns aggregated response).
- `list_chat_workspaces`: List chat workspaces.
- `get_chat_workspace`: Get chat workspace definition.
- `get_chat_workspace_settings`: Get workspace keys/settings.
- `update_chat_workspace`: Update workspace settings.
- `reset_chat_workspace`: Revert workspace settings.

#### 🛠️ Experimental & System
- `get_experimental_features`: Get enabled experimental flags.
- `update_experimental_features`: Toggle experimental capabilities (e.g. compositeEmbedders, dynamicSearchRules).
- `get_stats`: Get stats for a specific index or global instance database.
- `get_health`: Get health check details.
- `get_version`: Get instance version.
- `create_dump`: Trigger a database dump (async).
- `create_snapshot`: Trigger a snapshot backup (async).

---

### 2. Resources
Clients can access database context using standard `meilisearch://` URIs:
- **`meilisearch://health`**: Current server health payload.
- **`meilisearch://version`**: Server version info.
- **`meilisearch://stats`**: Global instance statistics.
- **`meilisearch://indexes`**: List of all index definitions.
- **`meilisearch://indexes/{index_name}/stats`**: Performance metrics and statistics for a specific index.
- **`meilisearch://indexes/{index_name}/settings`**: Settings configuration for a specific index.
- **`meilisearch://indexes/{index_name}/documents`**: Preview of the first 20 documents.

---

### 3. Prompts
Templates to assist LLMs in performing common workflows:
- **`search_index`**: Prompts the LLM to search for information in an index (`index_name`, `query`).
- **`import_documents`**: Prompts the LLM on preparing and sending records via tools (`index_name`, `primary_key`).
- **`configure_settings`**: Walks through pulling settings, discussing filterable/sortable needs, and applying them (`index_name`).
- **`multi_search_help`**: Guides the LLM in structuring and performing searches across multiple indexes simultaneously (`queries`).

---

## Docker

### Build Image
```sh
docker build -t meilisearch-mcp:latest .
```

### Run Container
```sh
docker run --rm -p 8080:8080 \
  -e MEILI_HOST=http://host.docker.internal:7700 \
  -e MEILI_API_KEY=masterKey \
  meilisearch-mcp:latest
```

### Environment Variables
- `MEILI_HOST` (default: http://localhost:7700)
- `MEILI_API_KEY` (optional)
- `MCP_ADDR` (default: :8080)
- `MCP_POOL_SIZE` (default: 100)
- `MCP_POOL_DURATION` (default: 5m)
- `MCP_RATE_LIMIT_RPS` (default: 300)
- `MCP_DEBUG` (1 to enable debug logging)

### docker-compose Example
```yaml
services:
  mcp:
    image: meilisearch-mcp:latest
    build: .
    environment:
      MEILI_HOST: http://meili:7700
      MEILI_API_KEY: ${MEILI_API_KEY:-masterKey}
      MCP_RATE_LIMIT_RPS: 200
    ports:
      - "8080:8080"
    depends_on:
      - meili
  meili:
    image: getmeili/meilisearch:latest
    environment:
      MEILI_MASTER_KEY: masterKey
    ports:
      - "7700:7700"
```

---

## Configuration Examples

### Generic HTTP (mcp-remote)
```json
{
  "command": "npx",
  "args": ["-y","mcp-remote@latest","https://meilisearch.javad.dev/mcp", "--header","X-Meili-Instance: ${MEILISEARCH_INSTANCE}", "--header","X-Meili-APIKey: ${MEILISEARCH_API_KEY}"],
  "env": {"MEILISEARCH_INSTANCE": "http://localhost:7700", "MEILISEARCH_API_KEY": "masterKey"},
  "active": true
}
```

### Local STDIO
```json
{
  "command": "/usr/bin/meilisearch-mcp",
  "args": ["serve", "stdio","--meili-host","http://localhost:7700","--meili-api-key","masterKey"],
  "env": {},
  "active": false
}
```

---

## Supported Clients
- **Claude Desktop** — Full tool, resource, and prompt integrations.
- **Cursor / VS Code (Continue)** — Configure as standard stdio or remote HTTP MCP servers.
- **Zed** — Native editor tool integration.
- **Jan** — Native desktop assistant MCP config.

## Security Notes
- Place behind TLS (reverse proxy) when exposing the HTTP server publicly.
- Forward required auth headers only (`X-Meili-Instance`, `X-Meili-APIKey`) and strip unknown headers.
- Enforce appropriate payload size limits & rate limiting at the entry proxy.

## Contributing
Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for conventions.
