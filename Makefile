all: rocket

.PHONY: rocket
rocket:
	go install

.PHONY: deps
deps:
	glide install

.PHONY: clean
clean:
	rm rocket
	pg_ctl -D /usr/local/var/postgres stop -s -m fast

.PHONY: test
test:
	go test ./... -cover

# Sets up a local database for testing
.PHONY: mock-db
mock-db:
	sh mock_db.sh

.PHONY: docker
docker:
	docker-compose up -d

.PHONY: build
build:
	docker-compose up -d --build rocket
