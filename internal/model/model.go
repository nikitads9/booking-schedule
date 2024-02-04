package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	// Идентификатор пользователя
	UserID int64
	// Номер апартаментов
	SuiteID int64
	// Дата и время начала бронировании
	StartDate time.Time
	// Дата и время окончания бронировании
	EndDate time.Time
	// Интервал времени для уведомления о бронировании
	NotificationPeriod string
}

type EventInfo struct {
	// Уникальный идентификатор бронирования
	EventID uuid.UUID `json:"EventID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid" db:"id"`
	// Номер апартаментов
	SuiteID int64 `json:"suiteID" db:"suite_id"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate" db:"start_date"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate" db:"end_date"`
	// Интервал времени для уведомления о бронировании
	NotificationPeriod string `json:"notificationPeriod" db:"notification_period"`
	// Дата и время создания
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	// Дата и время обновления
	UpdatedAt sql.NullTime `json:"updatedAt" db:"updated_at"`
}

// Содержимое Update/Modify запроса
type UpdateEventInfo struct {
	EventID uuid.UUID `json:"eventID" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	// Идентификатор пользователя в системе
	UserID int64
	// Номер апартаментов
	SuiteID sql.NullInt64 `json:"suiteID" example:"123"`
	// Дата и время начала бронировании
	StartDate sql.NullTime `json:"startDate" example:"2006-01-02T15:04:05-07:00"`
	// Дата и время окончания бронировании
	EndDate sql.NullTime `json:"endDate" example:"2006-01-02T15:04:05-07:00"`
	// Интервал времени для уведомления о бронировании
	NotificationPeriod sql.NullString `json:"notificationPeriod" example:"24h"`
}

type GetEventsInfo struct {
	// Уникальный идентификатор пользователя в системе
	UserID int64 `in:"path=user"`
	// Начало интервала поиска
	StartDate time.Time `in:"query=start"`
	// Конец интервала поиска
	EndDate time.Time `in:"query=end"`
}

type EventInfoDB struct {
	// номер апартаментов
	SuiteID int64 `db:"suite_id"`
	// Название номера
	SuiteName string `db:"name"`
	// Дата и время начала бронировании
	StartDate time.Time `db:"start_date"`
	// Дата и время окончания бронировании
	EndDate time.Time `db:"end_date"`
	// Интервал времени для уведомления о бронировании
	NotificationPeriod string `db:"notification_period"`
}

type Suite struct {
	SuiteID  int64  `json:"suiteID" example:"123"`
	Capacity int8   `json:"capacity" example:"4"`
	Name     string `json:"name" example:"Winston Churchill"`
}

type Interval struct {
	StartDate time.Time `json:"startDate" example:"2006-01-02T15:04:05-07:00"`
	EndDate   time.Time `json:"endDate" exaple:"2006-01-02T15:04:05-07:00"`
}
