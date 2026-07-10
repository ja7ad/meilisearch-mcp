package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetNetwork() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_network",
			mcp.WithDescription("Get the current value of the instance’s network topology configuration (Experimental). Reference: https://www.meilisearch.com/docs/reference/api/experimental-features/get-network-topology.md"),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			net, err := client.GetNetworkWithContext(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(net)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) UpdateNetwork() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_network",
			mcp.WithDescription("Update the instance’s network topology configuration (Experimental). Reference: https://www.meilisearch.com/docs/reference/api/experimental-features/configure-network-topology.md"),
			mcp.WithObject("network", mcp.Description("The network topology configuration object to update"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			netVal, ok := req.GetArguments()["network"]
			if !ok {
				return mcp.NewToolResultError("missing required parameter: network"), nil
			}

			argsBytes, err := sonic.Marshal(netVal)
			if err != nil {
				return nil, err
			}

			var netReq meilisearch.UpdateNetworkRequest
			if err := sonic.Unmarshal(argsBytes, &netReq); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.UpdateNetworkWithContext(ctx, &netReq)
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
