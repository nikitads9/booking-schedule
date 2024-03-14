package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateBooking godoc
//
//	@Summary		Updates booking info
//	@Description	Updates an existing booking with given BookingID, suiteID, startDate, endDate values (notificationPeriod being optional). Implemented with the use of transaction: first room availibility is checked. In case one attempts to alter his previous booking (i.e. widen or tighten its' limits) the booking is updated.  If it overlaps with smb else's booking or with clients' another booking the request is considered unsuccessful. startDate parameter  is to be before endDate and both should not be expired.
//	@ID				modifyBookingByJSON
//	@Tags			bookings
//	@Accept			json
//	@Produce		json
//
//	@Param			booking_id path	string	true	"booking_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Param          booking body		api.UpdateBookingRequest	true	"BookingEntry"
//	@Success		200	{object}	api.UpdateBookingResponse
//	@Failure		400	{object}	api.UpdateBookingResponse
//	@Failure		401	{object}	api.UpdateBookingResponse
//	@Failure		404	{object}	api.UpdateBookingResponse
//	@Failure		422	{object}	api.UpdateBookingResponse
//	@Failure		503	{object}	api.UpdateBookingResponse
//	@Router			/{booking_id}/update [patch]
//
// @Security Bearer
func (i *Implementation) UpdateBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.UpdateBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)
		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		id := auth.UserIDFromContext(ctx)
		if id == 0 {
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

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", id)))

		req := &api.UpdateBookingRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, err.Error())
				log.Error("some of the required values were not received", sl.Err(validateErr))
				err = render.Render(w, r, api.ErrValidationError(validateErr))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to decode request body", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("request body decoded")
		log.Info("request body decoded", slog.Any("req", req))

		bookingID := chi.URLParam(r, "booking_id")
		if bookingID == "" {
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

		span.AddEvent("booking uuid decoded")
		//TODO: getters
		mod, err := convert.ToBookingInfo(&api.Booking{
			BookingID: bookingUUID,
			UserID:    id,
			SuiteID:   req.SuiteID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			NotifyAt:  req.NotifyAt,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
		}

		span.AddEvent("converted to booking model")

		err = i.booking.UpdateBooking(ctx, mod)
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

		span.AddEvent("booking updated")
		log.Info("booking updated", slog.Any("id: ", mod.ID))

		err = render.Render(w, r, api.UpdateBookingResponseAPI())
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}
}
