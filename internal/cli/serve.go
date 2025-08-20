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
	"github.com/spf13/cobra"
)

func (c *CLI) serve(debug bool) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server",
		Long:  "Start the MCP server to handle requests for Meilisearch.",
	}

	c.root.AddCommand(cmd)

	c.http(cmd, debug)
	c.stdio(cmd, debug)
}

func (c *CLI) http(serve *cobra.Command, debug bool) {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Start the HTTP MCP server",
		Long:  "Start the HTTP MCP server to handle requests for Meilisearch.",
	}

	host := cmd.Flags().String("meili-host", "http://localhost:7700",
		"Meilisearch host (e.g. http://127.0.0.1:7700)")
	apiKey := cmd.Flags().String("meili-api-key", "", "Meilisearch API key (optional)")
	poolSize := cmd.Flags().Int("pool-size", 100, "Size of the connection pool for HTTP transport")
	poolDuration := cmd.Flags().Duration("pool-duration", 5*time.Minute,
		"Duration for which connections are kept in the pool (e.g. 30s, 5m)")
	addr := cmd.Flags().String("addr", ":8080",
		"Address to bind the MCP server (only for HTTP transport)")
	rateLimitReqPerSec := cmd.Flags().Float64("rate-limit-req-per-sec", 300, "Rate limit requests per second (Default: 300)")

	serve.AddCommand(cmd)

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

		mc := server.NewMCPServer("Meilisearch HTTP MCP",
			version.Version.String(),
			server.WithToolCapabilities(true),
			server.WithResourceCapabilities(true, true),
		)

		if *poolSize <= 0 {
			return errors.New("pool-size must be > 0")
		}
		if *poolDuration <= 0 {
			return errors.New("pool-duration must be > 0")
		}
		if *addr == "" {
			return errors.New("missing --addr for HTTP transport")
		}

		p := pool.New(*poolSize, *poolDuration)
		srv := transport.NewHTTP(mc, *addr)

		proto := protocol.New(protocol.TransportHTTP, *host, *apiKey, p)
		rt := transport.NewRoute(mc, proto,
			transport.NewRateLimitMiddleware(*rateLimitReqPerSec, 10).ToolMiddleware,
			transport.LoggerToolMiddleware(logger.NewSubLogger("_transport", nil)),
		)
		rt.Register()

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go srv.Start(ctx)

		<-ctx.Done()
		return nil
	}
}

func (c *CLI) stdio(serve *cobra.Command, debug bool) {
	cmd := &cobra.Command{
		Use:   "stdio",
		Short: "Start the Stdio MCP server",
		Long:  "Start the Stdio MCP server to handle requests for Meilisearch.",
	}

	host := cmd.Flags().String("meili-host", "http://localhost:7700",
		"Meilisearch host (e.g. http://127.0.0.1:7700)")
	apiKey := cmd.Flags().String("meili-api-key", "", "Meilisearch API key (optional)")

	serve.AddCommand(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg := logger.DefaultConfig()
		if debug {
			cfg.Levels["default"] = "debug"
			cfg.Levels["_transport"] = "debug"
			cfg.Levels["_protocol"] = "debug"
			cfg.Levels["_pool"] = "debug"
		}

		logger.InitGlobalLogger(cfg)
		logger.Info("Starting Meilisearch Stdio MCP server", "version", version.Version.String())

		mc := server.NewMCPServer("Meilisearch Stdio MCP",
			version.Version.String(),
			server.WithToolCapabilities(true),
			server.WithResourceCapabilities(true, true),
		)

		if *host == "" {
			return errors.New("missing --meili-host for stdio transport")
		}
		srv := transport.NewStdio(mc)

		proto := protocol.New(protocol.TransportStdio, *host, *apiKey, nil)
		rt := transport.NewRoute(mc, proto,
			transport.LoggerToolMiddleware(logger.NewSubLogger("_transport", nil)),
		)
		rt.Register()

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go srv.Start(ctx)

		<-ctx.Done()
		return nil
	}
}
