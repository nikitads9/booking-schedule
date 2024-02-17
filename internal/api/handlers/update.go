package handlers

import (
	"errors"
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
)

// UpdateEvent godoc
//
//	@Summary		Updates event info
//	@Description	Updates an existing event with given EventID, suiteID, startDate, endDate values (notificationPeriod being optional). Implemented with the use of transaction: first room availibility is checked. In case one attempts to alter his previous booking (i.e. widen or tighten its' limits) the booking is updated.  If it overlaps with smb else's booking the request is considered unsuccessful. startDate parameter  is to be before endDate and both should not be expired.
//	@ID				modifyEventByJSON
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param			event_id path	string	true	"event_id"	Format(uuid) default(550e8400-e29b-41d4-a716-446655440000)
//	@Param          event body		api.UpdateEventRequest	true	"UpdateEventRequest"
//	@Success		200	{object}	api.UpdateEventResponse
//	@Failure		400	{object}	api.UpdateEventResponse
//	@Failure		404	{object}	api.UpdateEventResponse
//	@Failure		422	{object}	api.UpdateEventResponse
//	@Failure		503	{object}	api.UpdateEventResponse
//	@Router			/{user_id}/{event_id}/update [patch]
func (i *Implementation) UpdateEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.UpdateEvent"

		ctx := r.Context()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		req := &api.UpdateEventRequest{}
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

		mod, err := convert.ToUpdateEventInfo(r, req)
		if err != nil {
			log.Error("invalid request", sl.Err(err)) //TODO: log real error
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		err = i.Service.UpdateEvent(ctx, mod)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("event updated", slog.Any("id:", mod.EventID))

		render.Render(w, r, api.UpdateEventResponseAPI())
	}
}
