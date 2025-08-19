package cli

import (
	"errors"
	"os/signal"
	"syscall"
	"time"

	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/ja7ad/meilisearch-mcp/internal/pool"
	"github.com/ja7ad/meilisearch-mcp/internal/protocol"
	"github.com/ja7ad/meilisearch-mcp/internal/transport"
	"github.com/ja7ad/meilisearch-mcp/version"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cobra"
)

func (c *CLI) serve(debug bool) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long:  "Start the MCP server to handle requests for Meilisearch.",
	}

	host := cmd.Flags().String("meili-host", "http://localhost:7700",
		"Meilisearch host (e.g. http://127.0.0.1:7700)")
	apiKey := cmd.Flags().String("meili-api-key", "", "Meilisearch API key (optional)")
	tran := cmd.Flags().String("transport", "stdio", "Transport protocol to use (e.g. stdio, http)")
	poolSize := cmd.Flags().Int("pool-size", 100, "Size of the connection pool for HTTP transport")
	poolDuration := cmd.Flags().Duration("pool-duration", 5*time.Minute,
		"Duration for which connections are kept in the pool (e.g. 30s, 5m)")
	addr := cmd.Flags().String("addr", ":8080",
		"Address to bind the MCP server (only for HTTP transport)")

	c.root.AddCommand(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg := logger.DefaultConfig()
		if debug {
			cfg.Levels["default"] = "debug"
			cfg.Levels["_transport"] = "debug"
			cfg.Levels["_protocol"] = "debug"
			cfg.Levels["_pool"] = "debug"
		}

		logger.InitGlobalLogger(cfg)
		logger.Info("Starting Meilisearch MCP server", "version", version.Version.String())

		t := protocol.Transport(*tran)
		var p *pool.Pool[string, meilisearch.ServiceManager]
		var srv transport.Server

		mc := server.NewMCPServer("Meilisearch MCP",
			version.Version.String(),
			server.WithToolCapabilities(true),
			server.WithResourceCapabilities(true, true),
		)

		switch t {
		case protocol.TransportStdio:
			if *host == "" {
				return errors.New("missing --meili-host for stdio transport")
			}
			srv = transport.NewStdio(mc)
		case protocol.TransportHTTP:
			if *poolSize <= 0 {
				return errors.New("pool-size must be > 0")
			}
			if *poolDuration <= 0 {
				return errors.New("pool-duration must be > 0")
			}
			if *addr == "" {
				return errors.New("missing --addr for HTTP transport")
			}

			p = pool.New(*poolSize, *poolDuration)
			srv = transport.NewHTTP(mc, *addr)
		default:
			return errors.New("invalid --transport (expected stdio|http)")
		}

		proto := protocol.New(t, *host, *apiKey, p)
		rt := transport.NewRoute(mc, proto)
		rt.Register()

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go srv.Start(ctx)

		<-ctx.Done()
		return nil
	}
}
