package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetVacantRooms godoc
//
//	@Summary		Get list of vacant rooms
//	@Description	Receives two dates as query parameters. start is to be before end and both should not be expired. Responds with list of vacant rooms and their parameters for given interval.
//	@ID				getRoomsByDates
//	@Tags			bookings
//	@Produce		json
//	@Param			start	query	string	true	"start"	Format(time.Time) default(2024-03-28T17:43:00)
//	@Param			end	query	string	true	"end"	Format(time.Time) default(2024-03-29T17:43:00)
//	@Success		200	{object}	api.GetVacantRoomsResponse
//	@Failure		400	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/get-vacant-rooms [get]
func (i *Implementation) GetVacantRooms(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.GetVacantRooms"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

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

		rooms, err := i.booking.GetVacantRooms(ctx, startDate, endDate)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("vacant rooms acquired", trace.WithAttributes(attribute.Int("quantity", len(rooms))))
		log.Info("vacant rooms acquired", slog.Int("quantity: ", len(rooms)))

		render.Status(r, http.StatusCreated)
		api.WriteWithStatus(w, http.StatusOK, api.GetVacantRoomsResponse{
			Rooms: convert.ToApiSuites(rooms),
		})
	}
}
