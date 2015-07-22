package router

import (
	"github.com/gorilla/mux"

	"github.com/keshavdv/docklet/handlers"
)

func API() *mux.Router {
	m := mux.NewRouter()

	m.HandleFunc("/", handlers.Home)
	m.HandleFunc("/version", handlers.GetAPIVersion)

	m.HandleFunc("/pull", handlers.Pull)
	m.HandleFunc("/launch", handlers.Launch)
	return m
}
