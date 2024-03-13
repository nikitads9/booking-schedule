package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v3"
)

type BookingInfo struct {
	ID        uuid.UUID     `db:"id"`
	SuiteID   int64         `db:"suite_id"`
	StartDate time.Time     `db:"start_date"`
	EndDate   time.Time     `db:"end_date"`
	NotifyAt  time.Duration `db:"notify_at"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt null.Time     `db:"updated_at"`
	UserID    int64         `db:"user_id"`
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
