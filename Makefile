.PHONY: build run dev test clean docker-dev docker-prod

build:
	go build -o ./tmp/main ./cmd/api

run:
	go run ./cmd/api

dev:
	air

test:
	go test ./...

clean:
	rm -rf ./tmp ./api

docker-dev:
	docker compose up api-dev

docker-prod:
	docker compose up api
