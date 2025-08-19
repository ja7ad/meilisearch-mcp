package transport

import (
	"context"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/mark3labs/mcp-go/server"
)

type HTTP struct {
	addr    string
	srv     *server.StreamableHTTPServer
	errCh   chan error
	logging *logger.SubLogger
}

func NewHTTP(mc *server.MCPServer, addr string) Server {
	srv := &HTTP{
		srv:   server.NewStreamableHTTPServer(mc),
		errCh: make(chan error),
		addr:  addr,
	}

	srv.logging = logger.NewSubLogger("_transport", srv)

	return srv
}

func (*HTTP) String() string {
	return "HTTP"
}

func (h *HTTP) Start(_ context.Context) {
	go func() {
		h.logging.Info("Starting http server", "addr", h.addr)
		h.errCh <- h.srv.Start(h.addr)
	}()
}

func (h *HTTP) Err() <-chan error {
	return h.errCh
}

func (h *HTTP) Stop(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
