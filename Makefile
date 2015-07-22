.PHONY: default bundle fmt all

default: all

bundle:
	browserify -t reactify -e public/js/app.js -o public/js/bundle.js

fmt:
	go fmt ./...

run: bundle fmt
	reflex --decoration=fancy -c reflex.conf
