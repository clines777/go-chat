.PHONY: build run test clean

build:
	go build -o bin/chatd ./cmd/chatd
	rm -rf bin/web bin/static
	cp -R web static bin/
	cp .env bin/ 2>/dev/null || true

run:
	go run ./cmd/chatd

test:
	go test ./...

clean:
	rm -rf bin/
