#!/bin/sh
set -euo pipefail

# Allow optional debug
DEBUG_FLAG=""
if [ "${MCP_DEBUG:-0}" = "1" ]; then
  DEBUG_FLAG="--debug"
fi

# Defaults (mirroring ENV defaults set in Dockerfile)
MEILI_HOST="${MEILI_HOST:-http://localhost:7700}"
MEILI_API_KEY="${MEILI_API_KEY:-}"
MCP_ADDR="${MCP_ADDR:-:8080}"
MCP_POOL_SIZE="${MCP_POOL_SIZE:-100}"
MCP_POOL_DURATION="${MCP_POOL_DURATION:-5m}"
MCP_RATE_LIMIT_RPS="${MCP_RATE_LIMIT_RPS:-300}"
MCP_ENABLE_SSE="${MCP_ENABLE_SSE:-false}"

CMD_ARGS="serve http --addr ${MCP_ADDR} --meili-host ${MEILI_HOST} --pool-size ${MCP_POOL_SIZE} --pool-duration ${MCP_POOL_DURATION} --rate-limit-req-per-sec ${MCP_RATE_LIMIT_RPS}"

if [ -n "${MEILI_API_KEY}" ]; then
  CMD_ARGS="${CMD_ARGS} --meili-api-key ${MEILI_API_KEY}"
fi

# Normalize and add --sse flag if enabled (accept: true|1|yes case-insensitive)
case "$(printf '%s' "${MCP_ENABLE_SSE}" | tr '[:upper:]' '[:lower:]')" in
  true|1|yes) CMD_ARGS="${CMD_ARGS} --sse" ;;
  *) : ;; # ignore otherwise
esac

# If user passes explicit command (e.g. `bash`), run that instead
if [ "${1:-}" = "sh" ] || [ "${1:-}" = "bash" ] || [ "${1:-}" = "meilisearch-mcp" ]; then
  exec "$@"
fi

exec /app/meilisearch-mcp ${DEBUG_FLAG} ${CMD_ARGS} "$@"
