include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d log-parser-postgres

env-down:
	@docker compose down log-parser-postgres

env-cleanup:
	@read -p "Do you want to clean all volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down log-parser-postgres port-forwarder && rm -rf ${PROJECT_ROOT}/out/pgdata && echo "Files was cleaned"; \
	else \
		echo "clean up env was cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "example make migrate-create seq={your value}"; \
		exit 1; \
	fi; \
	docker compose run --rm log-parser-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "example make migrate-action action={your value}"; \
		exit 1; \
	fi; \
	docker compose run --rm log-parser-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@log-parser-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

logparser-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/logparser/main.go