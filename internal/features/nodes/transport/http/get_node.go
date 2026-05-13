package nodes_transport_http

import (
	"net/http"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_transport_http_request "github.com/rrwwmq/log-parser/internal/core/transport/http/request"
	core_transport_http_response "github.com/rrwwmq/log-parser/internal/core/transport/http/response"
)

type NodeInfoDTO struct {
	SerialNumber           *string `json:"serial_number"`
	PartNumber             *string `json:"part_number"`
	Revision               *string `json:"revision"`
	ProductName            *string `json:"product_name"`
	Endianness             *int    `json:"endianness"`
	EnableEndiannessPerJob *int    `json:"enable_endianness_per_job"`
	ReproducibilityDisable *int    `json:"reproducibility_disable"`
}

type GetNodeResponse struct {
	ID       int          `json:"id"`
	NodeGUID string       `json:"node_guid"`
	NodeDesc string       `json:"node_desc"`
	NodeType string       `json:"node_type"`
	NumPorts int          `json:"num_ports"`
	Info     *NodeInfoDTO `json:"info"`
}

func (h *NodesHTTPHandler) GetNode(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetNode handler")

	id, err := core_transport_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get id from path")
		return
	}

	node, err := h.nodesService.GetNode(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get node")
		return
	}

	responseHandler.JSONResponse(getNodeResponseFromDomain(node), http.StatusOK)
}

func getNodeResponseFromDomain(n domain.Node) GetNodeResponse {
	resp := GetNodeResponse{
		ID:       n.ID,
		NodeGUID: n.NodeGUID,
		NodeDesc: n.NodeDesc,
		NodeType: "host",
		NumPorts: n.NumPorts,
	}

	if n.NodeType == domain.NodeTypeSwitch {
		resp.NodeType = "switch"
	}

	if n.Info != nil {
		resp.Info = &NodeInfoDTO{
			SerialNumber:           n.Info.SerialNumber,
			PartNumber:             n.Info.PartNumber,
			Revision:               n.Info.Revision,
			ProductName:            n.Info.ProductName,
			Endianness:             n.Info.Endianness,
			EnableEndiannessPerJob: n.Info.EnableEndiannessPerJob,
			ReproducibilityDisable: n.Info.ReproducibilityDisable,
		}
	}

	return resp
}
