package handlers

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
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
//	@Tags			bookings
//	@Produce		json
//	@Param			suite_id path	int	true	"suite_id"	Format(int64) default(1)
//	@Success		200	{object}	api.GetVacantDatesResponse
//	@Failure		400	{object}	api.GetVacantDatesResponse
//	@Failure		404	{object}	api.GetVacantDatesResponse
//	@Failure		422	{object}	api.GetVacantDatesResponse
//	@Failure		503	{object}	api.GetVacantDatesResponse
//	@Router			/{suite_id}/get-vacant-dates [get]
func (i *Implementation) GetVacantDates(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "bookings.api.handlers.GetVacantDates"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		suiteID := chi.URLParam(r, "suite_id")
		if suiteID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			err := render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		id, err := strconv.ParseInt(suiteID, 10, 64)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		if id == 0 {
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		dates, err := i.Booking.GetBusyDates(ctx, id)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		log.Info("vacant dates acquired", slog.Int("quantity: ", len(dates)))

		render.Status(r, http.StatusCreated)
		err = render.Render(w, r, api.GetVacantDatesAPI(convert.ToVacantIntervals(dates)))
		if err != nil {
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}
}
