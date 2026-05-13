package logs_service

import (
	"context"
	"fmt"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func (s *LogsService) GetLog(ctx context.Context, id int) (domain.Log, error) {
	logDomain, err := s.logsRepository.GetLog(ctx, id)
	if err != nil {
		return domain.Log{}, fmt.Errorf("get log: %w", err)
	}

	return logDomain, nil
}
