COLOUR_GREEN=\033[0;32m
COLOUR_BLUE=\033[0;34m

build:
	docker-compose build
run:
	docker-compose up -d

help:
	@echo "${COLOUR_BLUE}++ SMOOTHSTART ++${COLOUR_END}"
	@echo "${COLOUR_BLUE}Available options:${COLOUR_END}"
	@echo ""
	@echo "${COLOUR_GREEN}make build - build docker services${COLOUR_END}"
	@echo "${COLOUR_GREEN}make run   - run app${COLOUR_END}"
	
.DEFAULT_GOAL:= help
