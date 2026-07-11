package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) ChatCompletion() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"chat_completion",
			mcp.WithDescription("Get a chat completion from a Meilisearch chat workspace. Reference: https://www.meilisearch.com/docs/reference/api.md"),
			mcp.WithString("workspace", mcp.Description("The ID/UID of the chat workspace to use"), mcp.Required()),
			mcp.WithString("model", mcp.Description("Model to run the query on"), mcp.Required()),
			mcp.WithArray("messages", mcp.Description("Chat history messages in standard OpenAI format"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspace, err := RequiredParam[string](req, "workspace")
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

			var chatQuery meilisearch.ChatCompletionQuery
			if err := sonic.Unmarshal(argsBytes, &chatQuery); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			chatQuery.Stream = true

			stream, err := client.ChatCompletionStreamWithContext(ctx, workspace, &chatQuery)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			defer func() { _ = stream.Close() }()

			var fullText string
			for stream.Next() {
				chunk := stream.Current()
				if chunk != nil && len(chunk.Choices) > 0 {
					choice := chunk.Choices[0]
					if choice.Delta != nil && choice.Delta.Content != nil {
						fullText += *choice.Delta.Content
					}
				}
			}
			if err := stream.Err(); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(fullText), nil
		}
}

func (p *Protocol) ListChatWorkspaces() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_chat_workspaces",
			mcp.WithDescription("List all chat workspaces configured in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api.md"),
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
				limit = 20
			}

			resp, err := client.ListChatWorkspacesWithContext(ctx, &meilisearch.ListChatWorkSpaceQuery{
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

func (p *Protocol) GetChatWorkspace() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_chat_workspace",
			mcp.WithDescription("Get a chat workspace configuration by its ID. Reference: https://www.meilisearch.com/docs/reference/api.md"),
			mcp.WithString("workspace", mcp.Description("The ID/UID of the workspace to retrieve"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspace, err := RequiredParam[string](req, "workspace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.GetChatWorkspaceWithContext(ctx, workspace)
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

func (p *Protocol) GetChatWorkspaceSettings() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_chat_workspace_settings",
			mcp.WithDescription("Get chat workspace settings by workspace ID. Reference: https://www.meilisearch.com/docs/reference/api.md"),
			mcp.WithString("workspace", mcp.Description("The ID/UID of the workspace"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspace, err := RequiredParam[string](req, "workspace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.GetChatWorkspaceSettingsWithContext(ctx, workspace)
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

func (p *Protocol) UpdateChatWorkspace() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_chat_workspace",
			mcp.WithDescription("Update a chat workspace by its ID. Reference: https://www.meilisearch.com/docs/reference/api.md"),
			mcp.WithString("workspace", mcp.Description("The ID/UID of the workspace to update"), mcp.Required()),
			mcp.WithObject("settings", mcp.Description("The updated settings object for the chat workspace"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspace, err := RequiredParam[string](req, "workspace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			settingsVal, ok := req.GetArguments()["settings"]
			if !ok {
				return mcp.NewToolResultError("missing required parameter: settings"), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			settingsBytes, err := sonic.Marshal(settingsVal)
			if err != nil {
				return nil, err
			}

			var settings meilisearch.ChatWorkspaceSettings
			if err := sonic.Unmarshal(settingsBytes, &settings); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.UpdateChatWorkspaceWithContext(ctx, workspace, &settings)
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

func (p *Protocol) ResetChatWorkspace() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"reset_chat_workspace",
			mcp.WithDescription("Reset a chat workspace configuration by its ID. Reference: https://www.meilisearch.com/docs/reference/api.md"),
			mcp.WithString("workspace", mcp.Description("The ID/UID of the workspace to reset"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			workspace, err := RequiredParam[string](req, "workspace")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.ResetChatWorkspaceWithContext(ctx, workspace)
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
