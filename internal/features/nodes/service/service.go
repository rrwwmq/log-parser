package nodes_service

import (
	"context"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

type NodesService struct {
	nodesRepository NodesRepository
}

func NewNodesService(nodesRepository NodesRepository) *NodesService {
	return &NodesService{
		nodesRepository: nodesRepository,
	}
}

type NodesRepository interface {
	GetTopology(ctx context.Context, logID int) ([]domain.Node, error)
	GetNode(ctx context.Context, id int) (domain.Node, error)
}
