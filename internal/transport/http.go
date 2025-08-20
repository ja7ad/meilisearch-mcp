package transport

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/ja7ad/meilisearch-mcp/version"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed index.tpl
var indexTpl string

var indexTemplate = template.Must(template.New("index").Parse(indexTpl))

type HTTP struct {
	addr      string
	mcp       http.Handler
	sse       http.Handler
	httpSrv   *http.Server
	errCh     chan error
	enableSSE bool
	logging   *logger.SubLogger
}

func NewHTTP(mc *server.MCPServer, enableSSE bool, addr string) Server {
	s := &HTTP{
		addr:      addr,
		errCh:     make(chan error, 1),
		enableSSE: enableSSE,
	}

	s.logging = logger.NewSubLogger("_transport", s)

	httpHandler := server.NewStreamableHTTPServer(mc)

	var sseHandler http.Handler
	if enableSSE {
		sseHandler = server.NewSSEServer(mc)
		s.sse = sseHandler

		s.logging.Info("sse server enabled", "addr", addr+"/sse")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if enableSSE && r.URL.Path == "/sse" {
				sseHandler.ServeHTTP(w, r)
				return
			}
			httpHandler.ServeHTTP(w, r)
			return
		}

		// Exact "/"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Method == http.MethodHead {
			return
		}
		data := struct{ Version string }{Version: version.Version.String()}

		var buf bytes.Buffer
		if err := indexTemplate.Execute(&buf, data); err != nil {
			http.Error(w, "template render error", http.StatusInternalServerError)
			return
		}
		_, _ = buf.WriteTo(w)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		// Keep WriteTimeout generous when using SSE; the SSE handler
		// will keep the connection open. This timeout applies per write.
		WriteTimeout: 0, // 0 = no timeout; avoid closing long-lived SSE streams
		IdleTimeout:  120 * time.Second,
	}

	s.httpSrv = httpSrv
	s.mcp = httpHandler

	return s
}

func (*HTTP) String() string { return "HTTP" }

func (h *HTTP) Start(_ context.Context) {
	go func() {
		h.logging.Info("Starting http server", "addr", h.addr, "sse", h.enableSSE)
		if err := h.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.errCh <- err
			return
		}

		close(h.errCh)
	}()
}

func (h *HTTP) Err() <-chan error { return h.errCh }

func (h *HTTP) Stop(ctx context.Context) error {
	h.logging.Info("Shutting down http server", "addr", h.addr)
	return h.httpSrv.Shutdown(ctx)
}
