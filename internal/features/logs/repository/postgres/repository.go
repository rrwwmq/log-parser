package logs_postgres_repository

import core_repository_postgres_pool "github.com/rrwwmq/log-parser/internal/core/repository/postgres/pool"

type LogsRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewLogsRepository(pool core_repository_postgres_pool.Pool) *LogsRepository {
	return &LogsRepository{
		pool: pool,
	}
}