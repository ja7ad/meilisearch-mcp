package protocol

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (p *Protocol) GetTask() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_task",
			mcp.WithDescription("Get a task by its ID"),
			mcp.WithString(
				"task_id",
				mcp.Description("ID of the task to retrieve"),
				mcp.Required(),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			taskIDRaw, err := RequiredParam[string](req, "task_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			taskUID, err := strconv.ParseInt(taskIDRaw, 10, 64)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			task, err := client.GetTaskWithContext(ctx, taskUID)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultStructured(
				task,
				fmt.Sprintf("task with ID %d retrieved, status: %s", taskUID, task.Status),
			), nil
		}
}
