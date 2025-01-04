APP_NAME := samsung-tv-remote
VERSION := 1.0.0

.PHONY: build install clean test

build:
	go build -o bin/$(APP_NAME) ./cmd/samsung-tv-remote

install:
	go install ./cmd/samsung-tv-remote

clean:
	rm -rf bin

test:
	go test ./...
