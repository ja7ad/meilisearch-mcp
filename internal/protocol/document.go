package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetDocument() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_document",
			mcp.WithDescription("Get a single document from Meilisearch by its ID. Reference: https://www.meilisearch.com/docs/reference/api/documents/get-document.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithString("document_id", mcp.Description("ID of the document to retrieve"), mcp.Required()),
			mcp.WithArray("fields", mcp.Description("List of attributes to retrieve (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			docID, err := RequiredParam[string](req, "document_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			fields, err := OptionalStringArrayParam(req, "fields")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			var query *meilisearch.DocumentQuery
			if len(fields) > 0 {
				query = &meilisearch.DocumentQuery{
					Fields: fields,
				}
			}

			var doc map[string]any
			if err := idx.GetDocumentWithContext(ctx, docID, query, &doc); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(doc)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) GetDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"get_documents",
			mcp.WithDescription("Retrieve multiple documents from the index. Reference: https://www.meilisearch.com/docs/reference/api/documents/list-documents-with-get.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithNumber("limit", mcp.Description("Maximum number of documents to retrieve (default 20, optional)"), mcp.Min(1)),
			mcp.WithNumber("offset", mcp.Description("Offset for pagination (default 0, optional)"), mcp.Min(0)),
			mcp.WithArray("fields", mcp.Description("List of attributes to retrieve (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithArray("ids", mcp.Description("Retrieve only documents with the specified IDs (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			limit, err := OptionalInt64Param(req, "limit")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			offset, err := OptionalInt64Param(req, "offset")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			fields, err := OptionalStringArrayParam(req, "fields")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			ids, err := OptionalStringArrayParam(req, "ids")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if limit <= 0 {
				limit = 20
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			query := &meilisearch.DocumentsQuery{
				Limit:  limit,
				Offset: offset,
			}
			if len(fields) > 0 {
				query.Fields = fields
			}
			if len(ids) > 0 {
				query.Ids = ids
			}

			var result meilisearch.DocumentsResult
			if err := idx.GetDocumentsWithContext(ctx, query, &result); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(result)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) AddDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"add_documents",
			mcp.WithDescription("Add or replace documents in the index (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithArray("documents", mcp.Description("Array of JSON objects representing documents"), mcp.Required()),
			mcp.WithString("primary_key", mcp.Description("Primary key for the index (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			docsVal, ok := req.GetArguments()["documents"]
			if !ok {
				return mcp.NewToolResultError("missing required parameter: documents"), nil
			}
			primaryKey, err := OptionalParam[string](req, "primary_key")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var opts *meilisearch.DocumentOptions
			if primaryKey != "" {
				opts = &meilisearch.DocumentOptions{
					PrimaryKey: &primaryKey,
				}
			}

			idx := client.Index(uid)
			task, err := idx.AddDocumentsWithContext(ctx, docsVal, opts)
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

func (p *Protocol) UpdateDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"update_documents",
			mcp.WithDescription("Partially update existing documents in the index (asynchronously). Reference: https://www.meilisearch.com/docs/reference/api/documents/add-or-replace-documents.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithArray("documents", mcp.Description("Array of JSON objects representing documents to update"), mcp.Required()),
			mcp.WithString("primary_key", mcp.Description("Primary key for the index (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			docsVal, ok := req.GetArguments()["documents"]
			if !ok {
				return mcp.NewToolResultError("missing required parameter: documents"), nil
			}
			primaryKey, err := OptionalParam[string](req, "primary_key")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			var opts *meilisearch.DocumentOptions
			if primaryKey != "" {
				opts = &meilisearch.DocumentOptions{
					PrimaryKey: &primaryKey,
				}
			}

			idx := client.Index(uid)
			task, err := idx.UpdateDocumentsWithContext(ctx, docsVal, opts)
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

func (p *Protocol) DeleteDocument() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_document",
			mcp.WithDescription("Delete a single document from Meilisearch. Reference: https://www.meilisearch.com/docs/reference/api/documents/delete-document.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithString("document_id", mcp.Description("ID of the document to delete"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			docID, err := RequiredParam[string](req, "document_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.DeleteDocumentWithContext(ctx, docID, nil)
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

func (p *Protocol) DeleteDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_documents",
			mcp.WithDescription("Delete multiple documents from Meilisearch by their IDs. Reference: https://www.meilisearch.com/docs/reference/api/documents/delete-document.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithArray("document_ids", mcp.Description("IDs of the documents to delete"), mcp.Required(),
				mcp.Items(map[string]any{"type": "string"}),
			),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			docIDs, err := RequiredStringArrayParam(req, "document_ids")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.DeleteDocumentsWithContext(ctx, docIDs, nil)
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

func (p *Protocol) DeleteAllDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_all_documents",
			mcp.WithDescription("Delete all documents from Meilisearch index. Reference: https://www.meilisearch.com/docs/reference/api/documents/delete-all-documents.md"),
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
			task, err := idx.DeleteAllDocumentsWithContext(ctx, nil)
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

func (p *Protocol) DeleteDocumentsByFilter() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"delete_documents_by_filter",
			mcp.WithDescription("Delete documents matching a filter. Reference: https://www.meilisearch.com/docs/reference/api/documents/delete-documents-by-filter.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithString("filter", mcp.Description("Meilisearch filter string, e.g., 'genre = horror'"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			filter, err := RequiredParam[string](req, "filter")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			task, err := idx.DeleteDocumentsByFilterWithContext(ctx, filter, nil)
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
