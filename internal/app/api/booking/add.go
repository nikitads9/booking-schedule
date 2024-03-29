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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/add [post]
//
// @Security Bearer
func (i *Implementation) AddBooking(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.booking.AddBooking"

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
			span.SetStatus(codes.Error, api.ErrBadRequest.Error())
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			api.WriteWithError(w, http.StatusUnauthorized, api.ErrNoAuth.Error())
			return
		}

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", userID)))

		req := &api.AddBookingRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, api.ErrBadRequest.Error())
				log.Error("some of the required values were not received", sl.Err(validateErr))
				api.WriteValidationError(w, validateErr)
				return
			}

			span.RecordError(api.ErrBadRequest)
			span.SetStatus(codes.Error, api.ErrBadRequest.Error())
			log.Error("failed to decode request body", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		span.AddEvent("request body decoded")
		log.Info("request body decoded", slog.Any("req", req))
		//TODO: getters
		mod, err := convert.ToBookingInfo(&api.Booking{
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

		bookingID, err := i.booking.AddBooking(ctx, mod)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("booking created", trace.WithAttributes(attribute.String("id", bookingID.String())))
		log.Info("booking added", slog.Any("id: ", bookingID))

		render.Status(r, http.StatusCreated)
		api.WriteWithStatus(w, http.StatusOK, api.AddBookingResponse{
			BookingID: bookingID,
		})
	}

}
