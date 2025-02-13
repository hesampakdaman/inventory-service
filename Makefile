.PHONY: build lint run test start docker-build stop docker-clean setup-topics

build:
	go build ./cmd/main.go

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

setup-topics:
	docker exec kafka \
	  kafka-topics --create \
	    --bootstrap-server kafka:9092 \
	    --topic inventory-service.commands \
	    --partitions 12 \
	    --replication-factor 1 || true
