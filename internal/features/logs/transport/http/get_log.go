package logs_transport_http

import (
	"net/http"
	"time"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_transport_http_request "github.com/rrwwmq/log-parser/internal/core/transport/http/request"
	core_transport_http_response "github.com/rrwwmq/log-parser/internal/core/transport/http/response"
)

type GetLogResponse struct {
	ID         int       `json:"id"`
	FileName   string    `json:"file_name"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	NodeCount  int       `json:"node_count"`
	PortCount  int       `json:"parse_count"`
}

func (h *LogsHTTPHandler) GetLog(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetLog handler")

	id, err := core_transport_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from request")
		return
	}

	logDomain, err := h.logsService.GetLog(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get log")
	}

	responseHandler.JSONResponse(logDTOFromDomain(logDomain), http.StatusOK)
}

func logDTOFromDomain(log domain.Log) GetLogResponse {
	return GetLogResponse{
		ID:         log.ID,
		FileName:   log.FileName,
		Status:     string(log.Status),
		UploadedAt: log.UploadedAt,
		NodeCount:  log.NodeCount,
		PortCount:  log.PortCount,
	}
}
