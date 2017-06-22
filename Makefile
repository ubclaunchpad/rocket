.PHONY: deps clean rocket

all: rocket

run: rocket
	nohup rocket > /var/log/rocket.log &

rocket:
	go install

deps:
	glide install

clean:
	rm rocket