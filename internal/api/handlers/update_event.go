package handlers

import (
	"context"
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// UpdateEvent godoc
//
//	@Summary		Updates event info
//	@Description	Updates an existing Event with given EventID and several optional fields. At least one field should not be empty.
//	NotificationPeriod must look like {number}s,{number}m or {number}h.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			event_id path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Param          event body		api.UpdateEventRequest	true	"UpdateEventRequest"
//	@Success		200	{object}	api.UpdateEventResponse
//	@Failure		400	{object}	api.UpdateEventResponse
//	@Failure		404	{object}	api.UpdateEventResponse
//	@Failure		422	{object}	api.UpdateEventResponse
//	@Failure		503	{object}	api.UpdateEventResponse
//	@Router			/events/{user_id}/{event-id}/update [patch]
func (i *Implementation) UpdateEvent(log *slog.Logger, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.UpdateEvent"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		//TODO: getters for EventInfo
		event := r.Context().Value("event").(*model.EventInfo)

		req := &api.UpdateEventRequest{}
		if err := render.Bind(r, req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		err := i.Service.UpdateEvent(ctx, convert.ToUpdateEventInfo(req, event.EventID))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}
		//i.Service.UpdateEvent(event.Uuid, data.EventInfo) резервный вариант
		//обработка ошибки

		log.Info("event updated", slog.Any("id", event.EventID))

		render.Render(w, r, api.UpdateEventResponseAPI())
	}
}
