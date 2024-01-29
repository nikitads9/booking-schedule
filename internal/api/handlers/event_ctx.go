package handlers

import (
	"context"
	"event-schedule/internal/api"
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// EventCtx middleware is used to load an Event object from
// the URL parameters passed through as the request. In case
// the Event could not be found, we stop here and return a 404.
func (i *Implementation) EventCtx(log *slog.Logger, ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "handlers.events.api.EventCtx"
			var eventID string
			var eventUUID uuid.UUID
			var event *model.EventInfo
			var err error

			log = log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(ctx)),
			)

			if eventID = chi.URLParam(r, "eventID"); eventID == "" {
				log.Error("invalid request", sl.Err(api.ErrNoEventID))
				render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
				return
			}

			eventUUID, err = uuid.FromString(eventID)
			if err != nil {
				log.Error("invalid request", sl.Err(err))
				render.Render(w, r, api.ErrInvalidRequest(err))
				return
			}

			log.Info("decoded URL param", slog.Any("eventID", eventID))

			event, err = i.Service.GetEvent(ctx, eventUUID)
			if err != nil {
				log.Error("internal error", sl.Err(err))
				render.Render(w, r, api.ErrInternalError(err))
				return
			}

			log.Info("event acquired", slog.Any("event", event))

			ctx := context.WithValue(ctx, "event", event)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
