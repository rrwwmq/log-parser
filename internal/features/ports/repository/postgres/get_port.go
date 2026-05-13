package ports_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *PortsRepository) GetPorts(ctx context.Context, nodeID int) ([]domain.Port, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, node_id, port_guid, port_num, port_state, lid
		FROM logparser.ports
		WHERE node_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var ports []domain.Port
	for rows.Next() {
		var m PortModel
		if err := rows.Scan(&m.ID, &m.NodeID, &m.PortGUID, &m.PortNum, &m.PortState, &m.LID); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		ports = append(ports, domain.NewPort(m.ID, m.NodeID, m.PortGUID, m.PortNum, m.PortState, m.LID))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return ports, nil
}
