package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (p *Protocol) GetExperimentalFeatures() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_experimental_features",
			mcp.WithDescription("Get the state of all experimental features in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/experimental-features.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.ExperimentalFeatures().GetWithContext(ctx)
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

func (p *Protocol) UpdateExperimentalFeatures() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_experimental_features",
			mcp.WithDescription("Enable or disable specific experimental features in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/experimental-features.md"),
			mcp.WithBoolean("logsRoute", mcp.Description("Enable /logs route (optional)")),
			mcp.WithBoolean("metrics", mcp.Description("Enable /metrics route (optional)")),
			mcp.WithBoolean("editDocumentsByFunction", mcp.Description("Enable document editing by function (optional)")),
			mcp.WithBoolean("containsFilter", mcp.Description("Enable contains filter (optional)")),
			mcp.WithBoolean("network", mcp.Description("Enable network (optional)")),
			mcp.WithBoolean("compositeEmbedders", mcp.Description("Enable composite embedders (optional)")),
			mcp.WithBoolean("chatCompletions", mcp.Description("Enable chat completions (optional)")),
			mcp.WithBoolean("multimodal", mcp.Description("Enable multimodal support (optional)")),
			mcp.WithBoolean("dynamicSearchRules", mcp.Description("Enable dynamic search rules (optional)")),
			mcp.WithBoolean("getTaskDocumentsRoute", mcp.Description("Enable get task documents route (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			ef := client.ExperimentalFeatures()

			if val, ok := req.GetArguments()["logsRoute"].(bool); ok {
				ef.SetLogsRoute(val)
			}
			if val, ok := req.GetArguments()["metrics"].(bool); ok {
				ef.SetMetrics(val)
			}
			if val, ok := req.GetArguments()["editDocumentsByFunction"].(bool); ok {
				ef.SetEditDocumentsByFunction(val)
			}
			if val, ok := req.GetArguments()["containsFilter"].(bool); ok {
				ef.SetContainsFilter(val)
			}
			if val, ok := req.GetArguments()["network"].(bool); ok {
				ef.SetNetwork(val)
			}
			if val, ok := req.GetArguments()["compositeEmbedders"].(bool); ok {
				ef.SetCompositeEmbedders(val)
			}
			if val, ok := req.GetArguments()["chatCompletions"].(bool); ok {
				ef.SetChatCompletions(val)
			}
			if val, ok := req.GetArguments()["multimodal"].(bool); ok {
				ef.SetMultiModal(val)
			}
			if val, ok := req.GetArguments()["dynamicSearchRules"].(bool); ok {
				ef.SetDynamicSearchRules(val)
			}
			if val, ok := req.GetArguments()["getTaskDocumentsRoute"].(bool); ok {
				ef.SetGetTaskDocumentsRoute(val)
			}

			resp, err := ef.UpdateWithContext(ctx)
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
