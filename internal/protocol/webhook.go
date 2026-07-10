package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) AddWebhook() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"add_webhook",
			mcp.WithDescription("Add a new webhook to Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/webhooks/create-webhook.md"),
			mcp.WithString("url", mcp.Description("The target URL of the webhook"), mcp.Required()),
			mcp.WithObject("headers", mcp.Description("Optional HTTP headers to send with webhook events")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			url, err := RequiredParam[string](req, "url")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			headersVal, err := OptionalParam[map[string]any](req, "headers")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			headers := make(map[string]string)
			for k, v := range headersVal {
				if s, ok := v.(string); ok {
					headers[k] = s
				}
			}

			resp, err := client.AddWebhookWithContext(ctx, &meilisearch.AddWebhookRequest{
				URL:     url,
				Headers: headers,
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

func (p *Protocol) UpdateWebhook() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_webhook",
			mcp.WithDescription("Modify a previously existing webhook in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/webhooks/update-webhook.md"),
			mcp.WithString("uuid", mcp.Description("The UUID of the webhook to update"), mcp.Required()),
			mcp.WithString("url", mcp.Description("The new target URL of the webhook")),
			mcp.WithObject("headers", mcp.Description("Optional updated HTTP headers to send")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uuid, err := RequiredParam[string](req, "uuid")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			url, err := OptionalParam[string](req, "url")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			headersVal, err := OptionalParam[map[string]any](req, "headers")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var headers map[string]string
			if len(headersVal) > 0 {
				headers = make(map[string]string)
				for k, v := range headersVal {
					if s, ok := v.(string); ok {
						headers[k] = s
					}
				}
			}

			resp, err := client.UpdateWebhookWithContext(ctx, uuid, &meilisearch.UpdateWebhookRequest{
				URL:     url,
				Headers: headers,
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

func (p *Protocol) DeleteWebhook() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_webhook",
			mcp.WithDescription("Delete an existing webhook in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/webhooks/delete-webhook.md"),
			mcp.WithString("uuid", mcp.Description("The UUID of the webhook to delete"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uuid, err := RequiredParam[string](req, "uuid")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			err = client.DeleteWebhookWithContext(ctx, uuid)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(map[string]bool{"success": true})
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) GetWebhook() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_webhook",
			mcp.WithDescription("Get a webhook by its UUID in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/webhooks/get-webhook.md"),
			mcp.WithString("uuid", mcp.Description("The UUID of the webhook to retrieve"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uuid, err := RequiredParam[string](req, "uuid")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.GetWebhookWithContext(ctx, uuid)
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

func (p *Protocol) ListWebhooks() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_webhooks",
			mcp.WithDescription("List all webhooks configured in Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/webhooks/list-webhooks.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.ListWebhooksWithContext(ctx)
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
