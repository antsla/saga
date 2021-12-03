#!/bin/sh

until PGHOST=${HOST_DB} PGDATABASE=${POSTGRES_DB} PGPORT=${PORT_DB} PGUSER=${POSTGRES_USER} PGPASSWORD=${POSTGRES_PASSWORD} psql -c 'SELECT 1' >> /dev/null; do sleep 5; done;

migrate -source file://${PWD}/migrations/ -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST_DB}:${PORT_DB}/${POSTGRES_DB}?sslmode=disable up

go get -d ./...
go run -race ./cmd/main.go

eval "$@"
