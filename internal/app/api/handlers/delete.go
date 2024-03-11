package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/logger/sl"
	"event-schedule/internal/middleware/auth"
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
//	@Tags			bookings
//	@Produce		json
//
//	@Param			event_id path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.DeleteEventResponse
//	@Failure		400	{object}	api.DeleteEventResponse
//	@Failure		401	{object}	api.DeleteEventResponse
//	@Failure		404	{object}	api.DeleteEventResponse
//	@Failure		422	{object}	api.DeleteEventResponse
//	@Failure		503	{object}	api.DeleteEventResponse
//	@Router			/{event_id}/delete [delete]
//
// @Security Bearer
func (i *Implementation) DeleteEvent(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.DeleteEvent"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		userID := auth.UserIDFromContext(ctx)

		//id, err := strconv.ParseInt(userID, 10, 64)
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

		err = i.Booking.DeleteEvent(ctx, eventUUID, userID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("deleted event", slog.Any("id:", eventUUID))
		render.Render(w, r, api.DeleteEventResponseAPI())
	}
}
