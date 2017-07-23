.PHONY: deps clean rocket

all: rocket

rocket:
	go install

deps:
	glide install

clean:
	rm rocket