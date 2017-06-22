.PHONY: deps clean

all: rocket

run: clean rocket
	./rocket

rocket:
	go build

deps:
	glide install

clean:
	rm rocket