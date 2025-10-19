.PHONY: build run

build:
	go build -o bin/codegraph .

run: build
	./bin/codegraph parse --output graph.json .
