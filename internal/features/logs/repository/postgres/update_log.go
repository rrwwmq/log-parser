package logs_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *LogsRepository) UpdateLog(ctx context.Context, log domain.Log) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE logparser.logs
		SET
			status = $1,
			node_count = $2,
			port_count = $3
		WHERE id = $4;
	`

	_, err := r.pool.Exec(ctx, query, log.Status, log.NodeCount, log.PortCount, log.ID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
