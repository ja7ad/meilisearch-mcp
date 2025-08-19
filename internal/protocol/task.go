package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (p *Protocol) GetTask() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_task",
			mcp.WithDescription("Get a task by its ID"),
			mcp.WithNumber(
				"task_id",
				mcp.Description("ID of the task to retrieve"),
				mcp.Required(),
				mcp.Min(0),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskID, err := Required64Int(req, "task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			task, err := client.GetTaskWithContext(ctx, taskID)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(task)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}
