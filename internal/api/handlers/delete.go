package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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
		const op = "handlers.events.api.DeleteEvent"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)
		//TODO: check panic
		// Assume if we've reach this far, we can access the event
		// context because this handler is a child of the EventCtx
		// middleware. The worst case, the recoverer middleware will save us.

		event := r.Context().Value("event").(*model.EventInfo)
		if event == nil {
			log.Error("failed to load event from context", sl.Err(api.ErrEventNotFound))
			render.Render(w, r, api.ErrInternalError(api.ErrEventNotFound))
			return
		}

		err := i.Service.DeleteEvent(ctx, event.EventID)
		if err != nil {
			log.Error("failed to remove event", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("deleted event", slog.Any("id:", event.EventID))

		render.Render(w, r, api.DeleteEventResponseAPI())
	}
}
