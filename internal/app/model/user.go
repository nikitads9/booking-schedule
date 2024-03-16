package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type User struct {
	ID         int64      `db:"id"`
	TelegramID int64      `db:"telegram_id"`
	Nickname   string     `db:"telegram_nickname"`
	Name       string     `db:"name"`
	Password   string     `db:"password"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

type UpdateUserInfo struct {
	ID         int64       `db:"id"`
	TelegramID null.Int    `db:"telegram_id"`
	Nickname   null.String `db:"telegram_nickname"`
	Name       null.String `db:"name"`
	Password   null.String `db:"password"`
}
