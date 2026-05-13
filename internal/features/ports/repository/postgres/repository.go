package ports_postgres_repository

import core_repository_postgres_pool "github.com/rrwwmq/log-parser/internal/core/repository/postgres/pool"

type PortsRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewPortsRepository(pool core_repository_postgres_pool.Pool) *PortsRepository {
	return &PortsRepository{
		pool: pool,
	}
}
