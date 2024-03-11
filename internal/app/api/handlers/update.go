package handlers

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
//	@Param          booking body		api.UpdateBookingRequest	true	"UpdateBookingRequest"
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
		const op = "bookings.api.handlers.UpdateBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		req := &api.UpdateBookingRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				log.Error("some of the required values were not received", sl.Err(validateErr))
				render.Render(w, r, api.ErrValidationError(validateErr))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		id := auth.UserIDFromContext(ctx)
		//id, err := strconv.ParseInt(userID, 10, 64)
		if id == 0 {
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			return
		}

		bookingID := chi.URLParam(r, "booking_id")
		if bookingID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))

			return
		}

		bookingUUID, err := uuid.FromString(bookingID)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
		}

		if bookingUUID == uuid.Nil {
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			return
		}

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
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		err = i.Booking.UpdateBooking(ctx, mod)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("booking updated", slog.Any("id:", mod.ID))

		render.Render(w, r, api.UpdateBookingResponseAPI())
	}
}
