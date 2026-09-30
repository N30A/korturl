.PHONY: build run test templ

build: templ test
	mkdir -p ./bin
	go build -o ./bin/korturl ./cmd/korturl

run: build
	./bin/korturl

test:
	go test ./...

templ:
	templ generate
