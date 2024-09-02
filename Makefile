COLOUR_GREEN=\033[0;32m
COLOUR_BLUE=\033[0;34m

build:
	docker-compose build
run:
	docker-compose -p smoothstart-app up -d
run-prod:
	docker compose --file docker-compose.yml --env-file .env up --build -d
run-dev:
	docker compose --file docker-compose.dev.yml --env-file .env.dev up --build -d
	@echo ""
	@echo "Redis and Postgresql are running, now run commands below in separate terminal windows:"
	@echo "air"
	@echo "npx tailwindcss -i ./src/styles.css -o ./assets/css/styles.css --watch"
help:
	@echo "${COLOUR_BLUE}++ SMOOTHSTART ++${COLOUR_END}"
	@echo "${COLOUR_BLUE}Available options:${COLOUR_END}"
	@echo ""
	@echo "${COLOUR_GREEN}make run-dev 	  - run app in development environment${COLOUR_END}"
	@echo "${COLOUR_GREEN}make run-prod  	  - run app in production environment${COLOUR_END}"
	
.DEFAULT_GOAL:= help
