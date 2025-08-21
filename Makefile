.PHONY: up down logs ps restart build-images db-shell help

COMPOSE ?= docker compose

help:
	@echo "Repo-level Makefile"
	@echo "Targets:"
	@echo "  up            Build and start all services with docker compose"
	@echo "  down          Stop and remove containers and volumes"
	@echo "  logs          Tail logs for all services"
	@echo "  ps            Show compose service status"
	@echo "  restart       Restart app services (consumer, streamer)"
	@echo "  build-images  Build images via compose"
	@echo "  db-shell      Open MySQL shell inside db container"

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

restart:
	$(COMPOSE) restart consumer streamer

build-images:
	$(COMPOSE) build

db-shell:
	docker exec -it sensor-mysql mysql -uuser -ppassword sensordata



