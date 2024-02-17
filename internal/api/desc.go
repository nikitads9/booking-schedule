package api

import (
	"event-schedule/internal/model"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v3"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

type DeleteEventResponse struct {
	Response *Response `json:"response"`
}

func DeleteEventResponseAPI() *DeleteEventResponse {
	return &DeleteEventResponse{
		Response: OK(),
	}
}

func (rd *DeleteEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type AddEventRequest struct {
	// номер апаратаментов
	SuiteID int64 `json:"suiteID" validate:"required" example:"123"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00-03:00"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00-03:00"`
	// Интервал времени для предварительного уведомления о бронировании
	NotificationPeriod null.String `json:"notificationPeriod,omitempty" swaggertype:"primitive,string" example:"24h"`
}

type AddEventResponse struct {
	Response *Response `json:"response"`
	EventID  uuid.UUID `json:"eventID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

func AddEventResponseAPI(eventID uuid.UUID) *AddEventResponse {
	resp := &AddEventResponse{
		Response: OK(),
		EventID:  eventID,
	}

	return resp
}

func (a *AddEventRequest) Bind(r *http.Request) error {

	// Создаем объект валидатора
	// и передаем в него структуру, которую нужно провалидировать
	err := validator.New().Struct(a)
	if err != nil {
		return err
	}

	if a.StartDate.UTC().Before(time.Now().UTC()) || a.EndDate.UTC().Before(time.Now().UTC()) {
		return ErrExpiredDate
	}

	if a.EndDate.UTC().Sub(a.StartDate.UTC()) <= 0 {
		return ErrInvalidInterval
	}

	return nil
}

func (rd *AddEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetVacantDatesResponse struct {
	Response  *Response
	Intervals []*model.Interval `json:"intervals"`
}

func GetVacantDatesAPI(intervals []*model.Interval) *GetVacantDatesResponse {
	return &GetVacantDatesResponse{
		Response:  OK(),
		Intervals: intervals,
	}
}

func (rd *GetVacantDatesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetEventResponse struct {
	Response  *Response        `json:"response"`
	EventInfo *model.EventInfo `json:"event"`
}

func GetEventResponseAPI(event *model.EventInfo) *GetEventResponse {
	resp := &GetEventResponse{
		Response:  OK(),
		EventInfo: event,
	}

	return resp
}

func (rd *GetEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetEventsResponse struct {
	Response   *Response          `json:"response"`
	EventsInfo []*model.EventInfo `json:"events"`
}

func GetEventsResponseAPI(events []*model.EventInfo) *GetEventsResponse {
	resp := &GetEventsResponse{
		Response:   OK(),
		EventsInfo: events,
	}

	return resp
}

func (rd *GetEventsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetVacantRoomsResponse struct {
	Response *Response      `json:"response"`
	Rooms    []*model.Suite `json:"rooms"`
}

func GetVacantRoomsAPI(rooms []*model.Suite) *GetVacantRoomsResponse {
	return &GetVacantRoomsResponse{
		Response: OK(),
		Rooms:    rooms,
	}
}

func (rd *GetVacantRoomsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UpdateEventRequest struct {
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" validate:"required" format:"int64" example:"123"`
	// Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00-03:00"`
	// Дата и время окончания бронирования
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00-03:00"`
	// Интервал времени для уведомления о бронировании
	NotificationPeriod null.String `json:"notificationPeriod,omitempty" swaggertype:"primitive,string" example:"24h"`
}

type UpdateEventResponse struct {
	Response *Response `json:"response"`
}

// Функцияя для проверки поступающего запроса, чтобы удостовериться, что
// окончание бронирования не  происходит до его начала, чтобы даты не были истекшими.
// Также при изменении одной даты нужно изменять и другую, а также указывать, за сколько оповестить о бронировании.
func (u *UpdateEventRequest) Bind(r *http.Request) error {
	// Создаем объект валидатора
	// и передаем в него структуру, которую нужно провалидировать
	err := validator.New().Struct(u)
	if err != nil {
		return err
	}

	// Проверка, что обе даты еще не прошли
	if (u.StartDate.UTC().Before(time.Now().UTC())) || (u.EndDate.UTC().Before(time.Now().UTC())) {
		return ErrExpiredDate
	}

	//проверка, что дата окончания не находится перед датой начала и не совпадает с ней
	if u.EndDate.UTC().Sub(u.StartDate.UTC()) <= 0 {
		return ErrInvalidInterval
	}

	return nil
}

func (rd *UpdateEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UpdateEventResponseAPI() *UpdateEventResponse {
	return &UpdateEventResponse{
		Response: OK(),
	}
}
