package protocol

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) CreateIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_index",
			mcp.WithDescription("Create a new index in Meilisearch"),
			mcp.WithString("index_name", mcp.Description("Name of the index to create")),
			mcp.WithString("primary_key", mcp.Description("Primary key for the index (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			primaryKey, err := OptionalParam[string](req, "primary_key")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := p.validate(uid, "max=100,min=1"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := p.validate(primaryKey, "omitempty,max=100,min=1"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			task, err := client.CreateIndex(&meilisearch.IndexConfig{
				Uid:        uid,
				PrimaryKey: primaryKey,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultStructured(
				task,
				fmt.Sprintf("index with name %s created, task id %d and status %v", uid,
					task.TaskUID,
					task.Status,
				)), nil
		}
}
