package nodes_service

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (s *NodesService) GetTopology(ctx context.Context, logID int) ([]domain.Node, error) {
	nodes, err := s.nodesRepository.GetTopology(ctx, logID)
	if err != nil {
		return nil, fmt.Errorf("get topology: %w", err)
	}

	return nodes, nil
}
