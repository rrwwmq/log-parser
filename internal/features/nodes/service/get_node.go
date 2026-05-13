package nodes_service

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (s *NodesService) GetNode(ctx context.Context, id int) (domain.Node, error) {
	node, err := s.nodesRepository.GetNode(ctx, id)
	if err != nil {
		return domain.Node{}, fmt.Errorf("get node: %w", err)
	}

	return node, nil
}
