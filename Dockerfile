FROM golang:1.25.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o logparser ./cmd/log-parser/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/logparser .
COPY migrations/ ./migrations/
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./logparser"]