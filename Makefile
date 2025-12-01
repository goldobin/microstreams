default: fmt lint test build

fmt:
	@goimports -local "github.com/goldobin/microstreams" -l -w .
	@gofumpt -l -w .

lint:
	@golangci-lint run ./...

test:
	@go test ./...

build: build-pub build-sub

build-pub:
	@go build -o build/bin/pub cmd/pub/main.go

build-sub:
	@go build -o build/bin/sub cmd/sub/main.go

run-pub:
	@go run cmd/pub/main.go

run-sub:
	@go run cmd/sub/main.go $(ARGS)