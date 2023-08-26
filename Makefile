run:
	./main

migrate:
	goose -dir migrations postgres "postgres://postgres:228@postgres:5432/postgres?sslmode=disable" up

prod:
	make migrate	
	make run

env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.example .env
	@echo "AUTH_SALT=$$(openssl rand -hex 16)" >> .env
	@echo "AUTH_PRIVATE_KEY=$$(openssl rand -hex 32)" >> .env

start-db:
	docker compose -f db.docker-compose.yml up -d

test:
	go test -cover ./...
