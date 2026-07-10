package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetSettings() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_settings",
			mcp.WithDescription("Get the settings of a Meilisearch index. Reference: https://www.meilisearch.com/docs/reference/api/settings.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			resp, err := idx.GetSettingsWithContext(ctx)
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

func (p *Protocol) UpdateSettings() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_settings",
			mcp.WithDescription("Update settings of a Meilisearch index (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/settings.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithObject("settings", mcp.Description("The settings object to apply (e.g. filterableAttributes, sortableAttributes)"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
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

			// Unmarshal settingsVal into meilisearch.Settings
			bSettings, err := sonic.Marshal(settingsVal)
			if err != nil {
				return nil, err
			}

			var settings meilisearch.Settings
			if err := sonic.Unmarshal(bSettings, &settings); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.UpdateSettingsWithContext(ctx, &settings)
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

func (p *Protocol) ResetSettings() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"reset_settings",
			mcp.WithDescription("Reset all settings of a Meilisearch index to their default values (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/settings.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.ResetSettingsWithContext(ctx)
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
