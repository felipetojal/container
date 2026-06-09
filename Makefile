build:
	@go build -o bin/container main.go

run: build
	@./bin/container