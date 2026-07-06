FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /app/api .
COPY --from=build /app/worker .
COPY --from=build /go/bin/goose /usr/local/bin/goose

COPY migrations ./migrations
