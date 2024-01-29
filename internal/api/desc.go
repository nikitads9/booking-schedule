package api

import (
	"database/sql"
	"event-schedule/internal/model"
	"net/http"
	"time"

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
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-01-28T17:43:00Z03:00"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-01-29T17:43:00Z03:00"`
	// Интервал времени для предварительного уведомления о бронировании
	NotificationPeriod string `json:"notificationPeriod" example:"1h"`
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

func (e *AddEventRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	/* 	if e.OwnerID == "" {
		return errors.New("missing required event fields")
	} */

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
	Response *Response `json:"response"`
	//TODO: implement convert Event structs
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
	SuiteID sql.NullInt64 `json:"suiteID,omitempty" swaggertype:"primitive,integer" format:"int64" example:"123"`
	// Дата и время начала бронировании
	StartDate sql.NullTime `json:"startDate,omitempty" swaggertype:"primitive,string" example:"2006-01-02T15:04:05Z07:00"`
	// Дата и время окончания бронирования
	EndDate sql.NullTime `json:"endDate,omitempty" swaggertype:"primitive,string" example:"2006-01-02T15:04:05Z07:00"`
	// Интервал времени для уведомления о бронировании
	NotificationPeriod sql.NullString `json:"notificationPeriod,omitempty" swaggertype:"primitive,string" example:"24h"`
}

type UpdateEventResponse struct {
	Response *Response `json:"response"`
}

func (e *UpdateEventRequest) Bind(r *http.Request) error {
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
