package transport

import (
	"context"
	"os"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/mark3labs/mcp-go/server"
)

type Stdio struct {
	srv     *server.StdioServer
	errCh   chan error
	logging *logger.SubLogger
}

func NewStdio(mc *server.MCPServer) Server {
	srv := &Stdio{
		srv:   server.NewStdioServer(mc),
		errCh: make(chan error),
	}

	srv.logging = logger.NewSubLogger("_transport", srv)

	return srv
}

func (*Stdio) String() string {
	return "Stdio"
}

func (s *Stdio) Start(ctx context.Context) {
	go func() {
		s.logging.Info("Starting stdio server")
		s.errCh <- s.srv.Listen(ctx, os.Stdin, os.Stdout)
	}()
}

func (s *Stdio) Stop(_ context.Context) error {
	return nil
}

func (s *Stdio) Err() <-chan error {
	return s.errCh
}
