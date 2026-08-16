FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY server /app/server
COPY migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/server"]
