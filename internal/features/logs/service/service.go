package logs_service

import (
	"context"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

type LogsService struct {
	logsRepository LogsRepository
	nodesRepository NodesRepository
	portsRepository PortsRepository
}

func NewLogsService(logsRepository LogsRepository, nodesRepository NodesRepository, portsRepository PortsRepository) *LogsService {
	return &LogsService{
		logsRepository: logsRepository,
		nodesRepository: nodesRepository,
		portsRepository: portsRepository,
	}
}

type LogsRepository interface {
	CreateLog(ctx context.Context, log domain.Log) (domain.Log, error)
	UpdateLog(ctx context.Context, log domain.Log) error
	GetLog(ctx context.Context, id int) (domain.Log, error)
}

type NodesRepository interface {
	CreateNode(ctx context.Context, node domain.Node) (domain.Node, error)
	CreateNodeInfo(ctx context.Context, info domain.NodeInfo) error
}

type PortsRepository interface {
	CreatePort(ctx context.Context, port domain.Port) error
}