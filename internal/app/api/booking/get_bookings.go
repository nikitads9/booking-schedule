package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"log/slog"
	"time"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetBookings godoc
//
//	@Summary		Get several bookings info
//	@Description	Responds with series of booking info objects within given time period. The query parameters are start date and end date (start is to be before end and both should not be expired).
//	@ID				getMultipleBookingsByTag
//	@Tags			bookings
//	@Produce		json
//
//	@Param			start query		string	true	"start" Format(time.Time) default(2024-03-28T17:43:00)
//	@Param			end query		string	true	"end" Format(time.Time) default(2024-03-29T17:43:00)
//	@Success		200	{object}	api.GetBookingsResponse
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/get-bookings [get]
//
// @Security Bearer
func (i *Implementation) GetBookings(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.GetBookings"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			span.RecordError(api.ErrNoUserID)
			span.SetStatus(codes.Error, api.ErrNoUserID.Error())
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			api.WriteWithError(w, http.StatusUnauthorized, api.ErrNoAuth.Error())
			return
		}

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", userID)))

		start := r.URL.Query().Get("start")
		if start == "" {
			span.RecordError(errNoInterval)
			span.SetStatus(codes.Error, errNoInterval.Error())
			log.Error("invalid request", sl.Err(errNoInterval))
			api.WriteWithError(w, http.StatusBadRequest, errNoInterval.Error())
			return
		}

		span.AddEvent("startDate extracted from query", trace.WithAttributes(attribute.String("start", start)))

		end := r.URL.Query().Get("end")
		if end == "" {
			span.RecordError(errNoInterval)
			span.SetStatus(codes.Error, errNoInterval.Error())
			log.Error("invalid request", sl.Err(errNoInterval))
			api.WriteWithError(w, http.StatusBadRequest, errNoInterval.Error())
			return
		}

		span.AddEvent("endDate extracted from query", trace.WithAttributes(attribute.String("end", end)))

		startDate, err := time.Parse("2006-01-02T15:04:05", start)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, api.ErrParse.Error())
			return
		}
		endDate, err := time.Parse("2006-01-02T15:04:05", end)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, api.ErrParse.Error())
			return
		}

		span.AddEvent("start and end dates parsed")

		err = api.CheckDates(startDate, endDate)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		span.AddEvent("dates verified")
		log.Info("received request", slog.Any("params:", start+" to "+end))

		bookings, err := i.booking.GetBookings(ctx, startDate, endDate, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("bookings acquired", trace.WithAttributes(attribute.Int("quantity", len(bookings))))
		log.Info("bookings acquired", slog.Int("quantity: ", len(bookings)))

		render.Status(r, http.StatusCreated)
		api.WriteWithStatus(w, http.StatusOK, api.GetBookingsResponse{
			BookingsInfo: convert.ToApiBookingsInfo(bookings),
		})

	}

}
