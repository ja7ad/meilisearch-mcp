FROM golang:1.25 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/meilisearch-mcp ./cmd/mcp

FROM alpine:3.20
RUN adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=build /out/meilisearch-mcp /app/meilisearch-mcp
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh && chown -R app:app /app

USER app

ENV MEILI_HOST="http://localhost:7700" \
    MEILI_API_KEY="" \
    MCP_ADDR=":8080" \
    MCP_POOL_SIZE=100 \
    MCP_POOL_DURATION="5m" \
    MCP_RATE_LIMIT_RPS=300 \
    MCP_DEBUG=0

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -q --spider http://127.0.0.1:8080/ || exit 1

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD []
