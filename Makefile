.PHONY: default bundle fmt all

default: all

bundle:
	browserify -t reactify -e public/js/app.js -o public/js/bundle.js

fmt:
	go fmt ./...

run: bundle fmt
	reflex -g '*.go' -g '*.html' -g '*.tmpl' -g '*.js' -g '*.css' -R '^node_modules/' -s -- go run main.go
