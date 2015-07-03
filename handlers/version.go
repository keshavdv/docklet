package handlers

import (
	"net/http"

	"github.com/unrolled/render"
)

func GetAPIVersion(w http.ResponseWriter, req *http.Request) {
	r := render.New()
	r.JSON(w, http.StatusOK, map[string]string{"version": "0.0.1"})
}
