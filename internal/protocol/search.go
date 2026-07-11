package protocol

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) Search() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"search",
			mcp.WithDescription("Search for documents in a Meilisearch index. Reference: https://www.meilisearch.com/docs/reference/api/search/search-with-get.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index to search in"), mcp.Required()),
			mcp.WithString("q", mcp.Description("Query string (optional)")),
			WithPagination(),
			mcp.WithArray("attributesToRetrieve", mcp.Description("Attributes to retrieve (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithArray("attributesToSearchOn", mcp.Description("Attributes to search on (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithString("filter", mcp.Description("Meilisearch filter string (optional), e.g. 'category = books'")),
			mcp.WithArray("sort", mcp.Description("Sort criteria, e.g. ['price:asc'] (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithArray("facets", mcp.Description("Facets to retrieve (optional)"),
				mcp.Items(map[string]any{"type": "string"}),
			),
			mcp.WithBoolean("showRankingScore", mcp.Description("Show ranking score (optional)")),
			mcp.WithBoolean("showRankingScoreDetails", mcp.Description("Show ranking score details (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			// Extract parameters using JSON marshal/unmarshal for convenience and correctness
			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var searchReq meilisearch.SearchRequest
			if err := sonic.Unmarshal(argsBytes, &searchReq); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if searchReq.Limit <= 0 {
				searchReq.Limit = 20
			}

			idx := client.Index(uid)
			resp, err := idx.SearchWithContext(ctx, searchReq.Query, &searchReq)
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

func (p *Protocol) FacetSearch() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"facet_search",
			mcp.WithDescription("Search for facet values in a Meilisearch index. Reference: https://www.meilisearch.com/docs/reference/api/facet-search/search-in-facets.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithString("facetName", mcp.Description("Name of the facet to search in"), mcp.Required()),
			mcp.WithString("facetQuery", mcp.Description("Query string for the facet (optional)")),
			mcp.WithString("q", mcp.Description("Main query string (optional)")),
			mcp.WithString("filter", mcp.Description("Filter to apply to search results (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
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

			var facetReq meilisearch.FacetSearchRequest
			if err := sonic.Unmarshal(argsBytes, &facetReq); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			idx := client.Index(uid)
			resp, err := idx.FacetSearchWithContext(ctx, &facetReq)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(string(*resp)), nil
		}
}

func (p *Protocol) SearchSimilarDocuments() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"search_similar_documents",
			mcp.WithDescription("Retrieve documents similar to a given document using AI/vectors. Reference: https://www.meilisearch.com/docs/reference/api/similar-documents/get-similar-documents-with-get.md"),
			mcp.WithString("index_name", mcp.Description("Name of the index"), mcp.Required()),
			mcp.WithString("id", mcp.Description("ID of the reference document"), mcp.Required()),
			mcp.WithString("embedder", mcp.Description("Name of the embedder to use (configured in index settings)"), mcp.Required()),
			mcp.WithNumber("limit", mcp.Description("Maximum number of similar documents to retrieve (default 20)")),
			mcp.WithNumber("offset", mcp.Description("Offset for pagination")),
			mcp.WithString("filter", mcp.Description("Filter results (optional)")),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			uid, err := RequiredParam[string](req, "index_name")
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

			var similarQuery meilisearch.SimilarDocumentQuery
			if err := sonic.Unmarshal(argsBytes, &similarQuery); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			if similarQuery.Limit <= 0 {
				similarQuery.Limit = 20
			}

			idx := client.Index(uid)
			var resp meilisearch.SimilarDocumentResult
			if err := idx.SearchSimilarDocumentsWithContext(ctx, &similarQuery, &resp); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			b, err := sonic.Marshal(resp)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(b)), nil
		}
}

func (p *Protocol) MultiSearch() (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool(
			"multi_search",
			mcp.WithDescription("Perform multiple search queries across one or more indexes in a single request. Reference: https://www.meilisearch.com/docs/reference/api/search/perform-a-multi-search.md"),
			mcp.WithArray("queries", mcp.Description("Array of search queries, each containing at least indexUid and q"), mcp.Required()),
		), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			argsBytes, err := sonic.Marshal(req.GetArguments())
			if err != nil {
				return nil, err
			}

			var multiReq meilisearch.MultiSearchRequest
			if err := sonic.Unmarshal(argsBytes, &multiReq); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			resp, err := client.MultiSearchWithContext(ctx, &multiReq)
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
