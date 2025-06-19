APP_NAME=people-enrichment-service

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
