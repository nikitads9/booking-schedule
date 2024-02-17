package handlers

import (
	"errors"
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"

	"net/http"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// AddEvent godoc
//
//	@Summary		Adds event
//	@Description	Adds an even with given parameters associated with user. NotificationPeriod must look like {number}s,{number}m or {number}h. Implemented with the use of transaction: first the availibility is checked. In case one's new booking request intersects with and old one(even if belongs to him), the request is considered erratic.
//	@ID				addByEventJSON
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path	int	true	"user_id"	Format(int64) default(1234)
//	@Param          event	body	api.AddEventRequest	true	"AddEventRequest"
//	@Success		200	{object}	api.AddEventResponse
//	@Failure		400	{object}	api.AddEventResponse
//	@Failure		404	{object}	api.AddEventResponse
//	@Failure		422	{object}	api.AddEventResponse
//	@Failure		503	{object}	api.AddEventResponse
//	@Router			/{user_id}/add [post]
func (i *Implementation) AddEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.AddEvent"

		ctx := r.Context()

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		req := &api.AddEventRequest{}
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

		mod, err := convert.ToEvent(r, req)
		if err != nil {
			log.Error("invalid request", sl.Err(err)) //TODO: log real error
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		eventID, err := i.Service.AddEvent(ctx, mod)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("event added", slog.Any("id:", eventID))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.AddEventResponseAPI(eventID))
	}

}

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
