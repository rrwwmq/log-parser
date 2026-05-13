package ports_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *PortsRepository) CreatePort(ctx context.Context, port domain.Port) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO logparser.ports (node_id, port_guid, port_num, port_state, lid)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.pool.Exec(ctx, query,
		port.NodeID,
		port.PortGUID,
		port.PortNum,
		port.PortState,
		port.LID,
	)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
