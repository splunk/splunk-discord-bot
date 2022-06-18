.DEFAULT_GOAL := build
tools:
	go install github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen

generate:
	go generate ./...

build:
	go build -o bin/bot ./cmd/main.go

run:
	go run ./cmd/main.go