.PHONY: deps clean

all: rocket

run: clean deps rocket
	./rocket

rocket:
	go build

deps:
	glide install

clean:
	rm rocket