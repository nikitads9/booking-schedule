package api

import (
	"net/http"
)

func HandlePingCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteWithStatus(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "pong"})
	}
}
