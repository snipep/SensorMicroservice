.PHONY: up down restart build-images  help

COMPOSE ?= docker compose

help:
	@echo "Repo-level Makefile"
	@echo "Targets:"
	@echo "  up            Build and start all services with docker compose"
	@echo "  down          Stop and remove containers and volumes"
	@echo "  restart       Restart app services (consumer, streamer)"
	@echo "  build-images  Build images via compose"

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down -v


restart:
	$(COMPOSE) restart consumer streamer

build-images:
	$(COMPOSE) build


