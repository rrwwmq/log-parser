package ports_transport_http

import (
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_transport_http_request "github.com/rrwwmq/log-parser/internal/core/transport/http/request"
	core_transport_http_response "github.com/rrwwmq/log-parser/internal/core/transport/http/response"
)

type PortDTO struct {
	ID        int    `json:"id"`
	PortGUID  string `json:"port_guid"`
	PortNum   int    `json:"port_num"`
	PortState int    `json:"port_state"`
	LID       int    `json:"lid"`
}

type GetPortsResponse struct {
	Ports []PortDTO `json:"ports"`
}

func (h *PortsHTTPHandler) GetPorts(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetPorts handler")

	nodeID, err := core_transport_http_request.GetIntPathValue(r, "node_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get node_id from path")
		return
	}

	ports, err := h.portsService.GetPorts(ctx, nodeID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get ports")
		return
	}

	response := GetPortsResponse{
		Ports: portDTOsFromDomain(ports),
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}

func portDTOsFromDomain(ports []domain.Port) []PortDTO {
	dtos := make([]PortDTO, 0, len(ports))
	for _, p := range ports {
		dtos = append(dtos, PortDTO{
			ID:        p.ID,
			PortGUID:  p.PortGUID,
			PortNum:   p.PortNum,
			PortState: p.PortState,
			LID:       p.LID,
		})
	}
	return dtos
}
