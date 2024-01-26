package api

import (
	"context"
	"errors"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// EventCtx middleware is used to load an Event object from
// the URL parameters passed through as the request. In case
// the Event could not be found, we stop here and return a 404.
func (i *Implementation) EventCtx(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var eventID string
			var event *model.Event
			var err error

			if eventID = chi.URLParam(r, "eventID"); eventID != "" {
				//tofo getEvent
				//event, err = GetEvent(eventID)
			} else {
				render.Render(w, r, ErrInvalidRequest(errors.New("received no eventID")))
				return
			}
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), "event", event)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
