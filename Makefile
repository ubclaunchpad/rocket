.PHONY: deps clean

all: rocket

run: rocket
	./rocket

rocket:
	go build

deps:
	glide install

clean:
	rm rocket