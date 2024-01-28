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

migrate_create:
	@migrate create -seq -ext=.sql -dir=./migrations $(filename)

dsn := $(shell cat dsn.txt)

migrate_up:
	@migrate -path=./migrations -database="${dsn}" up

migrate_down:
	@migrate -path=./migrations -database="${dsn}" down

migrate_force:
	@migrate -path=./migrations -database="${dsn}" force "${n}"