package handlers

import (
	"event-schedule/internal/api"
	"event-schedule/internal/convert"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"
	"strconv"

	"net/http"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// AddEvent godoc
//
//	@Summary		Adds event
//	@Description	Adds an even with given parameters associated with user. NotificationPeriod must look like {number}s,{number}m or {number}h.
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
		const op = "handlers.events.api.AddEvent"

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
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		err = validator.New().Struct(req)
		if err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Render(w, r, api.ErrValidationError(validateErr))

			return
		}

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			log.Error("invalid request", sl.Err(api.ErrNoUserID))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoUserID))
			return
		}

		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		eventID, err := i.Service.AddEvent(ctx, convert.ToEvent(req, id))
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("event added", slog.Any("id", eventID))

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
