package nodes_transport_http

import (
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_transport_http_request "github.com/rrwwmq/log-parser/internal/core/transport/http/request"
	core_transport_http_response "github.com/rrwwmq/log-parser/internal/core/transport/http/response"
)

type NodeDTO struct {
	ID       int    `json:"id"`
	NodeGUID string `json:"node_guid"`
	NodeDesc string `json:"node_desc"`
	NodeType string `json:"node_type"`
	NumPorts int    `json:"num_ports"`
}

type GetTopologyResponse struct {
	Nodes []NodeDTO `json:"nodes"`
}

func (h *NodesHTTPHandler) GetTopology(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetTopology handler")

	logID, err := core_transport_http_request.GetIntPathValue(r, "log_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get log_id from path")
		return
	}

	nodes, err := h.nodesService.GetTopology(ctx, logID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get topology")
		return
	}

	response := GetTopologyResponse{
		Nodes: nodeDTOsFromDomain(nodes),
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}

func nodeDTOsFromDomain(nodes []domain.Node) []NodeDTO {
	dtos := make([]NodeDTO, 0, len(nodes))
	for _, n := range nodes {
		dtos = append(dtos, nodeDTOFromDomain(n))
	}
	return dtos
}

func nodeDTOFromDomain(n domain.Node) NodeDTO {
	nodeType := "host"
	if n.NodeType == domain.NodeTypeSwitch {
		nodeType = "switch"
	}

	return NodeDTO{
		ID:       n.ID,
		NodeGUID: n.NodeGUID,
		NodeDesc: n.NodeDesc,
		NodeType: nodeType,
		NumPorts: n.NumPorts,
	}
}
