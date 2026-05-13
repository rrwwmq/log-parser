package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/rrwwmq/log-parser/internal/core/logger"
	core_repository_postgres_pool "github.com/rrwwmq/log-parser/internal/core/repository/postgres/pool"
	core_transport_http_middleware "github.com/rrwwmq/log-parser/internal/core/transport/http/middleware"
	core_transport_http_server "github.com/rrwwmq/log-parser/internal/core/transport/http/server"
	logs_postgres_repository "github.com/rrwwmq/log-parser/internal/features/logs/repository/postgres"
	logs_service "github.com/rrwwmq/log-parser/internal/features/logs/service"
	logs_transport_http "github.com/rrwwmq/log-parser/internal/features/logs/transport/http"
	nodes_postgres_repository "github.com/rrwwmq/log-parser/internal/features/nodes/repository/postgres"
	nodes_service "github.com/rrwwmq/log-parser/internal/features/nodes/service"
	nodes_transport_http "github.com/rrwwmq/log-parser/internal/features/nodes/transport/http"
	ports_postgres_repository "github.com/rrwwmq/log-parser/internal/features/ports/repository/postgres"
	ports_service "github.com/rrwwmq/log-parser/internal/features/ports/service"
	ports_transport_http "github.com/rrwwmq/log-parser/internal/features/ports/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Initializing postgres connection pool")
	pool, err := core_repository_postgres_pool.NewConnectionPool(ctx, core_repository_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logsRepository := logs_postgres_repository.NewLogsRepository(pool)
	nodesRepository := nodes_postgres_repository.NewNodesRepository(pool)
	portsRepository := ports_postgres_repository.NewPortsRepository(pool)

	logger.Debug("initializing feature", zap.String("feature", "logs"))
	logsService := logs_service.NewLogsService(logsRepository, nodesRepository, portsRepository)

	logger.Debug("initializing feature", zap.String("feature", "nodes"))
	nodesService := nodes_service.NewNodesService(nodesRepository)

	logger.Debug("initializing feature", zap.String("feature", "ports"))
	portsService := ports_service.NewPortsService(portsRepository)

	logsTransportHTTP := logs_transport_http.NewLogsHTTPHandler(logsService)
	nodesTransportHTTP := nodes_transport_http.NewNodesHTTPHandler(nodesService)
	portsTransportHTTP := ports_transport_http.NewPortsHTTPHandler(portsService)

	logger.Debug("Initializing HTTP server")
	httpServer := core_transport_http_server.NewHTTPServer(
		core_transport_http_server.NewConfigMust(),
		logger,
		core_transport_http_middleware.RequestID(),
		core_transport_http_middleware.Logger(logger),
		core_transport_http_middleware.Trace(),
		core_transport_http_middleware.Panic(),
	)

	apiVersionRouter := core_transport_http_server.NewAPIVersionRouter(core_transport_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(logsTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(nodesTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(portsTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
