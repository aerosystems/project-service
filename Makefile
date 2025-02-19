.PHONY: lint lint-fix test help

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
