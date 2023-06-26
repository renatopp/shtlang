default:
	@go run main.go

dev: build
	@./bin/sht run examples/dev.sht

repl:
	@go run cmd/sht.go

build:
	@go build -o bin/$(BINARY_NAME) cmd/sht.go

test:
	@go test ./...