include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

deps:
	@ echo
	@ echo "Starting downloading dependencies..."
	@ echo
	@ go get -u ./...

mock:
	@ echo
	@ echo "Starting building mocks..."
	@ echo
	@ rm mocks/*.go || true && \
		$(GOPATH)/bin/mockgen -source=domain/contracts.go -destination=mocks/domain.go -package=mocks && \
		$(GOPATH)/bin/mockgen -source=infra/contracts.go -destination=mocks/infra.go -package=mocks

test:
	@ echo
	@ echo "Starting running tests..."
	@ echo
	@ go test -cover ./...

debug:
	@ echo
	@ echo "Starting the debug..."
	@ echo
	@ MONITOR_PORT=${WORKER_MONITOR_PORT} go run ./cmd/debug

%:
	@:
