DOCKER_IMAGE_TAG := jwt-mock:latest

lint:
	golangci-lint run ./...

test:
	go test ./... -v

.PHONY: build
build:
	go build -o jwt-mock ./build/main.go

run:
	go run build/main.go

docker/build:
	docker build . -t $(DOCKER_IMAGE_TAG)

docker/run:
	docker run -p 80:80 $(DOCKER_IMAGE_TAG)
