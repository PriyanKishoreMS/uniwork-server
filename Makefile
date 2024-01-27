all: build

build:
	@echo building
	@go build -o main cmd/api/*

run:
	@go run cmd/api/*

watch:
	@air

up:
	@docker compose up

down:
	@docker compose down	