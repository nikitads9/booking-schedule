package booking

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
//	@Failure		400	{object}	api.GetVacantDatesResponse
//	@Failure		404	{object}	api.GetVacantDatesResponse
//	@Failure		422	{object}	api.GetVacantDatesResponse
//	@Failure		503	{object}	api.GetVacantDatesResponse
//	@Router			/{suite_id}/get-vacant-dates [get]
func (i *Implementation) GetVacantDates(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.GetVacantDates"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)
		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		suiteID := chi.URLParam(r, "suite_id")
		if suiteID == "" {
			span.RecordError(api.ErrNoSuiteID)
			span.SetStatus(codes.Error, api.ErrNoSuiteID.Error())
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			err := render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("suiteID extracted from path", trace.WithAttributes(attribute.String("id", suiteID)))

		id, err := strconv.ParseInt(suiteID, 10, 64)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		if id == 0 {
			span.RecordError(api.ErrNoSuiteID)
			span.SetStatus(codes.Error, api.ErrNoSuiteID.Error())
			log.Error("invalid request", sl.Err(api.ErrNoSuiteID))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrNoSuiteID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("suiteID parsed")

		dates, err := i.booking.GetBusyDates(ctx, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("busy dates acquired", trace.WithAttributes(attribute.Int("quantity", len(dates))))
		log.Info("busy dates acquired", slog.Int("quantity: ", len(dates)))

		vacant := convert.ToVacantIntervals(dates)

		span.AddEvent("converted to vacant dates", trace.WithAttributes(attribute.Int("quantity", len(vacant))))
		log.Info("converted to vacant dates", slog.Int("quantity: ", len(vacant)))

		err = render.Render(w, r, api.GetVacantDatesAPI(vacant))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}
}
