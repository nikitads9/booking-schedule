package model

import "time"

type User struct {
	ID         int64      `db:"id"`
	TelegramID int64      `db:"telegram_id"`
	Nickname   string     `db:"telegram_nickname"`
	Name       string     `db:"name"`
	Password   *string    `db:"password"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}
