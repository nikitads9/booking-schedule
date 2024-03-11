package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
)

type BookingInfo struct {
	// Уникальный идентификатор бронирования
	ID uuid.UUID `db:"id"`
	// Номер апартаментов
	SuiteID int64 `db:"suite_id"`
	//Дата и время начала бронировании
	StartDate time.Time `db:"start_date"`
	// Дата и время окончания бронировании
	EndDate time.Time `db:"end_date"`
	// Интервал времени для уведомления о бронировании
	NotifyAt time.Duration `db:"notify_at"`
	// Дата и время создания
	CreatedAt time.Time `db:"created_at"`
	// Дата и время обновления
	UpdatedAt null.Time `db:"updated_at"`
	// Идентификатор владельца бронирования
	UserID int64 `db:"user_id"`
}

type Interval struct {
	StartDate time.Time `db:"start"`
	EndDate   time.Time `db:"end"`
}

type Suite struct {
	SuiteID  int64  `db:"suite_id"`
	Capacity int8   `db:"capacity"`
	Name     string `db:"name"`
}

type Availibility struct {
	Availible        bool `db:"availible"`
	OccupiedByClient bool `db:"occupied_by_client"`
}
