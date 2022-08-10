include $(shell git rev-parse --show-toplevel)/build/testenv.mk

.PHONY: all
all: test lint

.PHONY: version-check
version-check:
	./scripts/check-go-version.sh

.PHONY: deps-check
deps-check: version-check
	go mod tidy && go mod verify
	@make -s no-diff || { echo "go.mod and go.sum don't match the source code in the project. Please run \"go mod tidy && go mod verify\" locally and commit the changes"; exit 1; }

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: test
test: version-check generate
	go test -race -tags=integration -timeout=0 ./...

.PHONY: lint
lint: generate
	golangci-lint run
