package protocol

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (p *Protocol) SearchIndexPrompt() (prompt mcp.Prompt, handler server.PromptHandlerFunc) {
	return mcp.NewPrompt(
			"search_index",
			mcp.WithPromptDescription("Prompt to search documents in a specific index"),
			mcp.WithArgument("index_name", mcp.ArgumentDescription("The index you want to search in"), mcp.RequiredArgument()),
			mcp.WithArgument("query", mcp.ArgumentDescription("The search query (optional)")),
		), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			indexName := req.Params.Arguments["index_name"]
			query := req.Params.Arguments["query"]

			var promptText string
			if query != "" {
				promptText = fmt.Sprintf("Please use the 'search' tool to search inside the Meilisearch index '%s' for the query '%s'. Analyze and summarize the results.", indexName, query)
			} else {
				promptText = fmt.Sprintf("Please use the 'search' tool to search inside the Meilisearch index '%s' with an empty query (or perform a general retrieval). List the returned documents and explain their structure.", indexName)
			}

			return mcp.NewGetPromptResult(
				"Search Index Helper",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		}
}

func (p *Protocol) ImportDocumentsPrompt() (prompt mcp.Prompt, handler server.PromptHandlerFunc) {
	return mcp.NewPrompt(
			"import_documents",
			mcp.WithPromptDescription("Prompt to help import documents into Meilisearch"),
			mcp.WithArgument("index_name", mcp.ArgumentDescription("The index you want to import into"), mcp.RequiredArgument()),
			mcp.WithArgument("primary_key", mcp.ArgumentDescription("Primary key of the index (optional)")),
		), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			indexName := req.Params.Arguments["index_name"]
			primaryKey := req.Params.Arguments["primary_key"]

			promptText := fmt.Sprintf("Please guide the user in preparing their documents for import. Once the documents are formatted as a JSON array of objects, use the 'add_documents' tool to upload them to the index '%s'.", indexName)
			if primaryKey != "" {
				promptText += fmt.Sprintf(" Note that the primary key for this index is set to '%s'.", primaryKey)
			}

			return mcp.NewGetPromptResult(
				"Import Documents Helper",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		}
}

func (p *Protocol) ConfigureSettingsPrompt() (prompt mcp.Prompt, handler server.PromptHandlerFunc) {
	return mcp.NewPrompt(
			"configure_settings",
			mcp.WithPromptDescription("Prompt to guide through configuring settings for an index"),
			mcp.WithArgument("index_name", mcp.ArgumentDescription("The index to configure"), mcp.RequiredArgument()),
		), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			indexName := req.Params.Arguments["index_name"]

			promptText := fmt.Sprintf("First, retrieve the current settings of the index '%s' using the 'get_settings' tool. Then, ask the user what attributes they want to filter, sort, or search on, and configure those by calling the 'update_settings' tool with appropriate values.", indexName)

			return mcp.NewGetPromptResult(
				"Configure Settings Helper",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		}
}

func (p *Protocol) MultiSearchPrompt() (prompt mcp.Prompt, handler server.PromptHandlerFunc) {
	return mcp.NewPrompt(
			"multi_search_help",
			mcp.WithPromptDescription("Prompt to guide through performing a multi-index search in Meilisearch"),
			mcp.WithArgument("queries", mcp.ArgumentDescription("JSON array of queries containing indexUid and q (e.g. '[{\"indexUid\":\"movies\",\"q\":\"sci-fi\"}]')"), mcp.RequiredArgument()),
		), func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			queries := req.Params.Arguments["queries"]

			promptText := fmt.Sprintf("Please use the 'multi_search' tool with the following query array: %s. Parse the search results for each index and summarize the key matches.", queries)

			return mcp.NewGetPromptResult(
				"Multi Search Helper",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(promptText)),
				},
			), nil
		}
}
