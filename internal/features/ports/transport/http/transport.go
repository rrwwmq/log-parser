package ports_transport_http

import (
	"context"
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_transport_http_server "github.com/rrwwmq/log-parser/internal/core/transport/http/server"
)

type PortsHTTPHandler struct {
	portsService PortsService
}

type PortsService interface {
	GetPorts(ctx context.Context, nodeID int) ([]domain.Port, error)
}

func NewPortsHTTPHandler(portsService PortsService) *PortsHTTPHandler {
	return &PortsHTTPHandler{
		portsService: portsService,
	}
}

func (h *PortsHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/port/{node_id}",
			Handler: h.GetPorts,
		},
	}
}
