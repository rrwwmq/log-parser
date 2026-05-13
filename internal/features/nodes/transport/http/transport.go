package nodes_transport_http

import (
	"context"
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_transport_http_server "github.com/rrwwmq/log-parser/internal/core/transport/http/server"
)

type NodesHTTPHandler struct {
	nodesService NodesService
}

type NodesService interface {
	GetTopology(ctx context.Context, logID int) ([]domain.Node, error)
	GetNode(ctx context.Context, id int) (domain.Node, error)
}

func NewNodesHTTPHandler(nodesService NodesService) *NodesHTTPHandler {
	return &NodesHTTPHandler{
		nodesService: nodesService,
	}
}

func (h *NodesHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/topology/{log_id}",
			Handler: h.GetTopology,
		},
		{
			Method:  http.MethodGet,
			Path:    "/node/{id}",
			Handler: h.GetNode,
		},
	}
}
