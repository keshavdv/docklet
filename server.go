package main

import (
	"github.com/codegangsta/negroni"

	"github.com/keshavdv/docklet/router"
)

func main() {
	n := negroni.Classic()
	n.UseHandler(router.API())
	n.Run(":3000")
}