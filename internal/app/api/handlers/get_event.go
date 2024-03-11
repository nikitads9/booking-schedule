package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/logger/sl"
	"event-schedule/internal/middleware/auth"
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
//	@Description	Responds with booking info for booking with given EventID.
//	@ID				getEventbyTag
//	@Tags			bookings
//	@Produce		json
//
//	@Param			event_id	path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.GetEventResponse
//	@Failure		400	{object}	api.GetEventResponse
//	@Failure		401	{object}	api.GetEventResponse
//	@Failure		404	{object}	api.GetEventResponse
//	@Failure		422	{object}	api.GetEventResponse
//	@Failure		503	{object}	api.GetEventResponse
//	@Router			/{event_id}/get [get]
//
// @Security Bearer
func (i *Implementation) GetEvent(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetEvent"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			return
		}

		eventID := chi.URLParam(r, "event_id")
		if eventID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoEventID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
			return
		}

		eventUUID, err := uuid.FromString(eventID)
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

		event, err := i.Booking.GetEvent(ctx, eventUUID, userID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		err = render.Render(w, r, api.GetEventResponseAPI(convert.ToApiEventInfo(event)))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrRender(err))
			return
		}

		log.Info("event acquired", slog.Any("event:", event))
	}
}
