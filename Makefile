.PHONY: deps clean

all: rocket

run: rocket
	nohup rocket > /var/log/rocket.log &

rocket:
	go install

deps:
	glide install

clean:
	rm rocket