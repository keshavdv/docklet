package handlers

import (
	"net/http"

	"github.com/coreos/go-etcd/etcd"
	"github.com/unrolled/render"
)

func Launch(w http.ResponseWriter, req *http.Request) {
	r := render.New()

	machines := []string{"http://localhost:4001"}
	client := etcd.NewClient(machines)

	if _, err := client.Set("/docklet", "bar", 0); err != nil {
		r.JSON(w, http.StatusInternalServerError, map[string]string{"status": "failed"})
		returngit
	}

	r.JSON(w, http.StatusOK, map[string]string{"status": "launched"})
}
