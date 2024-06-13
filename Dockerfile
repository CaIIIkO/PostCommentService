# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk --no-cache add bash git make wget

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN go build -o ./bin/app cmd/main.go

# Install golang-migrate
RUN wget -O /usr/local/bin/migrate https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz
RUN tar -xzf /usr/local/bin/migrate -C /usr/local/bin
RUN chmod +x /usr/local/bin/migrate

# Stage 2: Runner
FROM alpine AS runner

WORKDIR /app
RUN apk --no-cache add bash postgresql-client

COPY --from=builder /app/bin/app /app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY internal/db/migrations/migrationsDocker /migrations
COPY docker.env .env
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
