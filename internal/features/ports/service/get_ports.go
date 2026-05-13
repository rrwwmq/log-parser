package ports_service

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (s *PortsService) GetPorts(ctx context.Context, nodeID int) ([]domain.Port, error) {
	ports, err := s.portsRepository.GetPorts(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("get ports: %w", err)
	}

	return ports, nil
}
