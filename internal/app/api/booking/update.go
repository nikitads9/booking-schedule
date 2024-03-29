package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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
//	@Success		200
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/{booking_id}/update [patch]
//
// @Security Bearer
func (i *Implementation) UpdateBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.UpdateBooking"

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

		req := &api.UpdateBookingRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, err.Error())
				log.Error("some of the required values were not received", sl.Err(validateErr))
				api.WriteValidationError(w, validateErr)
				return
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to decode request body", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		span.AddEvent("request body decoded")
		log.Info("request body decoded", slog.Any("req", req))

		bookingID := chi.URLParam(r, "booking_id")
		if bookingID == "" {
			span.RecordError(errNoBookingID)
			span.SetStatus(codes.Error, errNoBookingID.Error())
			log.Error("invalid request", sl.Err(errNoBookingID))
			api.WriteWithError(w, http.StatusBadRequest, errNoBookingID.Error())
			return
		}

		span.AddEvent("bookingID extracted from path", trace.WithAttributes(attribute.String("id", bookingID)))

		bookingUUID, err := uuid.FromString(bookingID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, api.ErrParse.Error())
			return
		}

		if bookingUUID == uuid.Nil {
			span.RecordError(errNoBookingID)
			span.SetStatus(codes.Error, errNoBookingID.Error())
			log.Error("invalid request", sl.Err(errNoBookingID))
			api.WriteWithError(w, http.StatusBadRequest, errNoBookingID.Error())
			return
		}

		span.AddEvent("booking uuid decoded")
		//TODO: getters
		mod, err := convert.ToBookingInfo(&api.Booking{
			BookingID: bookingUUID,
			UserID:    userID,
			SuiteID:   req.SuiteID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			NotifyAt:  req.NotifyAt,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		span.AddEvent("converted to booking model")

		err = i.booking.UpdateBooking(ctx, mod)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("booking updated")
		log.Info("booking updated", slog.Any("id: ", mod.ID))

		api.WriteWithStatus(w, http.StatusOK, nil)
	}
}
