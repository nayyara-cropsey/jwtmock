lint:
	golangci-lint run ./...

build:
	go build ./

run:
	go run main.go
