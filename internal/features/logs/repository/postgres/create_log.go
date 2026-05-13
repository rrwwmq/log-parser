package logs_postgres_repository

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (r *LogsRepository) CreateLog(ctx context.Context, log domain.Log) (domain.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO logparser.logs(
			file_name,
			status,
			uploaded_at,
			node_count,
			port_count
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, file_name, status, uploaded_at, node_count, port_count;
	`

	row := r.pool.QueryRow(ctx, query,
		log.FileName,
		log.Status,
		log.UploadedAt,
		log.NodeCount,
		log.PortCount,
	)

	var m LogModel
	if err := row.Scan(&m.ID, &m.FileName, &m.Status, &m.UploadedAt, &m.NodeCount, &m.PortCount); err != nil {
		return domain.Log{}, fmt.Errorf("scan: %w", err)
	}

	return domain.NewLog(m.ID, m.FileName, domain.LogStatus(m.Status), m.UploadedAt, m.NodeCount, m.PortCount), nil
}
