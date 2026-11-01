# Multi-stage multi-arch Dockerfile for FlagManagment Edge Proxy
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy backend for protobuf/SDK module dependencies if needed, or proxy files
COPY backend/go.mod backend/go.sum ./backend/
COPY proxy/go.mod proxy/go.sum* ./proxy/

WORKDIR /app/proxy
RUN go mod download || true

COPY backend/pkg/gen/sdk/v1 ../backend/pkg/gen/sdk/v1
COPY proxy/ ./

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -o /proxy ./cmd/proxy

FROM alpine:3.21

RUN apk update && apk upgrade --no-cache && apk add --no-cache ca-certificates tzdata

WORKDIR /
COPY --from=builder /proxy /proxy

EXPOSE 9091 8081

ENTRYPOINT ["/proxy"]
