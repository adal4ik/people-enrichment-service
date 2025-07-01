include .env
export $(shell sed 's/=.*//' .env)

APP_NAME=people-enrichment-service

MIGRATIONS_PATH=./migrations
DB_URL=postgres://postgres:yourpassword@localhost:5432/peopledb?sslmode=disable

## Run latest migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

## Roll back last migration
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

## Force set DB version (use with caution!)
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force 1

## Drop entire DB schema (dangerous!)
migrate-drop:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop -f

## Show current migration version
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version


up:
	docker-compose up --build

down:
	docker-compose down

ps:
	docker-compose ps

logs:
	docker-compose logs -f

restart:
	docker-compose restart

db:
	docker-compose exec db psql -U $(DB_USER) $(DB_NAME)

build:
	go build -o ./bin/$(APP_NAME) ./cmd/app

run:
	go run ./cmd/app

help:
	@echo "Available commands:"
	@echo "  up      - Start the containers"
	@echo "  down    - Stop the containers"
	@echo "  ps      - List the containers"
	@echo "  db      - Connect to the PostgreSQL database"
	@echo "  logs    - Show logs for the containers"
	@echo "  restart - Restart the containers"
	@echo "  build   - Build Go binary"
	@echo "  run     - Run the app locally"
	@echo "  help    - Show this help message"
