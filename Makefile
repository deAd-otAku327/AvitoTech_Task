DB_NAME=avitoshop

DB_URL=postgres://postgres:qwerty123@localhost:5432/$(DB_NAME)?sslmode=disable

ifeq ($(os), win)
#win path format
MIGRATIONS_PATH=.\internal\db\migrations
else 
#unix-like path format (default)
MIGRATIONS_PATH=./internal/db/migrations
endif

MIGRATE=docker run --rm -v $(MIGRATIONS_PATH):/migrations --network host migrate/migrate -path=/migrations/

all: build run

build:
	docker-compose build server

run:
	docker-compose up

stop:
	docker-compose stop

remove:
	docker-compose down

migrate-up:
	$(MIGRATE) -database $(DB_URL) up
migrate-down:
	$(MIGRATE) -database $(DB_URL) down -all