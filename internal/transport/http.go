package transport

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/ja7ad/meilisearch-mcp/version"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed index.tpl
var indexTpl string

var indexTemplate = template.Must(template.New("index").Parse(indexTpl))

type HTTP struct {
	addr    string
	mcp     http.Handler
	httpSrv *http.Server
	errCh   chan error
	logging *logger.SubLogger
}

func NewHTTP(mc *server.MCPServer, addr string) Server {
	mcpHandler := server.NewStreamableHTTPServer(mc)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
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
			return
		}
		mcpHandler.ServeHTTP(w, r)
	})

	httpSrv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	s := &HTTP{
		addr:    addr,
		mcp:     mcpHandler,
		httpSrv: httpSrv,
		errCh:   make(chan error, 1),
	}

	s.logging = logger.NewSubLogger("_transport", s)
	return s
}

func (*HTTP) String() string {
	return "HTTP"
}

func (h *HTTP) Start(_ context.Context) {
	go func() {
		h.logging.Info("Starting http server", "addr", h.addr)
		h.errCh <- h.httpSrv.ListenAndServe()
	}()
}

func (h *HTTP) Err() <-chan error {
	return h.errCh
}

func (h *HTTP) Stop(ctx context.Context) error {
	return h.httpSrv.Shutdown(ctx)
}
