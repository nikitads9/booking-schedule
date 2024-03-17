package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func HandlePingCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = render.Render(w, r, PingResponseAPI())
	}
}
