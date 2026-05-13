package ports_service

import (
	"context"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

type PortsService struct {
	portsRepository PortsRepository
}

func NewPortsService(portsRepository PortsRepository) *PortsService {
	return &PortsService{
		portsRepository: portsRepository,
	}
}

type PortsRepository interface {
	GetPorts(ctx context.Context, nodeID int) ([]domain.Port, error)
	CreatePort(ctx context.Context, port domain.Port) error
}
