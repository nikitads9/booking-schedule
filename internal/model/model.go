package model

import (
	"database/sql"
	"time"
)

type EventInfo struct {
	// номер апаратаментов
	SuiteID int64 `json:"suiteID"`
	//Дата и время начала бронировании
	StartDate time.Time `json:"startDate"`
	// Дата и время окончания бронировании
	EndDate time.Time `json:"endDate"`
	// Интервал времени для уведомления о бронировании
	NotificationInterval time.Duration `json:"notificationInterval"`
	// telegram ID покупателя
	OwnerId string
}

type Event struct {
	//уникальный идентификатор бронирования
	Uuid string `json:"uuid"`
	//инфо о бронировании
	EventInfo *EventInfo `json:"eventInfo"`
	//датаи время создания
	CreatedAt time.Time `json:"createdAt"`
	//дата и время обновления
	UpdatedAt sql.NullTime `json:"updatedAt"`
}
