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
//	@Description	Responds with series of event info objects within given time period. The query parameters are start date and end date (start is to be before end and both should not be expired).
//	@ID				getMultipleEventsByTag
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			start query		string	true	"start" Format(time.Time) default(2024-03-28T17:43:00-03:00)
//	@Param			end query		string	true	"end" Format(time.Time) default(2024-03-29T17:43:00-03:00)
//	@Success		200	{object}	api.GetEventsResponse
//	@Failure		400	{object}	api.GetEventsResponse
//	@Failure		404	{object}	api.GetEventsResponse
//	@Failure		422	{object}	api.GetEventsResponse
//	@Failure		503	{object}	api.GetEventsResponse
//	@Router			/{user_id}/get-events [get]
func (i *Implementation) GetEvents(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetEvents"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		getEventsInfo, err := convert.ToGetEventsInfo(r)
		if err != nil {
			log.Error("invalid request", sl.Err(err)) //TODO: log real errors
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

		log.Info("events acquired", slog.Int("quantity:", len(events)))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetEventsResponseAPI(events))
	}

}
