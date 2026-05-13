package logs_transport_http

import (
	"context"
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_transport_http_server "github.com/rrwwmq/log-parser/internal/core/transport/http/server"
)

type LogsHTTPHandler struct {
	logsService LogsService
}

type LogsService interface {
	Parse(ctx context.Context, filePath string) (domain.Log, error)
	GetLog(ctx context.Context, id int) (domain.Log, error)
}

func NewLogsHTTPHandler(logsService LogsService) *LogsHTTPHandler {
	return &LogsHTTPHandler{
		logsService: logsService,
	}
}

func (h *LogsHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/parse",
			Handler: h.Parse,
		},

		{
			Method:  http.MethodGet,
			Path:    "/log/{id}",
			Handler: h.GetLog,
		},
	}
}
