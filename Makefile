all: compile run

compile:
	@echo compile
	go build .

run:
	./git-trends search

docker_build:
	@echo building docker image

docker_run:
	@echo run docker container

push:
	@echo push

clean:
	@echo clean
	rm -rf *.json git-trends coverage.out

video:
	@echo make video!

gif:
	@echo make gif from video!

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out