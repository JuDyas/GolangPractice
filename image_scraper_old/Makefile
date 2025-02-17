COMPOSE_FILE=docker-compose.yml

ifneq (,$(wildcard .env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

build:
	docker-compose -f $(COMPOSE_FILE) build

up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down

clean:
	docker-compose -f $(COMPOSE_FILE) down --volumes --remove-orphans
	docker system prune -f

logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

restart: down up

status:
	docker-compose -f $(COMPOSE_FILE) ps
