package handlers

import (
	"context"
	"event-schedule/internal/api"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"strconv"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// GetEvents godoc
//
//	@Summary		Get several events info
//	@Description	Responds with series of event info objects within given time period.
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			event_id	path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Param			interval	path	string	true	"interval" default(all)
//	@Success		200	{object}	api.GetEventsResponse
//	@Failure		400	{object}	api.GetEventsResponse
//	@Failure		404	{object}	api.GetEventsResponse
//	@Failure		422	{object}	api.GetEventsResponse
//	@Failure		503	{object}	api.GetEventsResponse
//	@Router			/events/{user_id}/{interval} [get]
func (i *Implementation) GetEvents(log *slog.Logger, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetEvents"
		var userID, period string
		var id int64
		var err error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		if userID = chi.URLParam(r, "userID"); userID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoUserID))
			return
		}

		if id, err = strconv.ParseInt(userID, 10, 64); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		if period = chi.URLParam(r, "period"); period == "" {
			log.Error("invalid request", sl.Err(api.ErrNoInterval))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoInterval))
			return
		}

		events, err := i.Service.GetEvents(ctx, id, period) //TODO: filter
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("events acquired", slog.Int("quantity", len(events)))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetEventsResponseAPI(events))
	}

}
