package router

import (
	"github.com/gorilla/mux"

	"github.com/keshavdv/docklet/handlers"
)

func API() *mux.Router {
	m := mux.NewRouter()

	m.HandleFunc("/version", handlers.GetAPIVersion)
	return m
}