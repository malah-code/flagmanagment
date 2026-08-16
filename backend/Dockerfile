FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY server /app/server
COPY migrations /app/migrations

RUN chmod +x /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
