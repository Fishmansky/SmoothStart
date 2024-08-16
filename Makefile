COLOUR_GREEN=\033[0;32m
COLOUR_BLUE=\033[0;34m

build:
	docker-compose build
run:
	docker-compose -p smoothstart-app up -d
switch-to-dev:
	mv .env env_prod
	mv .env.dev .env
	mv docker-compose.yml docker-compose.yml_prod
	mv docker-compose.dev.yml docker-compose.yml
switch-to-prod:
	mv .env .env.dev
	mv env_prod .env
	mv docker-compose.yml docker-compose.dev.yml
	mv docker-compose.yml_prod docker-compose.yml

help:
	@echo "${COLOUR_BLUE}++ SMOOTHSTART ++${COLOUR_END}"
	@echo "${COLOUR_BLUE}Available options:${COLOUR_END}"
	@echo ""
	@echo "${COLOUR_GREEN}make build 	  - build docker services${COLOUR_END}"
	@echo "${COLOUR_GREEN}make run   	  - run app${COLOUR_END}"
	@echo "${COLOUR_GREEN}make switch-to-prod - change env and docker-compose files to prod${COLOUR_END}"
	@echo "${COLOUR_GREEN}make switch-to-dev  - change env and docker-compose files to dev${COLOUR_END}"
	
.DEFAULT_GOAL:= help
