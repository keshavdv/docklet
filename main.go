package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"

	"github.com/keshavdv/docklet/router"
)

func main() {
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(router.API())

	n.Run(":3000")
}