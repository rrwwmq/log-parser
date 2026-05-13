package logs_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_errors "github.com/rrwwmq/log-parser/internal/core/errors"
)

func (r *LogsRepository) GetLog(ctx context.Context, id int) (domain.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			id,
			file_name,
			status,
			uploaded_at,
			node_count,
			port_count
		FROM logparser.logs
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var m LogModel
	if err := row.Scan(&m.ID, &m.FileName, &m.Status, &m.UploadedAt, &m.NodeCount, &m.PortCount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Log{}, fmt.Errorf("log id=%d: %w", id, core_errors.ErrNotFound)
		}

		return domain.Log{}, fmt.Errorf("scan: %w", err)
	}

	return domain.NewLog(m.ID, m.FileName, domain.LogStatus(m.Status), m.UploadedAt, m.NodeCount, m.PortCount), nil
}
