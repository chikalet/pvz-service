FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pvz-service ./cmd/server

FROM alpine:latest
WORKDIR /app

COPY --from=builder /pvz-service /app/pvz-service
COPY migrations /app/migrations

RUN apk add --no-cache postgresql-client

EXPOSE 8080

CMD ["/app/pvz-service"]
RUN apk add --no-cache postgresql-client curl
HEALTHCHECK --interval=30s --timeout=3s \
    CMD curl -f http://localhost:8080/health || exit 1