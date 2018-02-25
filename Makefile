.PHONY: deps clean rocket test

all: rocket

rocket:
	go install

deps:
	glide install

clean:
	rm rocket

test:
	go test ./... -cover

docker:
	docker-compose up -d