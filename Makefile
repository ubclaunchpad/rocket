MIGRATION ?= tables

all: deps

# Install rocket as an executable
rocket:
	go install

# Install dependencies
.PHONY: deps
deps:
	dep ensure

# Clean up unwanted files
.PHONY: clean
clean:
	rm -f rocket
	docker-compose -f docker-compose.test.yml down

# Run all unit tests and report coverage
.PHONY: test
test:
	go test ./... -short -cover

# Run all integration tests and report coverage
.PHONY: test-integration
test-integration: mock-db
	go test ./... -cover
	make clean

# Sets up a local database for testing
.PHONY: mock-db
mock-db:
	docker-compose -f docker-compose.test.yml up -d

# Start the Rocket and Postgres containers
.PHONY: docker
docker:
	docker-compose up -d

# Build and start the Rocket container
.PHONY: build
build:
	docker-compose up -d --build rocket

# Run a migration specified by MIGRATION
# e.g: $ make migration MIGRATION=6_add_is_tech_lead
.PHONY: migrate
migrate:
	@docker-compose exec postgres bash -c \
		"psql -U \$$POSTGRES_USER -d \$$POSTGRES_DB -f /etc/rocket/migrations/${MIGRATION}.sql"

# Dump Rocket's DB to a file called rocket_dump.sql in the CWD
.PHONY: dump
dump:
	@docker-compose exec postgres bash -c \
		"pg_dump -f /tmp/rocket_dump.sql -U \$$POSTGRES_USER -d \$$POSTGRES_DB"
	@docker cp rocket_postgres_1:/tmp/rocket_dump.sql rocket_dump.sql
	@echo Created dump file "$(shell pwd)/rocket_dump.sql"
