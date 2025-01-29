FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN make build
RUN make build-migrate

FROM golang:1.23
WORKDIR /app
COPY --from=builder /app/bin/swift-service.exe .
COPY --from=builder /app/bin/swift-migrate.exe .

EXPOSE 8080
CMD ["/bin/sh", "-c", "./swift-migrate.exe && ./swift-service.exe"]