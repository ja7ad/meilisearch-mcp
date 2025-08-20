package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) ListKeys() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_keys",
			mcp.WithDescription("List all API keys in Meilisearch"),
			WithPagination(),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			limit, err := OptionalInt64Param(req, "limit")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			offset, err := OptionalInt64Param(req, "offset")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if limit <= 0 {
				limit = 20 // Default limit if not specified
			}

			if offset <= 0 {
				offset = 0 // Default offset if not specified
			}

			resp, err := client.GetKeysWithContext(ctx, &meilisearch.KeysQuery{
				Limit:  limit,
				Offset: offset,
			})
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

func (p *Protocol) GetKey() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_key",
			mcp.WithDescription("Get API key"),
			mcp.WithString("key", mcp.Description("API key")),
			mcp.WithString("uid", mcp.Description("Key UID")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			key, err := OptionalParam[string](req, "key")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			uid, err := OptionalParam[string](req, "uid")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if key == "" && uid == "" {
				return mcp.NewToolResultError("API key or UID is required"), nil
			} else if key != "" {
				if err := p.validate(key, "max=250,min=1"); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
			} else if uid != "" {
				if err := p.validate(uid, "max=250,min=1"); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.GetKeyWithContext(ctx, key)
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
