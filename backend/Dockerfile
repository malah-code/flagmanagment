FROM golang:alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/flagmanagment/backend/internal/health.Version=0.1.0" -o /app/server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server /server
COPY --from=builder /app/migrations /migrations

EXPOSE 8080
ENTRYPOINT ["/server"]
