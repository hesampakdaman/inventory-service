.PHONY: build lint run test start docker-build stop docker-clean

build:
	go build -o banking-service ./cmd/main.go

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

run:
	go run ./cmd/main.go

test:
	go test ./...

docker-build:
	docker build .

start:
	docker compose up --build -d

stop:
	docker compose down

docker-clean:
	docker compose down --volumes
