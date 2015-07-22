package context

import (
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"net/http"
)

type RenderMiddleware struct {
	Render *render.Render
}

func NewRender(render *render.Render) *RenderMiddleware {
	return &RenderMiddleware{render}
}

func (render RenderMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Render", render.Render)
	next(rw, r)
}
