# Development shortcuts for Docker Compose
.PHONY: help up down restart logs ps build clean

help:
	@echo "make up       - Start MySQL and RabbitMQ"
	@echo "make down     - Stop containers"
	@echo "make restart  - Restart containers"
	@echo "make logs     - Follow service logs"
	@echo "make ps       - Show container status"
	@echo "make build    - Pull/build services and start them"
	@echo "make clean    - Stop containers and delete local volumes (DATA LOSS)"

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

logs:
	docker compose logs -f

ps:
	docker compose ps

build:
	docker compose up -d --build

clean:
	docker compose down -v
