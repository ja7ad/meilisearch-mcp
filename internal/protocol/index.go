package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) CreateIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"create_index",
			mcp.WithDescription("Create a new index in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/indexes/create-index.md"),
			mcp.WithString("index_name",
				mcp.Description("Name of the index to create"),
				mcp.Required(),
			),
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

			if err := p.validate(uid, "max=250,min=1"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := p.validate(primaryKey, "omitempty,max=250,min=1"); err != nil {
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

			b, err := sonic.Marshal(task)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) UpdateIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_index",
			mcp.WithDescription("Update an primary key or UID of index in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/indexes/update-index.md"),
			mcp.WithString("index_name",
				mcp.Description("Name of the index to update"),
				mcp.Required(),
			),
			mcp.WithString("primary_key", mcp.Description("Primary key for the index (optional)")),
			mcp.WithString("new_uid", mcp.Description("New UID for the index (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			primaryKey, err := OptionalParam[string](req, "primary_key")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			newUID, err := OptionalParam[string](req, "new_uid")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := p.validate(uid, "max=250,min=1"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if primaryKey != "" {
				if err := p.validate(primaryKey, "max=250,min=1"); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
			}

			if newUID != "" {
				if err := p.validate(newUID, "max=250,min=1"); err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.UpdateIndex(&meilisearch.UpdateIndexRequestParams{
				PrimaryKey: primaryKey,
				UID:        newUID,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(task)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) DeleteIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_index",
			mcp.WithDescription("Delete an index in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/indexes/delete-index.md"),
			mcp.WithString("index_name",
				mcp.Description("Name of the index for delete"),
				mcp.Required(),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			task, err := client.DeleteIndex(uid)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(task)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) GetIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_index",
			mcp.WithDescription("Get an index by its name. Reference: https://www.meilisearch.com/docs/reference/api/indexes/get-index.md"),
			mcp.WithString(
				"index_name",
				mcp.Description("Name of the index to retrieve"),
				mcp.Required(),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if err := p.validate(uid, "max=250,min=1"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			index, err := client.GetIndexWithContext(ctx, uid)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(index)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) ListIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_indexes",
			mcp.WithDescription("List all indexes in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/indexes/list-indexes.md"),
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

			res, err := client.ListIndexesWithContext(ctx, &meilisearch.IndexesQuery{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(res)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) SwapIndex() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("swap_index",
			mcp.WithDescription("Swap two indexes in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/indexes/swap-indexes.md"),
			mcp.WithArray("indexes",
				mcp.Required(),
				mcp.Description("Indexes to swap"),
				mcp.Items(map[string]any{
					"type":     "array",
					"minItems": 2,
					"maxItems": 2,
					"items": map[string]any{
						"type": "string",
					},
				}),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			idxObj, ok := req.GetArguments()["indexes"].([]interface{})
			if !ok {
				return mcp.NewToolResultError("indexes parameter must be an array of string two pairs, " +
					"for example: { \"indexes\": [ [ \"foobar1\", \"foobar2\" ], [ \"foobar3\", \"foobar4\" ] ] }"), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			task, err := client.SwapIndexesWithContext(ctx, swapIndexes(idxObj))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(task)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}
