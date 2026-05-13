package nodes_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_errors "github.com/rrwwmq/log-parser/internal/core/errors"
)

func (r *NodesRepository) GetNode(ctx context.Context, id int) (domain.Node, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			n.id, n.log_id, n.node_guid, n.node_desc, n.node_type, n.num_ports,
			i.id, i.serial_number, i.part_number, i.revision, i.product_name,
			i.endianness, i.enable_endianness_per_job, i.reproducibility_disable
		FROM logparser.nodes n
		LEFT JOIN logparser.nodes_info i ON i.node_id = n.id
		WHERE n.id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var m NodeModel
	var i NodeInfoModel
	err := row.Scan(
		&m.ID, &m.LogID, &m.NodeGUID, &m.NodeDesc, &m.NodeType, &m.NumPorts,
		&i.ID, &i.SerialNumber, &i.PartNumber, &i.Revision, &i.ProductName,
		&i.Endianness, &i.EnableEndiannessPerJob, &i.ReproducibilityDisable,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Node{}, fmt.Errorf("node id=%d: %w", id, core_errors.ErrNotFound)
		}
		return domain.Node{}, fmt.Errorf("scan: %w", err)
	}

	var info *domain.NodeInfo
	if i.ID != nil {
		info = &domain.NodeInfo{
			ID:                     *i.ID,
			NodeID:                 m.ID,
			SerialNumber:           i.SerialNumber,
			PartNumber:             i.PartNumber,
			Revision:               i.Revision,
			ProductName:            i.ProductName,
			Endianness:             i.Endianness,
			EnableEndiannessPerJob: i.EnableEndiannessPerJob,
			ReproducibilityDisable: i.ReproducibilityDisable,
		}
	}

	return domain.NewNode(m.ID, m.LogID, m.NodeGUID, m.NodeDesc, domain.NodeType(m.NodeType), m.NumPorts, info), nil
}
