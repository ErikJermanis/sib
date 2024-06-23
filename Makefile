run: templ css build
	@./bin/sibweb

deploy: templ css buildlinux
	scp ./bin/sibweblinux erik@erikjermanis.me:/home/erik/sib-web
	scp ./public/styles.css erik@erikjermanis.me:/home/erik/sib-web/public
	scp ./public/htmx.min.js erik@erikjermanis.me:/home/erik/sib-web/public
	scp ./public/favicon.png erik@erikjermanis.me:/home/erik/sib-web/public

buildlinux:
	@GOOS=linux GOARCH=amd64 go build -o bin/sibweblinux

build:
	@go build -o bin/sibweb

css:
	@npx tailwindcss -i ./views/input.css -o ./public/styles.css

templ:
	@templ generate