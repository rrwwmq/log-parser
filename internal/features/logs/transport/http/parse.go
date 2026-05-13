package logs_transport_http

import (
	"net/http"

	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_transport_http_request "github.com/rrwwmq/log-parser/internal/core/transport/http/request"
	core_transport_http_response "github.com/rrwwmq/log-parser/internal/core/transport/http/response"
)

type ParseRequest struct {
	FilePath string `json:"file_path" validate:"required"`
}

type ParseResposne struct {
	LogID int `json:"log_id"`
}

func (h *LogsHTTPHandler) Parse(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke Parse handler")

	var request ParseRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	logDomain, err := h.logsService.Parse(ctx, request.FilePath)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed tp parse log file")
		return
	}

	responseHandler.JSONResponse(ParseResposne{LogID: logDomain.ID}, http.StatusCreated)
}
