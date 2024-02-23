package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// DeleteEvent godoc
//
//	@Summary		Deletes an event
//	@Description	Deletes an event with given UUID.
//	@ID				removeByEventID
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			event_id path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.DeleteEventResponse
//	@Failure		400	{object}	api.DeleteEventResponse
//	@Failure		404	{object}	api.DeleteEventResponse
//	@Failure		422	{object}	api.DeleteEventResponse
//	@Failure		503	{object}	api.DeleteEventResponse
//	@Router			/{user_id}/{event_id}/delete [delete]
func (i *Implementation) DeleteEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.DeleteEvent"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		id := chi.URLParam(r, "event_id")
		if id == "" {
			log.Error("invalid request", sl.Err(api.ErrNoEventID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
			return
		}

		eventUUID, err := uuid.FromString(id)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}

		if eventUUID == uuid.Nil {
			log.Error("invalid request", sl.Err(api.ErrNoEventID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
			return
		}

		log.Info("decoded URL param", slog.Any("eventID:", eventUUID))

		err = i.Service.DeleteEvent(ctx, eventUUID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("deleted event", slog.Any("id:", eventUUID))

		render.Render(w, r, api.DeleteEventResponseAPI())
	}
}
