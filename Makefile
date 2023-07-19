mockery := go run github.com/vektra/mockery/cmd/mockery


all: test


install:
	@go install github.com/vektra/mockery/v2@v2.32.0


mock:
ifndef name
	$(error 'name is required, make mock name=InterfaceToMock')
endif
	@$(mockery) -name $(name) -recursive -case underscore


test:
	@go test -v -failfast -cover -coverprofile=cover.out ./...
	@go tool cover -html=cover.out
