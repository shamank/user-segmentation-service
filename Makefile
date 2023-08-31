include .env
export $(shell sed 's/=.*//' .env)
.PHONY:
.SILENT:

build:
	go mod download && go build -o ./.bin/app ./cmd/main.go

run: build
	./.bin/app

r:
	go run ./cmd/main.go

docker:
	docker-compose build && docker-compose up

make swag:
	swag init -g internal/app/app.go

migrate-up:
	migrate -path ./migrations -database 'postgres://pguser:${PG_PASSWORD}@localhost:5432/devdb?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://pguser:${PG_PASSWORD}@localhost:5432/devdb?sslmode=disable' down

