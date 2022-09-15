#!make

GOPATH=$(shell go env GOPATH)

.make.setup:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen
	touch .make.setup

mock: .make.setup
	rm -f mocks/pubsubmocks.go
	$(GOPATH)/bin/mockgen -source=contracts.go -destination=mocks/pubsubmocks.go -package=mocks
	$(GOPATH)/bin/mockgen -source=subscriber_steps/contracts.go -destination=mocks/subscriberreceivermocks.go -package=mocks

test: mock
	go test ./... -cover

coverage:
	go test -coverprofile=cover.out $(path)
	go tool cover -func=cover.out -o cover || true
	rm cover.out
