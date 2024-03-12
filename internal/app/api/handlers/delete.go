package handlers

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

// DeleteBooking godoc
//
//	@Summary		Deletes an booking
//	@Description	Deletes an booking with given UUID.
//	@ID				removeByBookingID
//	@Tags			bookings
//	@Produce		json
//
//	@Param			booking_id path	string	true	"booking_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	api.DeleteBookingResponse
//	@Failure		400	{object}	api.DeleteBookingResponse
//	@Failure		401	{object}	api.DeleteBookingResponse
//	@Failure		404	{object}	api.DeleteBookingResponse
//	@Failure		422	{object}	api.DeleteBookingResponse
//	@Failure		503	{object}	api.DeleteBookingResponse
//	@Router			/{booking_id}/delete [delete]
//
// @Security Bearer
func (i *Implementation) DeleteBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "bookings.api.handlers.DeleteBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		userID := auth.UserIDFromContext(ctx)

		if userID == 0 {
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			err := render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		bookingID := chi.URLParam(r, "booking_id")
		if bookingID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			err := render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		bookingUUID, err := uuid.FromString(bookingID)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		if bookingUUID == uuid.Nil {
			log.Error("invalid request", sl.Err(api.ErrNoBookingID))
			err = render.Render(w, r, api.ErrInvalidRequest(api.ErrNoBookingID))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		log.Info("decoded URL param", slog.Any("bookingID:", bookingUUID))

		err = i.Booking.DeleteBooking(ctx, bookingUUID, userID)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		log.Info("deleted booking", slog.Any("id: ", bookingUUID))

		err = render.Render(w, r, api.DeleteBookingResponseAPI())
		if err != nil {
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}
}
