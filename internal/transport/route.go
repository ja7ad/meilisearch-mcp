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
	r.registerTaskRoute()
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

	tool, handler = r.proto.GetIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))

	tool, handler = r.proto.ListIndex()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}

func (r *Route) registerTaskRoute() {
	tool, handler := r.proto.GetTask()
	r.mc.AddTool(tool, r.apply(handler, r.middlewares...))
}
