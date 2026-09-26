.PHONY: build run test

build: test
	mkdir -p ./bin
	go build -o ./bin/korturl ./cmd/korturl

run: build
	./bin/korturl

test:
	go test ./...
