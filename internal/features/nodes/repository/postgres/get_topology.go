package nodes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *NodesRepository) GetTopology(ctx context.Context, logID int) ([]domain.Node, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT 
			id, 
			log_id, 
			node_guid, 
			node_desc, 
			node_type, 
			num_ports
		FROM logparser.nodes
		WHERE log_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, logID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var m NodeModel
		if err := rows.Scan(&m.ID, &m.LogID, &m.NodeGUID, &m.NodeDesc, &m.NodeType, &m.NumPorts); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		nodes = append(nodes, domain.NewNode(m.ID, m.LogID, m.NodeGUID, m.NodeDesc, domain.NodeType(m.NodeType), m.NumPorts, nil))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return nodes, nil
}
