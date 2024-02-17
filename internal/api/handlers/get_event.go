package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// GetEvent godoc
//
//	@Summary		Get event info
//	@Description	Responds with event info with given EventID.
//	@ID				getEventbyTag
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			event_id	path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.GetEventResponse
//	@Failure		400	{object}	api.GetEventResponse
//	@Failure		404	{object}	api.GetEventResponse
//	@Failure		422	{object}	api.GetEventResponse
//	@Failure		503	{object}	api.GetEventResponse
//	@Router			/{user_id}/{event_id}/get [get]
func (i *Implementation) GetEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetEvent"

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

		event, err := i.Service.GetEvent(ctx, eventUUID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		err = render.Render(w, r, api.GetEventResponseAPI(event))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrRender(err))
			return
		}

		log.Info("event acquired", slog.Any("event:", event))
	}
}
