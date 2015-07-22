package router

import (
	"github.com/gorilla/mux"

	"github.com/keshavdv/docklet/handlers"
)

func API() *mux.Router {
	m := mux.NewRouter()

	m.HandleFunc("/", handlers.Home)
	m.HandleFunc("/version", handlers.GetAPIVersion)

	m.HandleFunc("/status", handlers.Status)
	m.HandleFunc("/create", handlers.Create)
	m.HandleFunc("/inspect", handlers.Inspect)
	m.HandleFunc("/start", handlers.Start)
	m.HandleFunc("/terminal", handlers.Attach)

	m.Handle("/terminal-ws/", handlers.CreateTerminalServer())

	return m
}
