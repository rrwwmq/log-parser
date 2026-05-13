package nodes_postgres_repository

import core_repository_postgres_pool "github.com/rrwwmq/log-parser/internal/core/repository/postgres/pool"

type NodesRepository struct {
	pool core_repository_postgres_pool.Pool
}

func NewNodesRepository(pool core_repository_postgres_pool.Pool) *NodesRepository {
	return &NodesRepository{
		pool: pool,
	}
}
