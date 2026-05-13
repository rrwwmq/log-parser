package nodes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *NodesRepository) CreateNode(ctx context.Context, node domain.Node) (domain.Node, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO logparser.nodes (
			log_id, 
			node_guid, 
			node_desc, 
			node_type, 
			num_ports)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, log_id, node_guid, node_desc, node_type, num_ports;
	`

	row := r.pool.QueryRow(ctx, query,
		node.LogID,
		node.NodeGUID,
		node.NodeDesc,
		node.NodeType,
		node.NumPorts,
	)

	var m NodeModel
	if err := row.Scan(&m.ID, &m.LogID, &m.NodeGUID, &m.NodeDesc, &m.NodeType, &m.NumPorts); err != nil {
		return domain.Node{}, fmt.Errorf("scan: %w", err)
	}

	return domain.NewNode(m.ID, m.LogID, m.NodeGUID, m.NodeDesc, domain.NodeType(m.NodeType), m.NumPorts, nil), nil
}

func (r *NodesRepository) CreateNodeInfo(ctx context.Context, info domain.NodeInfo) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO logparser.nodes_info (
			node_id, 
			serial_number, 
			part_number, 
			revision, 
			product_name, 
			endianness, 
			enable_endianness_per_job, 
			reproducibility_disable)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := r.pool.Exec(ctx, query,
		info.NodeID,
		info.SerialNumber,
		info.PartNumber,
		info.Revision,
		info.ProductName,
		info.Endianness,
		info.EnableEndiannessPerJob,
		info.ReproducibilityDisable,
	)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
