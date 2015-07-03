package main

import (
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"

	"github.com/keshavdv/docklet/router"
)

func main() {
	n := negroni.New()
	n.UseHandler(router.API())

	n.Use(negronilogrus.NewMiddleware())
	n.Run(":3000")
}