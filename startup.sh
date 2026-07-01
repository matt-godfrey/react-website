#!/bin/bash

docker compose up -d

# backend
cd backend
go run ./cmd/worker/*.go &
go run ./cmd/api/*.go &

# frontend
cd ../frontend && npm run dev
