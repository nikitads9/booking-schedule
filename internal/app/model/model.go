package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
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
	NotifyAt null.Time
}

func (e *Event) GetEventID() uuid.UUID {
	if e != nil {
		return e.EventID
	}

	return uuid.Nil
}

func (e *Event) SetEventID(id uuid.UUID) {
	if e != nil {
		e.EventID = id
	}
}

func (e *Event) GetNotifyAt() null.Time {
	if e != nil {
		return e.NotifyAt
	}
	return null.Time{}
}

func (e *Event) SetNotifyAt(t time.Time) {
	if e != nil {
		e.NotifyAt = null.Time{
			Time:  t,
			Valid: true,
		}
	}
}

/* // Содержимое Update/Modify запроса
type UpdateEventInfo struct {
	EventID uuid.UUID
	// Идентификатор пользователя в системе
	UserID int64
	// Номер апартаментов
	SuiteID int64
	// Дата и время начала бронировании
	StartDate time.Time
	// Дата и время окончания бронировании
	EndDate time.Time
	// Интервал времени для уведомления о бронировании
	NotifyAt null.Time
} */

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
	NotifyAt time.Time `json:"notifyAt,omitempty" db:"notify_at"`
	// Дата и время создания
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	// Дата и время обновления
	UpdatedAt null.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

type GetEventsInfo struct {
	// Уникальный идентификатор пользователя в системе
	UserID int64 `in:"path=user"`
	// Начало интервала поиска
	StartDate time.Time `in:"query=start"`
	// Конец интервала поиска
	EndDate time.Time `in:"query=end"`
}

type Interval struct {
	StartDate time.Time `json:"start" db:"start" example:"2024-03-02T15:04:05-07:00"`
	EndDate   time.Time `json:"end" db:"end" exaple:"2024-04-02T15:04:05-07:00"`
}

type Suite struct {
	SuiteID  int64  `json:"suiteID" db:"suite_id" example:"123"`
	Capacity int8   `json:"capacity" db:"capacity" example:"4"`
	Name     string `json:"name" db:"name" example:"Winston Churchill"`
}

type Availibility struct {
	Availible        bool `db:"availible"`
	OccupiedByClient bool `db:"occupied_by_client"`
}