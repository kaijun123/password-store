APP_NAME := password_store

# ==============================================================================
# Docker support for server

build:
	@docker build -t password_store .

run:
	@docker run -p 8080:8080 password_store


# ==============================================================================
# Docker-compose support

compose-build:
	@docker compose build --no-cache

compose-up:
	@docker compose up --quiet-pull --remove-orphans

compose-down:
	@docker compose down --remove-orphans

start: compose-build compose-up