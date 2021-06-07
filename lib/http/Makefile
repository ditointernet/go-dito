#!make

GOPATH=$(shell go env GOPATH)

.make.setup:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen
	touch .make.setup

mock: .make.setup
	rm -f mocks/inframocks.go
	$(GOPATH)/bin/mockgen -source=infra/contracts.go -destination=mocks/inframocks.go -package=mocks

test: mock
	go test ./... -cover

coverage:
	go test -coverprofile=cover.out $(path)
	go tool cover -func=cover.out -o cover || true
	rm cover.out
