all: build

build:
	@echo building
	@go build -o main cmd/api/*

run:
	@go run cmd/api/*

watch:
	@air

mysql:
	@docker compose up mysql

psql:
	@docker compose up psql

down:
	@docker compose down	

dsn := $(shell cat dsn.txt)

goose_create:
	@goose -s -dir='./migrations' mysql "${dsn}" create "${fn}" sql

goose_one:
	@goose -dir='./migrations' mysql "${dsn}" up-by-one

goose_down:
	@goose -dir='./migrations' mysql "${dsn}" down

goose_up:
	@goose -dir='./migrations' mysql "${dsn}" up