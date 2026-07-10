package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
)

// Mock ServiceManager implementation
type mockServiceManager struct {
	meilisearch.ServiceManager
	createIndexFunc func(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error)
	getIndexFunc    func(ctx context.Context, uid string) (*meilisearch.IndexResult, error)
	listIndexesFunc func(ctx context.Context, param *meilisearch.IndexesQuery) (*meilisearch.IndexesResults, error)
	indexFunc       func(uid string) meilisearch.IndexManager
	healthFunc      func(ctx context.Context) (*meilisearch.Health, error)
	versionFunc     func(ctx context.Context) (*meilisearch.Version, error)
	getStatsFunc    func(ctx context.Context, param *meilisearch.StatsParams) (*meilisearch.Stats, error)
	getNetworkFunc  func(ctx context.Context) (*meilisearch.Network, error)
	listWebhooksFunc func(ctx context.Context) (*meilisearch.WebhookResults, error)
}

func (m *mockServiceManager) CreateIndex(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error) {
	if m.createIndexFunc != nil {
		return m.createIndexFunc(config)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) CreateIndexWithContext(ctx context.Context, config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error) {
	return m.CreateIndex(config)
}

func (m *mockServiceManager) GetIndexWithContext(ctx context.Context, uid string) (*meilisearch.IndexResult, error) {
	if m.getIndexFunc != nil {
		return m.getIndexFunc(ctx, uid)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) ListIndexesWithContext(ctx context.Context, param *meilisearch.IndexesQuery) (*meilisearch.IndexesResults, error) {
	if m.listIndexesFunc != nil {
		return m.listIndexesFunc(ctx, param)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) Index(uid string) meilisearch.IndexManager {
	if m.indexFunc != nil {
		return m.indexFunc(uid)
	}
	return nil
}

func (m *mockServiceManager) HealthWithContext(ctx context.Context) (*meilisearch.Health, error) {
	if m.healthFunc != nil {
		return m.healthFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) VersionWithContext(ctx context.Context) (*meilisearch.Version, error) {
	if m.versionFunc != nil {
		return m.versionFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) GetStatsWithContext(ctx context.Context, param *meilisearch.StatsParams) (*meilisearch.Stats, error) {
	if m.getStatsFunc != nil {
		return m.getStatsFunc(ctx, param)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) GetNetworkWithContext(ctx context.Context) (*meilisearch.Network, error) {
	if m.getNetworkFunc != nil {
		return m.getNetworkFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockServiceManager) ListWebhooksWithContext(ctx context.Context) (*meilisearch.WebhookResults, error) {
	if m.listWebhooksFunc != nil {
		return m.listWebhooksFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

// Mock IndexManager implementation
type mockIndexManager struct {
	meilisearch.IndexManager
	getDocumentFunc  func(ctx context.Context, id string, query *meilisearch.DocumentQuery, dst interface{}) error
	getDocumentsFunc func(ctx context.Context, query *meilisearch.DocumentsQuery, resp *meilisearch.DocumentsResult) error
	addDocumentsFunc func(ctx context.Context, documents interface{}, opts *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error)
	searchFunc       func(ctx context.Context, query string, req *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error)
	getStatsFunc     func(ctx context.Context, param *meilisearch.StatsParams) (*meilisearch.StatsIndex, error)
	getSettingsFunc  func(ctx context.Context) (*meilisearch.Settings, error)
}

func (m *mockIndexManager) GetDocumentWithContext(ctx context.Context, id string, query *meilisearch.DocumentQuery, dst interface{}) error {
	if m.getDocumentFunc != nil {
		return m.getDocumentFunc(ctx, id, query, dst)
	}
	return errors.New("not implemented")
}

func (m *mockIndexManager) GetDocumentsWithContext(ctx context.Context, query *meilisearch.DocumentsQuery, resp *meilisearch.DocumentsResult) error {
	if m.getDocumentsFunc != nil {
		return m.getDocumentsFunc(ctx, query, resp)
	}
	return errors.New("not implemented")
}

func (m *mockIndexManager) AddDocumentsWithContext(ctx context.Context, documents interface{}, opts *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error) {
	if m.addDocumentsFunc != nil {
		return m.addDocumentsFunc(ctx, documents, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockIndexManager) SearchWithContext(ctx context.Context, query string, req *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query, req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockIndexManager) GetStatsWithContext(ctx context.Context, param *meilisearch.StatsParams) (*meilisearch.StatsIndex, error) {
	if m.getStatsFunc != nil {
		return m.getStatsFunc(ctx, param)
	}
	return nil, errors.New("not implemented")
}

func (m *mockIndexManager) GetSettingsWithContext(ctx context.Context) (*meilisearch.Settings, error) {
	if m.getSettingsFunc != nil {
		return m.getSettingsFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestCreateIndex(t *testing.T) {
	mockSrv := &mockServiceManager{
		createIndexFunc: func(config *meilisearch.IndexConfig) (*meilisearch.TaskInfo, error) {
			assert.Equal(t, "movies", config.Uid)
			assert.Equal(t, "id", config.PrimaryKey)
			return &meilisearch.TaskInfo{TaskUID: 123}, nil
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.CreateIndex()
	req := mcp.CallToolRequest{}
	req.Params.Name = "create_index"
	req.Params.Arguments = map[string]any{
		"index_name":  "movies",
		"primary_key": "id",
	}

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	var taskInfo meilisearch.TaskInfo
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &taskInfo)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), taskInfo.TaskUID)
}

func TestGetDocument(t *testing.T) {
	mockIndex := &mockIndexManager{
		getDocumentFunc: func(ctx context.Context, id string, query *meilisearch.DocumentQuery, dst interface{}) error {
			assert.Equal(t, "456", id)
			assert.Equal(t, []string{"title"}, query.Fields)
			m := dst.(*map[string]any)
			*m = map[string]any{"id": "456", "title": "Inception"}
			return nil
		},
	}
	mockSrv := &mockServiceManager{
		indexFunc: func(uid string) meilisearch.IndexManager {
			assert.Equal(t, "movies", uid)
			return mockIndex
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.GetDocument()
	req := mcp.CallToolRequest{}
	req.Params.Name = "get_document"
	req.Params.Arguments = map[string]any{
		"index_name":  "movies",
		"document_id": "456",
		"fields":      []any{"title"},
	}

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	var doc map[string]any
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &doc)
	assert.NoError(t, err)
	assert.Equal(t, "Inception", doc["title"])
}

func TestSearch(t *testing.T) {
	mockIndex := &mockIndexManager{
		searchFunc: func(ctx context.Context, query string, req *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
			assert.Equal(t, "sci-fi", query)
			assert.Equal(t, int64(5), req.Limit)
			return &meilisearch.SearchResponse{
				Hits: []meilisearch.Hit{
					{"id": json.RawMessage(`"1"`)},
				},
				Limit: 5,
			}, nil
		},
	}
	mockSrv := &mockServiceManager{
		indexFunc: func(uid string) meilisearch.IndexManager {
			assert.Equal(t, "movies", uid)
			return mockIndex
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.Search()
	req := mcp.CallToolRequest{}
	req.Params.Name = "search"
	req.Params.Arguments = map[string]any{
		"index_name": "movies",
		"q":          "sci-fi",
		"limit":      5,
	}

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	var searchResp meilisearch.SearchResponse
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &searchResp)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), searchResp.Limit)
}

func TestGetHealthResource(t *testing.T) {
	mockSrv := &mockServiceManager{
		healthFunc: func(ctx context.Context) (*meilisearch.Health, error) {
			return &meilisearch.Health{Status: "available"}, nil
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.GetHealthResource()
	req := mcp.ReadResourceRequest{}
	req.Params.URI = "meilisearch://health"

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	contents, ok := mcp.AsTextResourceContents(res[0])
	assert.True(t, ok)
	assert.Equal(t, "application/json", contents.MIMEType)

	var health meilisearch.Health
	err = sonic.Unmarshal([]byte(contents.Text), &health)
	assert.NoError(t, err)
	assert.Equal(t, "available", health.Status)
}

func TestSearchIndexPrompt(t *testing.T) {
	proto := New(TransportStdio, "http://localhost:7700", "key", nil)

	_, handler := proto.SearchIndexPrompt()
	req := mcp.GetPromptRequest{}
	req.Params.Name = "search_index"
	req.Params.Arguments = map[string]string{
		"index_name": "movies",
		"query":      "action",
	}

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Search Index Helper", res.Description)
	assert.Len(t, res.Messages, 1)
	assert.Equal(t, mcp.RoleUser, res.Messages[0].Role)

	textContents, ok := mcp.AsTextContent(res.Messages[0].Content)
	assert.True(t, ok)
	assert.Contains(t, textContents.Text, "movies")
	assert.Contains(t, textContents.Text, "action")
}

func TestGetNetwork(t *testing.T) {
	mockSrv := &mockServiceManager{
		getNetworkFunc: func(ctx context.Context) (*meilisearch.Network, error) {
			return &meilisearch.Network{
				Self:    "node1",
				Leader:  "node1",
				Version: "1.0.0",
			}, nil
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.GetNetwork()
	req := mcp.CallToolRequest{}
	req.Params.Name = "get_network"

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	var net meilisearch.Network
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &net)
	assert.NoError(t, err)
	assert.Equal(t, "node1", net.Self)
	assert.Equal(t, "node1", net.Leader)
}

func TestListWebhooks(t *testing.T) {
	mockSrv := &mockServiceManager{
		listWebhooksFunc: func(ctx context.Context) (*meilisearch.WebhookResults, error) {
			return &meilisearch.WebhookResults{
				Result: []*meilisearch.Webhook{
					{UUID: "w1", URL: "http://webhook.internal"},
				},
			}, nil
		},
	}

	proto := New(TransportStdio, "http://localhost:7700", "key", nil)
	proto.service = mockSrv

	_, handler := proto.ListWebhooks()
	req := mcp.CallToolRequest{}
	req.Params.Name = "list_webhooks"

	res, err := handler(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	var webhooks meilisearch.WebhookResults
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &webhooks)
	assert.NoError(t, err)
	assert.Len(t, webhooks.Result, 1)
	assert.Equal(t, "w1", webhooks.Result[0].UUID)
	assert.Equal(t, "http://webhook.internal", webhooks.Result[0].URL)
}
