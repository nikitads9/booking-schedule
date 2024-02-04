package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// GetVacantDates godoc
//
//	@Summary		Get vacant intervals
//	@Description	Responds with free dates within month for selected suite.
//	@ID				getDatesBySuiteID
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			suite_id path	int	true	"suite_id"	Format(int64) default(1234)
//	@Success		200	{object}	api.GetVacantDatesResponse
//	@Failure		400	{object}	api.GetVacantDatesResponse
//	@Failure		404	{object}	api.GetVacantDatesResponse
//	@Failure		422	{object}	api.GetVacantDatesResponse
//	@Failure		503	{object}	api.GetVacantDatesResponse
//	@Router			/{user_id}/{suite_id}/get-vacant-dates [get]
func (i *Implementation) GetVacantDates(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetVacantDates"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		suiteID := chi.URLParam(r, "suite_id")
		if suiteID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			return
		}

		id, err := strconv.ParseInt(suiteID, 10, 64)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		intervals, err := i.Service.GetVacantDates(ctx, id) //TODO:GetVacantDates
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("vacant dates acquired", slog.Any("intervals", intervals))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetVacantDatesAPI(intervals))
	}
}