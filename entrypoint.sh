#!/bin/sh

set -e

echo "Waiting for postgres to be ready..."
until pg_isready -h database -U "$POSTGRES_USER"; do
  sleep 1
done

echo "Running migrations..."
migrate -path /migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@database:5432/$POSTGRES_DB?sslmode=disable" up

echo "Starting the application..."
./app
