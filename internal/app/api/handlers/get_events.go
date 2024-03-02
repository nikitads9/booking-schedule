package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/logger/sl"
	"log/slog"
	"strconv"
	"time"

	"net/http"

	"github.com/go-chi/chi"
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
//	@Param			start query		string	true	"start" Format(time.Time) default(2024-03-28T17:43:00Z)
//	@Param			end query		string	true	"end" Format(time.Time) default(2024-03-29T17:43:00Z)
//	@Success		200	{object}	api.GetEventsResponse
//	@Failure		400	{object}	api.GetEventsResponse
//	@Failure		404	{object}	api.GetEventsResponse
//	@Failure		422	{object}	api.GetEventsResponse
//	@Failure		503	{object}	api.GetEventsResponse
//	@Router			/{user_id}/get-events [get]
func (i *Implementation) GetEvents(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetEvents"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoUserID))
			return
		}

		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}

		start := r.URL.Query().Get("start")
		if start == "" {
			log.Error("invalid request", sl.Err(api.ErrNoInterval))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoInterval))
			return
		}

		end := r.URL.Query().Get("end")
		if end == "" {
			log.Error("invalid request", sl.Err(api.ErrNoInterval))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoInterval))
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04:05-07:00", start)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}
		endDate, err := time.Parse("2006-01-02T15:04:05-07:00", end)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}

		err = api.CheckDates(startDate, endDate)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		log.Info("received request", slog.Any("params:", start+" to "+end))

		events, err := i.Service.GetEvents(ctx, startDate, endDate, id)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("events acquired", slog.Int("quantity:", len(events)))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetEventsResponseAPI(convert.ToApiEventsInfo(events)))
	}

}
