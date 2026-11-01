# Multi-stage build for backend
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.21

RUN apk update && apk upgrade --no-cache && apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY migrations /app/migrations

RUN chmod +x /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
