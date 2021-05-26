#!make

GOPATH=$(shell go env GOPATH)

setup:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen

mock:
	rm -f mocks/inframocks.go
	$(GOPATH)/bin/mockgen -source=contracts.go -destination=mocks/authmocks.go -package=mocks

test: mock
	go test ./... -cover

coverage:
	go test -coverprofile=cover.out $(path)
	go tool cover -func=cover.out -o cover || true
	rm cover.out
