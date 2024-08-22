ESL_BINARY=eslApp
MSSQL_BINARY=mssqlApp
REDIS_BINARY=redisApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker compose (if running), builds all projects and starts docker compose
up_build: build_esl build_sql
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_esl: builds the esl binary as a linux executable
build_esl:
	@echo "Building esl binary..."
	cd ../esl-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ESL_BINARY} ./cmd/api
	@echo "Done!"

## build_sql: builds the sql binary as a linux executable
build_sql:
	@echo "Building sql binary..."
	cd ../mssql-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MSSQL_BINARY} ./cmd/api
	@echo "Done!"

## build_redis: builds the redis binary as a linux executable
build_redis:
	@echo "Building redis binary..."
	cd ../redis-service && env GOOS=linux CGO_ENABLED=0 go build -o ${REDIS_BINARY} ./cmd/api
	@echo "Done!"
	