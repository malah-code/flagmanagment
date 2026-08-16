FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/flagmanagment/backend/internal/health.Version=0.1.0" -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
