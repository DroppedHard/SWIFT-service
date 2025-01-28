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

build-docker:
	@docker build -t swift-service:latest . --progress=plain

run-docker: build-docker
	@docker run --rm -p 8080:8080 --name swift-service swift-service:latest


DOCKER_COMPOSE = docker-compose
MIGRATION_IMAGE = migration-app:latest
REDIS_URL = redis://localhost:6379

compose-up:
	$(DOCKER_COMPOSE) up --build -d

compose-down:
	$(DOCKER_COMPOSE) down

compose-clean:
	$(DOCKER_COMPOSE) down -v