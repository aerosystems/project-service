.PHONY: proto lint lint-fix test help

##proto: generates proto files
#make sure to install the following dependencies
#go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
proto:
	@protoc --go_out=internal/common/protobuf --go-grpc_out=internal/common/protobuf -I api/protobuf project.proto

##lint-fix: runs linter with fix some issues
lint-fix:
	@golangci-lint run --fix
	@gofumpt -w -extra .

##lint: runs linter
lint:
	@golangci-lint run

##test: runs tests
test:
	@go test -v ./... -count=1 -cover

##help: displays help
help: Makefile
	@echo " Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
