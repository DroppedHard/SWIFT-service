FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN make build
RUN make migrate
RUN ls -l /app/bin

FROM golang:1.23
WORKDIR /app
COPY --from=builder /app/bin/swift-service.exe .
EXPOSE 8080
CMD ["./swift-service.exe"]