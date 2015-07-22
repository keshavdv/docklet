package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"

	"github.com/keshavdv/docklet/context"
	"github.com/keshavdv/docklet/handlers"
	"github.com/keshavdv/docklet/router"
	"github.com/unrolled/render"
)

func main() {
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewStatic(http.Dir("public")))
	render := render.New(render.Options{})
	renderer := context.NewRender(render)
	n.Use(renderer)
	n.UseHandler(router.API())

	go handlers.Hub.Run()
	n.Run(":3000")
}
