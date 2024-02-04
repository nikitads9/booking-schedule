package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// UpdateEvent godoc
//
//	@Summary		Updates event info
//	@Description	Updates an existing Event with given EventID and several optional fields. At least one field should not be empty. NotificationPeriod must look like {number}s,{number}m or {number}h.
//	@ID				modifyEventByJSON
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
//	@Router			/{user_id}/{event_id}/update [patch]
func (i *Implementation) UpdateEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.UpdateEvent"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		//TODO: getters for EventInfo
		event := r.Context().Value("event").(*model.EventInfo)
		if event == nil {
			log.Error("failed to load event from context", sl.Err(api.ErrEventNotFound))
			render.Render(w, r, api.ErrInternalError(api.ErrEventNotFound))
			return
		}

		req := &api.UpdateEventRequest{}
		err := render.Bind(r, req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoUserID))
		}

		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		err = i.Service.UpdateEvent(ctx, convert.ToUpdateEventInfo(req, event.EventID, id))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("event updated", slog.Any("id", event.EventID))

		render.Render(w, r, api.UpdateEventResponseAPI())
	}
}
