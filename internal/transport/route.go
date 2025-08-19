package transport

import (
	"github.com/ja7ad/meilisearch-mcp/internal/protocol"
	"github.com/mark3labs/mcp-go/server"
)

type Route struct {
	mc    *server.MCPServer
	proto *protocol.Protocol
}

func NewRoute(mc *server.MCPServer, proto *protocol.Protocol) *Route {
	return &Route{
		mc:    mc,
		proto: proto,
	}
}

func (r *Route) Register() {
	r.registerIndexRoute()
	r.registerTaskRoute()
}

func (r *Route) registerIndexRoute() {
	r.mc.AddTool(r.proto.CreateIndex())
}

func (r *Route) registerTaskRoute() {
	r.mc.AddTool(r.proto.GetTask())
}
