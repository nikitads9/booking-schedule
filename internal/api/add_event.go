package api

import (
	"encoding/json"
	"errors"
	"event-schedule/internal/lib/logger/sl"
	"log/slog"

	"net/http"
	"time"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type AddEventRequest struct {
	// номер апаратаментов
	SuiteID int64 `json:"suiteID" validate:"required"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required"`
	// Интервал времени для предварительного уведомления о бронировании
	NotificationInterval time.Duration `json:"notificationInterval"`
	// telegram ID покупателя
	OwnerID string `json:"ownerID"`
}

type AddEventResponse struct {
	Response *Response `json:"response"`
	EventID  string    `json:"eventID"`
}

func (i *Implementation) AddEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.AddEvent"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := &AddEventRequest{}
		if err := render.Bind(r, req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Render(w, r, ErrValidationError(validateErr))

			return
		}

		id, err := i.Service.AddEvent(r.Context())
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
		}
		res, _ := json.Marshal(req)
		render.Status(r, http.StatusCreated)
		render.Render(w, r, AddEventResponseAPI("received AddEvent with THIS body:"+string(res)+id))
	}

}

func AddEventResponseAPI(eventID string) *AddEventResponse {
	resp := &AddEventResponse{
		Response: OK(),
		EventID:  eventID,
	}

	return resp
}

func (e *AddEventRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if e.OwnerID == "" {
		return errors.New("missing required event fields")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or further nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	//a.ProtectedID = ""                                 // unset the protected ID
	//e.Event = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

func (rd *AddEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	//rd.Elapsed = 10
	return nil
}

//код для того чтобы зафиксировать объект
/* if userID = chi.URLParam(r, "userID"); userID == "" {
	//event, err = db.AddEvent(userID)
} else {
	render.Render(w, r, ErrNotFound)
	return
}
if err != nil {
	render.Render(w, r, ErrNotFound)
	return
} */

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
