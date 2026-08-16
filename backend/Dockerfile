FROM golang:1.24-bookworm AS builder

WORKDIR /app

ARG TARGETARCH

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w -X github.com/flagmanagment/backend/internal/health.Version=0.1.0" -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
