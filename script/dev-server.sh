#!/usr/bin/env bash
export DB_ADDRESS=postgres://localhost:5432/postgres?sslmode=disable
pg_ctl -D ./data -l logfile start

go run main.go