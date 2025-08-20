package transport

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/ja7ad/meilisearch-mcp/internal/util"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/time/rate"
)

type ToolMiddleware func(server.ToolHandlerFunc) server.ToolHandlerFunc

func ChainToolMiddleware(mw ...ToolMiddleware) ToolMiddleware {
	return func(final server.ToolHandlerFunc) server.ToolHandlerFunc {
		for i := len(mw) - 1; i >= 0; i-- {
			final = mw[i](final)
		}
		return final
	}
}

type RateLimitMiddleware struct {
	limiters map[string]*rate.Limiter
	mutex    sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimitMiddleware(requestsPerSecond float64, burst int) *RateLimitMiddleware {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 100 // Default rate limit if not specified
	}

	return &RateLimitMiddleware{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

func (m *RateLimitMiddleware) getLimiter(sessionID string) *rate.Limiter {
	m.mutex.RLock()
	limiter, exists := m.limiters[sessionID]
	m.mutex.RUnlock()

	if !exists {
		m.mutex.Lock()
		limiter = rate.NewLimiter(m.rate, m.burst)
		m.limiters[sessionID] = limiter
		m.mutex.Unlock()
	}

	return limiter
}

func (m *RateLimitMiddleware) ToolMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host, _, hash := util.MeilisearchHeaders(req.Header)
		if hash == "" {
			return nil, fmt.Errorf("missing x-meili-instance in request headers")
		}

		limiter := m.getLimiter(hash)

		if !limiter.Allow() {
			return nil, fmt.Errorf("rate limit exceeded for host %s", host)
		}

		return next(ctx, req)
	}
}

func LoggerToolMiddleware(logging *logger.SubLogger) func(server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			fields := []interface{}{
				"method", req.Method,
				"name", req.Params.Name,
				"duration_ms", duration.Milliseconds(),
			}

			if err != nil {
				logging.Error("Error tool request", "error", err, fields)
				return nil, err
			}

			logging.Info("Success tool request", fields...)

			return resp, nil
		}
	}
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
