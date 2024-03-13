package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64      `db:"id"`
	TelegramID int64      `db:"telegram_id"`
	Nickname   string     `db:"telegram_nickname"`
	Name       string     `db:"name"`
	Password   *string    `db:"password"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

type UpdateUserInfo struct {
	ID       int64          `db:"id"`
	Nickname sql.NullString `db:"telegram_nickname"`
	Name     sql.NullString `db:"name"`
	Password sql.NullString `db:"password"`
}
