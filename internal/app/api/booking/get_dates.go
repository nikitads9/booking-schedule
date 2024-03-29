package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/logger/sl"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
//	@Failure		400	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/{suite_id}/get-vacant-dates [get]
func (i *Implementation) GetVacantDates(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.GetVacantDates"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

		suiteID := chi.URLParam(r, "suite_id")
		if suiteID == "" {
			span.RecordError(errNoSuiteID)
			span.SetStatus(codes.Error, errNoSuiteID.Error())
			log.Error("invalid request", sl.Err(errNoSuiteID))
			api.WriteWithError(w, http.StatusBadRequest, errNoSuiteID.Error())
			return
		}

		span.AddEvent("suiteID extracted from path", trace.WithAttributes(attribute.String("id", suiteID)))

		id, err := strconv.ParseInt(suiteID, 10, 64)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, api.ErrParse.Error())
			return
		}

		if id == 0 {
			span.RecordError(errNoSuiteID)
			span.SetStatus(codes.Error, errNoSuiteID.Error())
			log.Error("invalid request", sl.Err(errNoSuiteID))
			api.WriteWithError(w, http.StatusBadRequest, errNoSuiteID.Error())
			return
		}

		span.AddEvent("suiteID parsed")

		dates, err := i.booking.GetVacantDates(ctx, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("vacant dates acquired", trace.WithAttributes(attribute.Int("quantity", len(dates))))
		log.Info("vacant dates acquired", slog.Int("quantity: ", len(dates)))

		api.WriteWithStatus(w, http.StatusOK, api.GetVacantDatesResponse{
			Intervals: dates,
		})
	}
}
