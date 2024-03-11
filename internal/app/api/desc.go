package api

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
)

type Event struct {
	EventID uuid.UUID
	// Идентификатор пользователя
	UserID int64
	// Номер апартаментов
	SuiteID int64
	// Дата и время начала бронировании
	StartDate time.Time
	// Дата и время окончания бронировании
	EndDate time.Time
	// Интервал времени для уведомления о бронировании
	NotifyAt null.String
}

type EventInfo struct {
	// Уникальный идентификатор бронирования
	ID uuid.UUID `json:"EventID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" example:"1"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" example:"2024-03-28T17:43:00Z"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" example:"2024-03-29T17:43:00Z"`
	// Интервал времени для уведомления о бронировании
	NotifyAt string `json:"notifyAt,omitempty" example:"24h00m00s"`
	// Дата и время создания
	CreatedAt time.Time `json:"createdAt" example:"2024-03-27T17:43:00Z"`
	// Дата и время обновления
	UpdatedAt time.Time `json:"updatedAt,omitempty" example:"2024-03-27T18:43:00Z"`
	// Идентификатор владельца бронирования
	UserID int64 `json:"userID,omitempty" example:"1"`
} //@name EventInfo

type AddEventRequest struct {
	// Номер апаратаментов
	SuiteID int64 `json:"suiteID" validate:"required" example:"1"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00Z"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00Z"`
	// Интервал времени для предварительного уведомления о бронировании
	NotifyAt null.String `json:"notifyAt,omitempty" swaggertype:"primitive,string" example:"24h"`
} //@name AddEventRequest

type AddEventResponse struct {
	Response *Response `json:"response"`
	EventID  uuid.UUID `json:"eventID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
} //@name AddEventResponse

func AddEventResponseAPI(eventID uuid.UUID) *AddEventResponse {
	resp := &AddEventResponse{
		Response: OK(),
		EventID:  eventID,
	}

	return resp
}

func CheckDates(start time.Time, end time.Time) error {
	if start.Before(time.Now()) || end.Before(time.Now()) {
		return ErrExpiredDate
	}

	if end.Sub(start) <= 0 {
		return ErrInvalidInterval
	}

	return nil
}

func (arq *AddEventRequest) Bind(req *http.Request) error {
	err := validator.New().Struct(arq)
	if err != nil {
		return err
	}

	return CheckDates(arq.StartDate, arq.EndDate)
}

func (ar *AddEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetEventResponse struct {
	Response  *Response  `json:"response"`
	EventInfo *EventInfo `json:"event"`
} //@name GetEventResponse

func GetEventResponseAPI(event *EventInfo) *GetEventResponse {
	resp := &GetEventResponse{
		Response:  OK(),
		EventInfo: event,
	}

	return resp
}

func (ge *GetEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type GetEventsResponse struct {
	Response   *Response    `json:"response"`
	EventsInfo []*EventInfo `json:"events"`
} //@name GetEventsResponse

func GetEventsResponseAPI(events []*EventInfo) *GetEventsResponse {
	resp := &GetEventsResponse{
		Response:   OK(),
		EventsInfo: events,
	}

	return resp
}

func (ges *GetEventsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Interval struct {
	// Номер свободен с
	StartDate time.Time `json:"start" example:"2024-03-10T15:04:05Z"`
	// Номер свободен по
	EndDate time.Time `json:"end" example:"2024-04-10T15:04:05Z"`
} //@name Interval

type GetVacantDatesResponse struct {
	Response  *Response
	Intervals []*Interval `json:"intervals"`
} //@name GetVacantDateResponse

func GetVacantDatesAPI(intervals []*Interval) *GetVacantDatesResponse {
	return &GetVacantDatesResponse{
		Response:  OK(),
		Intervals: intervals,
	}
}

func (gd *GetVacantDatesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Suite struct {
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" example:"1"`
	// Вместимость в персонах
	Capacity int8 `json:"capacity" example:"4"`
	// Название апартаментов
	Name string `json:"name" example:"Winston Churchill"`
} //@name Suite

type GetVacantRoomsResponse struct {
	Response *Response `json:"response"`
	Rooms    []*Suite  `json:"rooms"`
} //@name GetVacantRoomsResponse

func GetVacantRoomsAPI(rooms []*Suite) *GetVacantRoomsResponse {
	return &GetVacantRoomsResponse{
		Response: OK(),
		Rooms:    rooms,
	}
}

func (gr *GetVacantRoomsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UpdateEventRequest struct {
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" validate:"required" format:"int64" example:"1"`
	// Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00-03:00"`
	// Дата и время окончания бронирования
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00-03:00"`
	// Интервал времени для уведомления о бронировании
	NotifyAt null.String `json:"notifyAt,omitempty" swaggertype:"primitive,string" example:"24h"`
} //@name UpdateEventRequest

func (urq *UpdateEventRequest) Bind(r *http.Request) error {
	err := validator.New().Struct(urq)
	if err != nil {
		return err
	}

	return CheckDates(urq.StartDate, urq.EndDate)
}

type UpdateEventResponse struct {
	Response *Response `json:"response"`
} //@name UpdateEventResponse

func (ur *UpdateEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UpdateEventResponseAPI() *UpdateEventResponse {
	return &UpdateEventResponse{
		Response: OK(),
	}
}

type DeleteEventResponse struct {
	Response *Response `json:"response"`
} //@name DeleteEventResponse

func DeleteEventResponseAPI() *DeleteEventResponse {
	return &DeleteEventResponse{
		Response: OK(),
	}
}

func (dr *DeleteEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type AuthResponse struct {
	Response *Response `json:"response"`
	// JWT токен для доступа
	Token string `json:"token"`
} //@name SignInResponse

func AuthResponseAPI(token string) *AuthResponse {
	return &AuthResponse{
		Response: OK(),
		Token:    token,
	}
}

func (auth *AuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type User struct {
	// Телеграм ID пользователя
	TelegramID int64 `json:"telegramID" validate:"required,notblank" example:"1235678"`
	// Никнейм пользователя в телеграме
	Nickname string `json:"telegramNickname" validate:"required,notblank" example:"pavel_durov"`
	// Имя пользователя
	Name string `json:"name" validate:"required,notblank" example:"Pavel Durov"`
	// Пароль
	Password string `json:"password" validate:"required,notblank" example:"12345"`
} //@name User

func (u *User) Bind(req *http.Request) error {
	v := validator.New()
	err := v.RegisterValidation("notblank", NotBlank)
	if err != nil {
		return err
	}

	err = v.Struct(u)
	if err != nil {
		return err
	}

	return nil
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Int64:
		return !field.IsZero()
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
