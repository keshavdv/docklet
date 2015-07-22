.PHONY: default bundle fmt all

default: all

deps:
	go get -d -v ./...

bundle:
	browserify -t reactify -e public/js/app.js -o public/js/bundle.js

fmt:
	go fmt ./...

run: deps bundle fmt
	reflex --decoration=fancy -c reflex.conf
