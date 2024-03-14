package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetBooking godoc
//
//	@Summary		Get booking info
//	@Description	Responds with booking info for booking with given BookingID.
//	@ID				getBookingbyTag
//	@Tags			bookings
//	@Produce		json
//
//	@Param			booking_id	path	string	true	"booking_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.GetBookingResponse
//	@Failure		400	{object}	api.GetBookingResponse
//	@Failure		401	{object}	api.GetBookingResponse
//	@Failure		404	{object}	api.GetBookingResponse
//	@Failure		422	{object}	api.GetBookingResponse
//	@Failure		503	{object}	api.GetBookingResponse
//	@Router			/{booking_id}/get [get]
//
// @Security Bearer
func (i *Implementation) GetBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.GetBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)
		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			span.RecordError(api.ErrNoUserID)
			span.SetStatus(codes.Error, api.ErrNoUserID.Error())
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			err := render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", userID)))

		bookingID := chi.URLParam(r, "booking_id")
		if bookingID == "" {
			span.RecordError(api.ErrNoBookingID)
			span.SetStatus(codes.Error, api.ErrNoBookingID.Error())
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			err := render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("bookingID extracted from path", trace.WithAttributes(attribute.String("id", bookingID)))

		bookingUUID, err := uuid.FromString(bookingID)
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

		if bookingUUID == uuid.Nil {
			span.RecordError(api.ErrNoBookingID)
			span.SetStatus(codes.Error, api.ErrNoBookingID.Error())
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("bookingID decoded")
		log.Info("decoded URL param", slog.Any("bookingID:", bookingUUID))

		booking, err := i.booking.GetBooking(ctx, bookingUUID, userID)
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

		span.AddEvent("booking acquired")
		log.Info("booking acquired", slog.Any("booking: ", booking))

		err = render.Render(w, r, api.GetBookingResponseAPI(convert.ToApiBookingInfo(booking)))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrRender(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}
	}
}
