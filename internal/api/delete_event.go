package api

import (
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type DeleteEventResponse struct {
	Response *Response `json:"response"`
}

func (i *Implementation) DeleteEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.DeleteEvent"
		var err error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		//TODO: check panic
		// Assume if we've reach this far, we can access the event
		// context because this handler is a child of the EventCtx
		// middleware. The worst case, the recoverer middleware will save us.

		event := r.Context().Value("event").(*model.Event)

		//TODO: create repo method
		//event, err = dbRemoveEvent(event.Uuid)
		if err != nil {
			log.Error("failed to remove event", sl.Err(err))
			render.Render(w, r, ErrInternalError(err))
			return
		}

		log.Info("deleted event", slog.Any("id:", event.UUID))

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
