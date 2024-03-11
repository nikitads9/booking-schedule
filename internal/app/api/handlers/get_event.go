package handlers

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
		const op = "bookings.api.handlers.GetBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
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
			return
		}

		if bookingUUID == uuid.Nil {
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			return
		}

		log.Info("decoded URL param", slog.Any("bookingID:", bookingUUID))

		booking, err := i.Booking.GetBooking(ctx, bookingUUID, userID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		err = render.Render(w, r, api.GetBookingResponseAPI(convert.ToApiBookingInfo(booking)))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrRender(err))
			return
		}

		log.Info("booking acquired", slog.Any("booking:", booking))
	}
}
