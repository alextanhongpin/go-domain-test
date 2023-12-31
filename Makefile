mockery := go run github.com/vektra/mockery/cmd/mockery


all: test


install:
	@go install github.com/vektra/mockery/v2@v2.32.0


# usage: make mock name=purchaseRepository
mock:
	mockery


test:
	@go test -v -failfast -cover -covermode=count -coverprofile=cover.out ./...
	@go tool cover -html=cover.out
