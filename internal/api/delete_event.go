package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type DeleteEventResponse struct {
	Response *Response `json:"response"`
}

func (i *Implementation) DeleteEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		//TODO: check panic
		// Assume if we've reach this far, we can access the event
		// context because this handler is a child of the EventCtx
		// middleware. The worst case, the recoverer middleware will save us.

		//event := r.Context().Value("event").(*model.Event)

		//event, err = dbRemoveEvent(event.Uuid)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, DeleteEventResponseAPI())
	}
}

func DeleteEventResponseAPI() *DeleteEventResponse {
	return &DeleteEventResponse{
		Response: OK(),
	}
}

func (rd *DeleteEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
