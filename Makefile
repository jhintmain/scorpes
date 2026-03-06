.PHONY: build run dev test clean docker-dev docker-prod sqlc-install sqlc-version sqlc-generate

build:
	go build -o ./tmp/main ./cmd/server

run:
	go run ./cmd/server

dev:
	air

test:
	go test ./...

clean:
	rm -rf ./tmp ./api

docker-dev:
	docker compose up api-dev

sqlc-install:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc-version:
	sqlc version

sqlc-generate:
	sqlc generate
