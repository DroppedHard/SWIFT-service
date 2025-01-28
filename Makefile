build: 
	@go build -o bin/swift-service.exe cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/swift-service.exe

lint: 
	@gofmt -s -w .

migrate:
	@go run cmd/migrate/main.go