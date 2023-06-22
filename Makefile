default:
	@go run main.go

dev: build
	@./bin/sht run examples/dev.sht

build:
	@go build -o bin/$(BINARY_NAME) cmd/sht.go

test:
	@go test ./...