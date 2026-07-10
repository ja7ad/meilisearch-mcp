package protocol

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/meilisearch/meilisearch-go"
)

func (p *Protocol) GetHealthResource() (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	return mcp.NewResource(
			"meilisearch://health",
			"Meilisearch Health Status",
			mcp.WithMIMEType("application/json"),
			mcp.WithResourceDescription("Get the health status of the Meilisearch instance"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			health, err := client.HealthWithContext(ctx)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(health)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "meilisearch://health",
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetVersionResource() (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	return mcp.NewResource(
			"meilisearch://version",
			"Meilisearch Version Info",
			mcp.WithMIMEType("application/json"),
			mcp.WithResourceDescription("Get the version of the Meilisearch instance"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			ver, err := client.VersionWithContext(ctx)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(ver)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "meilisearch://version",
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetStatsResource() (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	return mcp.NewResource(
			"meilisearch://stats",
			"Meilisearch Global Stats",
			mcp.WithMIMEType("application/json"),
			mcp.WithResourceDescription("Get global database statistics"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			stats, err := client.GetStatsWithContext(ctx, nil)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(stats)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "meilisearch://stats",
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetIndexesResource() (resource mcp.Resource, handler server.ResourceHandlerFunc) {
	return mcp.NewResource(
			"meilisearch://indexes",
			"Meilisearch Indexes List",
			mcp.WithMIMEType("application/json"),
			mcp.WithResourceDescription("List all indexes in Meilisearch"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			indexes, err := client.ListIndexesWithContext(ctx, nil)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(indexes)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      "meilisearch://indexes",
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetIndexStatsResourceTemplate() (template mcp.ResourceTemplate, handler server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"meilisearch://indexes/{index_name}/stats",
			"Meilisearch Index Stats Template",
			mcp.WithTemplateMIMEType("application/json"),
			mcp.WithTemplateDescription("Get stats for a specific index by name"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			indexName, ok := req.Params.Arguments["index_name"].(string)
			if !ok || indexName == "" {
				return nil, fmt.Errorf("missing required template argument: index_name")
			}

			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			idx := client.Index(indexName)
			stats, err := idx.GetStatsWithContext(ctx, nil)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(stats)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetIndexSettingsResourceTemplate() (template mcp.ResourceTemplate, handler server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"meilisearch://indexes/{index_name}/settings",
			"Meilisearch Index Settings Template",
			mcp.WithTemplateMIMEType("application/json"),
			mcp.WithTemplateDescription("Get settings for a specific index by name"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			indexName, ok := req.Params.Arguments["index_name"].(string)
			if !ok || indexName == "" {
				return nil, fmt.Errorf("missing required template argument: index_name")
			}

			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			idx := client.Index(indexName)
			settings, err := idx.GetSettingsWithContext(ctx)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(settings)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}

func (p *Protocol) GetIndexDocumentsResourceTemplate() (template mcp.ResourceTemplate, handler server.ResourceTemplateHandlerFunc) {
	return mcp.NewResourceTemplate(
			"meilisearch://indexes/{index_name}/documents",
			"Meilisearch Index Documents Template",
			mcp.WithTemplateMIMEType("application/json"),
			mcp.WithTemplateDescription("Get documents from a specific index by name (up to 20 documents)"),
		), func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			indexName, ok := req.Params.Arguments["index_name"].(string)
			if !ok || indexName == "" {
				return nil, fmt.Errorf("missing required template argument: index_name")
			}

			client, err := p.client(req.Header)
			if err != nil {
				return nil, err
			}

			idx := client.Index(indexName)
			var docs meilisearch.DocumentsResult
			err = idx.GetDocumentsWithContext(ctx, &meilisearch.DocumentsQuery{Limit: 20}, &docs)
			if err != nil {
				return nil, err
			}

			b, err := sonic.Marshal(docs)
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		}
}
