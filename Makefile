run: templ build
	@./bin/main.go

all: templ css build
	@./bin/main.go

build:
	@go build -o bin/main.go

css:
	@npx tailwindcss -i ./views/input.css -o ./public/styles.css

templ:
	@templ generate