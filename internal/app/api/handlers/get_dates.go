package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/logger/sl"
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
//	@Description	Responds with list of vacant intervals within month for selected suite.
//	@ID				getDatesBySuiteID
//	@Tags			events
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1)
//	@Param			suite_id path	int	true	"suite_id"	Format(int64) default(1)
//	@Success		200	{object}	api.GetVacantDatesResponse
//	@Failure		400	{object}	api.GetVacantDatesResponse
//	@Failure		404	{object}	api.GetVacantDatesResponse
//	@Failure		422	{object}	api.GetVacantDatesResponse
//	@Failure		503	{object}	api.GetVacantDatesResponse
//	@Router			/{user_id}/{suite_id}/get-vacant-dates [get]
func (i *Implementation) GetVacantDates(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetVacantDates"

		ctx := r.Context()

		log := logger.With(
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
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}

		if id == 0 {
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			return
		}

		intervals, err := i.Service.GetVacantDates(ctx, id)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("vacant dates acquired", slog.Any("quantity:", len(intervals)))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetVacantDatesAPI(convert.ToFreeIntervals(intervals)))
		//TODO: ошибка и проверка на ошибку
	}
}
