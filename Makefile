.PHONY: build run test clean

build:
	go build -o bin/chatd ./cmd/chatd

run:
	go run ./cmd/chatd

test:
	go test ./...

clean:
	rm -rf bin/
