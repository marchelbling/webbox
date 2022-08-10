POSTGRES_SERVICE_PORT := $(or ${POSTGRES_SERVICE_PORT}, 5433)
REDIS_SERVICE_PORT := $(or ${REDIS_SERVICE_PORT}, 6379)
REPO_ROOT := $(shell git rev-parse --show-toplevel)

.PHONY: testenv-up
testenv-up: redis-up postgres-up
	docker ps --all

.PHONY: testenv-down
testenv-down: redis-down postgres-down

.PHONY: redis-up
redis-up:
	docker top redis>/dev/null 2>&1 || \
		( docker rm -f redis; docker run --name redis --publish ${REDIS_SERVICE_PORT}:6379 --detach redis:5 )

.PHONY: redis-down
redis-down:
	-docker rm -f redis

.PHONY: postgres-up
postgres-up:
	docker top postgres>/dev/null 2>&1 || \
		( docker rm -f postgres; docker run --name postgres --publish ${POSTGRES_SERVICE_PORT}:5432 --detach postgres:11.1-alpine )

.PHONY: postgres-down
postgres-down:
	-docker rm -f postgres
