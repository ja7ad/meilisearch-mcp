package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetStats() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_stats",
			mcp.WithDescription("Get statistical information about an index or the entire database. Reference: https://www.meilisearch.com/docs/reference/api/stats.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index (optional, gets global stats if omitted)")),
			mcp.WithBoolean("showInternalDatabaseSizes", mcp.Description("Include internal database sizes (optional)")),
			mcp.WithString("sizeFormat", mcp.Description("Format for database size: 'human' or 'raw' (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			indexName, err := OptionalParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var statsParams meilisearch.StatsParams
			if err := sonic.Unmarshal(argsBytes, &statsParams); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if indexName != "" {
				idx := client.Index(indexName)
				resp, err := idx.GetStatsWithContext(ctx, &statsParams)
				if err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
				b, err := sonic.Marshal(resp)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(b)), nil
			}

			resp, err := client.GetStatsWithContext(ctx, &statsParams)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) GetHealth() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_health",
			mcp.WithDescription("Get the health status of the Meilisearch instance. Reference: https://www.meilisearch.com/docs/reference/api/health.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.HealthWithContext(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) GetVersion() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_version",
			mcp.WithDescription("Get the version of the Meilisearch instance. Reference: https://www.meilisearch.com/docs/reference/api/version.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.VersionWithContext(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) CreateDump() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_dump",
			mcp.WithDescription("Create a database dump (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/backups/create-dump.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.CreateDumpWithContext(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) CreateSnapshot() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_snapshot",
			mcp.WithDescription("Create a database snapshot (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/backups/create-snapshot.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.CreateSnapshotWithContext(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(b)), nil
		}
}
