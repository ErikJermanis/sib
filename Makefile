run: build
	@./bin/main.go

build:
	@go build -o bin/main.go