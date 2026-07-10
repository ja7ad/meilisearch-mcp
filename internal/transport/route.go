package transport

import (
	"github.com/ja7ad/meilisearch-mcp/internal/protocol"
	"github.com/mark3labs/mcp-go/server"
)

type Route struct {
	mc          *server.MCPServer
	proto       *protocol.Protocol
	middlewares []ToolMiddleware
}

func NewRoute(mc *server.MCPServer, proto *protocol.Protocol, middlewares ...ToolMiddleware) *Route {
	return &Route{
		mc:          mc,
		proto:       proto,
		middlewares: middlewares,
	}
}

func (r *Route) Register() {
	r.registerIndexRoute()
	r.registerDocumentRoute()
	r.registerSearchRoute()
	r.registerSettingsRoute()
	r.registerTaskRoute()
	r.registerKeyRoute()
	r.registerSystemRoute()
	r.registerNetworkRoute()
	r.registerWebhookRoute()
	r.registerChatRoute()
	r.registerExperimentalRoute()
	r.registerResources()
	r.registerPrompts()
}

func (r *Route) apply(handler server.ToolHandlerFunc, mw ...ToolMiddleware) server.ToolHandlerFunc {
	if len(mw) == 0 {
		return handler
	}
	return ChainToolMiddleware(mw...)(handler)
}

func (r *Route) registerIndexRoute() {
	tool, handler := r.proto.CreateIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ListIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.SwapIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerDocumentRoute() {
	tool, handler := r.proto.GetDocument()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.AddDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteDocument()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteAllDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteDocumentsByFilter()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerSearchRoute() {
	tool, handler := r.proto.Search()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.FacetSearch()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.SearchSimilarDocuments()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.MultiSearch()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerSettingsRoute() {
	tool, handler := r.proto.GetSettings()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateSettings()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ResetSettings()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerTaskRoute() {
	tool, handler := r.proto.GetTask()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ListTasks()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.CancelTasks()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteTasks()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerKeyRoute() {
	tool, handler := r.proto.ListKeys()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetKey()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.CreateKey()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateKey()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteKey()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerSystemRoute() {
	tool, handler := r.proto.GetStats()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetHealth()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetVersion()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.CreateDump()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.CreateSnapshot()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerNetworkRoute() {
	tool, handler := r.proto.GetNetwork()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateNetwork()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerWebhookRoute() {
	tool, handler := r.proto.AddWebhook()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateWebhook()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.DeleteWebhook()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetWebhook()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ListWebhooks()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerChatRoute() {
	tool, handler := r.proto.ChatCompletion()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ListChatWorkspaces()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetChatWorkspace()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.GetChatWorkspaceSettings()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateChatWorkspace()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ResetChatWorkspace()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerExperimentalRoute() {
	tool, handler := r.proto.GetExperimentalFeatures()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.UpdateExperimentalFeatures()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerResources() {
	res, resHandler := r.proto.GetHealthResource()
	r.mc.AddResource(res, resHandler)

	res, resHandler = r.proto.GetVersionResource()
	r.mc.AddResource(res, resHandler)

	res, resHandler = r.proto.GetStatsResource()
	r.mc.AddResource(res, resHandler)

	res, resHandler = r.proto.GetIndexesResource()
	r.mc.AddResource(res, resHandler)

	tpl, tplHandler := r.proto.GetIndexStatsResourceTemplate()
	r.mc.AddResourceTemplate(tpl, tplHandler)

	tpl, tplHandler = r.proto.GetIndexSettingsResourceTemplate()
	r.mc.AddResourceTemplate(tpl, tplHandler)

	tpl, tplHandler = r.proto.GetIndexDocumentsResourceTemplate()
	r.mc.AddResourceTemplate(tpl, tplHandler)
}

func (r *Route) registerPrompts() {
	pr, prHandler := r.proto.SearchIndexPrompt()
	r.mc.AddPrompt(pr, prHandler)

	pr, prHandler = r.proto.ImportDocumentsPrompt()
	r.mc.AddPrompt(pr, prHandler)

	pr, prHandler = r.proto.ConfigureSettingsPrompt()
	r.mc.AddPrompt(pr, prHandler)
}
