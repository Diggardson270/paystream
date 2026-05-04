.PHONY: bootstrap dev build test lint

bootstrap:
	go mod download
	cd web/dashboard && pnpm install

build:
	go build -o bin/paystream-api    ./cmd/paystream-api
	go build -o bin/paystream-worker ./cmd/paystream-worker

dev:
	go run ./cmd/paystream-api &
	go run ./cmd/paystream-worker

test:
	go test ./...

lint:
	golangci-lint run ./...
