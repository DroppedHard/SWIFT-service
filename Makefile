build: 
	@go build -o bin/swift-service.exe ./cmd/main.go

build-migrate: 
	@go build -o bin/swift-migrate.exe ./cmd/migrate

test:
	@go test -v ./...

run: build
	@./bin/swift-service.exe

run-migrate: build-migrate
	@./bin/swift-migrate.exe

lint: 
	@gofmt -s -w .

migrate:
	@go run cmd/migrate/main.go

build-docker:
	@docker build -t swift-service:latest .

run-docker: build-docker
	@docker run --rm -p 8080:8080 --name swift-service swift-service:latest

DOCKER_COMPOSE = docker-compose

compose-up:
	$(DOCKER_COMPOSE) up --build -d

compose-down:
	$(DOCKER_COMPOSE) down

compose-clean:
	$(DOCKER_COMPOSE) down -v