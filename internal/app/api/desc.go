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

type Booking struct {
	BookingID uuid.UUID
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

type AddBookingRequest struct {
	// Номер апаратаментов
	SuiteID int64 `json:"suiteID" validate:"required" example:"1"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00Z"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00Z"`
	// Интервал времени для предварительного уведомления о бронировании
	NotifyAt null.String `json:"notifyAt,omitempty" swaggertype:"primitive,string" example:"24h"`
} //@name AddBookingRequest

type AddBookingResponse struct {
	BookingID uuid.UUID `json:"bookingID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
} //@name AddBookingResponse

type BookingInfo struct {
	// Уникальный идентификатор бронирования
	ID uuid.UUID `json:"BookingID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" example:"1"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" example:"2024-03-28T17:43:00Z"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" example:"2024-03-29T17:43:00Z"`
	// Интервал времени для уведомления о бронировании
	NotifyAt *string `json:"notifyAt,omitempty" example:"24h00m00s"`
	// Дата и время создания
	CreatedAt time.Time `json:"createdAt" example:"2024-03-27T17:43:00Z"`
	// Дата и время обновления
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:"2024-03-27T18:43:00Z"`
	// Идентификатор владельца бронирования
	UserID int64 `json:"userID,omitempty" example:"1"`
} //@name BookingInfo

type GetBookingResponse struct {
	BookingInfo *BookingInfo `json:"booking"`
} //@name GetBookingResponse

type GetBookingsResponse struct {
	BookingsInfo []*BookingInfo `json:"bookings"`
} //@name GetBookingsResponse

type UpdateBookingRequest struct {
	// Номер апаратаментов
	SuiteID int64 `json:"suiteID" validate:"required" example:"1"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" validate:"required" example:"2024-03-28T17:43:00Z"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" validate:"required" example:"2024-03-29T17:43:00Z"`
	// Интервал времени для предварительного уведомления о бронировании
	NotifyAt null.String `json:"notifyAt,omitempty" swaggertype:"primitive,string" example:"24h"`
} //@name UpdateBookingRequest

type Interval struct {
	// Номер свободен с
	StartDate time.Time `json:"start" example:"2024-03-10T15:04:05Z"`
	// Номер свободен по
	EndDate time.Time `json:"end" example:"2024-04-10T15:04:05Z"`
} //@name Interval

type GetVacantDatesResponse struct {
	Intervals []*Interval `json:"intervals"`
} //@name GetVacantDateResponse

type Suite struct {
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" example:"1"`
	// Вместимость в персонах
	Capacity int8 `json:"capacity" example:"4"`
	// Название апартаментов
	Name string `json:"name" example:"Winston Churchill"`
} //@name Suite

type GetVacantRoomsResponse struct {
	Rooms []*Suite `json:"rooms"`
} //@name GetVacantRoomsResponse

type AuthResponse struct {
	// JWT токен для доступа
	Token string `json:"token"`
} //@name AuthResponse

type SignUpRequest struct {
	// Телеграм ID пользователя
	TelegramID int64 `json:"telegramID" validate:"required,notblank" example:"1235678"`
	// Никнейм пользователя в телеграме
	Nickname string `json:"telegramNickname" validate:"required,notblank" example:"pavel_durov"`
	// Имя пользователя
	Name string `json:"name" validate:"required,notblank" example:"Pavel Durov"`
	// Пароль
	Password string `json:"password" validate:"required,notblank" example:"12345"`
} //@name SignUpRequest

type UserInfo struct {
	//ID пользователя в системе
	ID int64 `json:"id"`
	// Телеграм ID пользователя
	TelegramID int64 `json:"telegramID"`
	// Никнейм пользователя в телеграме
	Nickname string `json:"telegramNickname"`
	// Имя пользователя
	Name string `json:"name"`
	// Дата и время регистрации
	CreatedAt time.Time `json:"createdAt"`
	// Дата и время обновления профиля
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
} //@name UserInfo

type GetMyProfileResponse struct {
	// Профиль пользователя
	Profile *UserInfo `json:"profile"`
} //@name GetMyProfileResponse

type EditMyProfileRequest struct {
	// Имя пользователя
	Name null.String `json:"name" swaggertype:"primitive,string" validate:"notblank" example:"Kolya Durov"`
	// Телеграм ID пользователя
	TelegramID null.Int `json:"telegramID" swaggertype:"primitive,integer" validate:"notblank" example:"1235678"`
	// Никнейм пользователя в телеграме
	Nickname null.String `json:"telegramNickname" swaggertype:"primitive,string" validate:"notblank" example:"kolya_durov"`
	// Пароль
	Password null.String `json:"password" swaggertype:"primitive,string" validate:"notblank" example:"123456"`
} // @name EditMyProfileRequest

func (arq *AddBookingRequest) Bind(req *http.Request) error {
	err := validator.New().Struct(arq)
	if err != nil {
		return err
	}

	return CheckDates(arq.StartDate, arq.EndDate)
}

func (urq *UpdateBookingRequest) Bind(req *http.Request) error {
	err := validator.New().Struct(urq)
	if err != nil {
		return err
	}

	return CheckDates(urq.StartDate, urq.EndDate)
}

func (srq *SignUpRequest) Bind(req *http.Request) error {
	v := validator.New()
	err := v.RegisterValidation("notblank", NotBlank)
	if err != nil {
		return err
	}

	err = v.Struct(srq)
	if err != nil {
		return err
	}

	return nil
}

func (empr *EditMyProfileRequest) Bind(req *http.Request) error {
	v := validator.New()
	err := v.RegisterValidation("notblank", NotBlank)
	if err != nil {
		return err
	}

	err = v.Struct(empr)
	if err != nil {
		return err
	}

	if (empr.TelegramID.Valid && !empr.Nickname.Valid) || (!empr.TelegramID.Valid && empr.Nickname.Valid) {
		return ErrIncompleteRequest
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

func CheckDates(start time.Time, end time.Time) error {
	if start.Before(time.Now()) || end.Before(time.Now()) {
		return ErrExpiredDate
	}

	if end.Sub(start) <= 0 {
		return ErrInvalidInterval
	}

	return nil
}
