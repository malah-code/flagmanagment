FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY server /app/server
COPY migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/server"]
