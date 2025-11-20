default: fmt lint test

fmt:
	@goimports -local "github.com/goldobin/microstreams" -l -w .
	@gofumpt -l -w .

lint:
	@golangci-lint run ./...

test:
	@go test ./...

run-pub:
	@go run cmd/pub/main.go

run-sub:
	@go run cmd/sub/main.go $(ARGS)