package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetTask() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_task",
			mcp.WithDescription("Get a task by its ID. Reference: https://www.meilisearch.com/docs/reference/api/async-task-management/get-task.md"),
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
				return mcp.NewToolResultError(err.Error()), nil
			}

			task, err := client.GetTaskWithContext(ctx, taskID)
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

func (p *Protocol) ListTasks() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"list_tasks",
			mcp.WithDescription("List Meilisearch tasks with filters and pagination. Reference: https://www.meilisearch.com/docs/reference/api/async-task-management/list-tasks.md"),
			mcp.WithNumber("limit", mcp.Description("Limit the number of tasks returned (default 20)")),
			mcp.WithNumber("from", mcp.Description("Retrieve tasks starting from a specific task ID")),
			mcp.WithArray("index_uids", mcp.Description("Filter tasks by index UIDs"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("statuses", mcp.Description("Filter tasks by status (e.g. enqueued, processing, succeeded, failed)"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("types", mcp.Description("Filter tasks by type (e.g. indexCreation, indexUpdate, documentAdditionOrUpdate)"), mcp.Items(map[string]any{"type": "string"})),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var query meilisearch.TasksQuery
			if err := sonic.Unmarshal(argsBytes, &query); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if query.Limit <= 0 {
				query.Limit = 20
			}

			res, err := client.GetTasksWithContext(ctx, &query)
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

func (p *Protocol) CancelTasks() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"cancel_tasks",
			mcp.WithDescription("Cancel one or more tasks based on filters. Reference: https://www.meilisearch.com/docs/reference/api/async-task-management/cancel-tasks.md"),
			mcp.WithArray("statuses", mcp.Description("Cancel tasks matching statuses (e.g. enqueued, processing)"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("types", mcp.Description("Cancel tasks matching types"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("index_uids", mcp.Description("Cancel tasks matching index UIDs"), mcp.Items(map[string]any{"type": "string"})),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var query meilisearch.CancelTasksQuery
			if err := sonic.Unmarshal(argsBytes, &query); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			res, err := client.CancelTasksWithContext(ctx, &query)
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

func (p *Protocol) DeleteTasks() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_tasks",
			mcp.WithDescription("Delete completed tasks based on filters. Reference: https://www.meilisearch.com/docs/reference/api/async-task-management/delete-tasks.md"),
			mcp.WithArray("statuses", mcp.Description("Delete tasks matching statuses (e.g. succeeded, failed, canceled)"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("types", mcp.Description("Delete tasks matching types"), mcp.Items(map[string]any{"type": "string"})),
			mcp.WithArray("index_uids", mcp.Description("Delete tasks matching index UIDs"), mcp.Items(map[string]any{"type": "string"})),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var query meilisearch.DeleteTasksQuery
			if err := sonic.Unmarshal(argsBytes, &query); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			res, err := client.DeleteTasksWithContext(ctx, &query)
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
