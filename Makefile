DOCKER_IMAGE_TAG := jwt-mock:latest

lint:
	golangci-lint run ./...

build:
	go build ./

run:
	go run main.go

docker/build:
	docker build . -t $(DOCKER_IMAGE_TAG)

docker/run:
	docker run -p 80:80 $(DOCKER_IMAGE_TAG)
