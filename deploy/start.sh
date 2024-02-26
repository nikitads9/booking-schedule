#!/bin/sh

#migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
#TODO: toml and air
if [ "$ENV" = "development" ]; then
    echo "Starting backend service using hot reload"
    air
else
    echo "Starting backend service using build binary"
    /app/main
fi