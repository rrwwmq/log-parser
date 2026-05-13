package logs_service

import (
	"context"
	"fmt"
	"time"

	"github.com/rrwwmq/log-parser/internal/core/domain"
	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	"github.com/rrwwmq/log-parser/internal/parser"
	"go.uber.org/zap"
)

func (s *LogsService) Parse(ctx context.Context, filePath string) (domain.Log, error) {
	log := core_logger.FromContext(ctx)

	logDomain := domain.NewUninitializedLog(filePath)
	logDomain, err := s.logsRepository.CreateLog(ctx, logDomain)
	if err != nil {
		return domain.Log{}, fmt.Errorf("create log: %w", err)
	}

	start := time.Now()
	log.Debug("start parsing", zap.String("file_path", filePath))

	parsed, err := parser.ParseZip(filePath)
	if err != nil {
		logDomain.Status = domain.LogStatusFailed
		_ = s.logsRepository.UpdateLog(ctx, logDomain)
		return domain.Log{}, fmt.Errorf("parse zip: %w", err)
	}

	log.Debug(
		"parsed zip",
		zap.Duration("duration", time.Since(start)),
		zap.Int("nodes", len(parsed.Nodes)),
		zap.Int("ports", len(parsed.Ports)),
	)

	guidToNodeID := map[string]int{}
	for _, node := range parsed.Nodes {
		node.LogID = logDomain.ID
		savedNode, err := s.nodesRepository.CreateNode(ctx, node)
		if err != nil {
			logDomain.Status = domain.LogStatusFailed
			_ = s.logsRepository.UpdateLog(ctx, logDomain)
			return domain.Log{}, fmt.Errorf("create node info: %w", err)
		}

		guidToNodeID[node.NodeGUID] = savedNode.ID
	}

	for _, port := range parsed.Ports {
		nodeID, ok := guidToNodeID[port.PortGUID]
		if !ok {
			continue
		}

		port.NodeID = nodeID

		if err := s.portsRepository.CreatePort(ctx, port); err != nil {
			logDomain.Status = domain.LogStatusFailed
			_ = s.logsRepository.UpdateLog(ctx, logDomain)
			return domain.Log{}, fmt.Errorf("create port: %w", err)
		}
	}

	logDomain.Status = domain.LogStatusDone
	logDomain.NodeCount = len(parsed.Nodes)
	logDomain.PortCount = len(parsed.Ports)
	if err := s.logsRepository.UpdateLog(ctx, logDomain); err != nil {
		return domain.Log{}, fmt.Errorf("update log: %w", err)
	}

	return logDomain, nil
}
