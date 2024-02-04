package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// GetEvents godoc
//
//	@Summary		Get several events info
//	@Description	Responds with series of event info objects within given time period. The parameters are start date and end date.
//	@ID				getMultipleEventsByTag
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			interval	path	string	true	"interval" default(all)
//	@Success		200	{object}	api.GetEventsResponse
//	@Failure		400	{object}	api.GetEventsResponse
//	@Failure		404	{object}	api.GetEventsResponse
//	@Failure		422	{object}	api.GetEventsResponse
//	@Failure		503	{object}	api.GetEventsResponse
//	@Router			/{user_id}/?start={start}&end={end} [get]
func (i *Implementation) GetEvents(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetEvents"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		getEventsInfo, err := convert.ToGetEventsInfo(r)
		if err != nil {
			log.Error("invalid request", sl.Err(api.ErrMissingValues))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		log.Info("received request", slog.Any("params:", getEventsInfo))

		events, err := i.Service.GetEvents(ctx, getEventsInfo)
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
