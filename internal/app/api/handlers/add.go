package handlers

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
//	@Param          booking	body	api.AddBookingRequest	true	"AddBookingRequest"
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
				render.Render(w, r, api.ErrValidationError(validateErr))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		userID := auth.UserIDFromContext(ctx)
		//id, err := strconv.ParseInt(userID, 10, 64)
		if userID == 0 {
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
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
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		bookingID, err := i.Booking.AddBooking(ctx, mod)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("booking added", slog.Any("id:", bookingID))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.AddBookingResponseAPI(bookingID))
	}

}

//TODO:
//для презентации времени в нормальном виде
/* type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date

	return
}

func (t *CustomTime) ExcelDate() string {
    return t.Format("01/02/2006")
}
*/
