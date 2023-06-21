default:
	@go run main.go

build:
	@go build -o bin/$(BINARY_NAME) cmd/sht.go

test:
	@go test ./...