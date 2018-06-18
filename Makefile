all: deps

rocket:
	go install

.PHONY: deps
deps:
	dep ensure

.PHONY: clean
clean:
	rm -f rocket
	docker-compose -f docker-compose.test.yml down

.PHONY: test
test:
	go test ./... -short -cover

.PHONY: test-integration
test-integration: mock-db
	go test ./... -cover
	make clean

# Sets up a local database for testing
.PHONY: mock-db
mock-db:
	docker-compose -f docker-compose.test.yml up -d

.PHONY: docker
docker:
	docker-compose up -d

.PHONY: build
build:
	docker-compose up -d --build rocket
