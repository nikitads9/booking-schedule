package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"errors"
	"log/slog"

	"net/http"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// AddBooking godoc
//
//	@Summary		Adds booking
//	@Description	Adds an  associated with user with given parameters. NotificationPeriod is optional and must look like {number}s,{number}m or {number}h. Implemented with the use of transaction: first rooms availibility is checked. In case one's new booking request intersects with and old one(even if belongs to him), the request is considered erratic. startDate is to be before endDate and both should not be expired.
//	@ID				addByBookingJSON
//	@Tags			bookings
//	@Accept			json
//	@Produce		json
//
//	@Param          booking	body	api.AddBookingRequest	true	"BookingEntry"
//	@Success		200	{object}	api.AddBookingResponse
//	@Failure		400	{object}	api.AddBookingResponse
//	@Failure		401	{object}	api.AddBookingResponse
//	@Failure		404	{object}	api.AddBookingResponse
//	@Failure		422	{object}	api.AddBookingResponse
//	@Failure		503	{object}	api.AddBookingResponse
//	@Router			/add [post]
//
// @Security Bearer
func (i *Implementation) AddBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "bookings.api.handlers.AddBooking"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		req := &api.AddBookingRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				// Приводим ошибку к типу ошибки валидации
				validateErr := err.(validator.ValidationErrors)
				log.Error("some of the required values were not received", sl.Err(validateErr))
				err = render.Render(w, r, api.ErrValidationError(validateErr))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			err = render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}
		//TODO: getters
		mod, err := convert.ToBookingInfo(&api.Booking{
			UserID:    userID,
			SuiteID:   req.SuiteID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			NotifyAt:  req.NotifyAt,
		})

		if err != nil {
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		bookingID, err := i.Booking.AddBooking(ctx, mod)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		log.Info("booking added", slog.Any("id: ", bookingID))

		render.Status(r, http.StatusCreated)
		err = render.Render(w, r, api.AddBookingResponseAPI(bookingID))
		if err != nil {
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}

}
