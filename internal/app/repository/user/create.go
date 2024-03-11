package user

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	t "booking-schedule/internal/app/repository/table"

	"github.com/go-chi/chi/middleware"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	const op = "users.repository.CreateUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Insert(t.UserTable).
		Columns(t.TelegramID, t.TelegramNickname, t.Name, t.Password, t.CreatedAt).
		Values(user.TelegramID, user.Nickname, user.Name, *user.Password, time.Now())

	query, args, err := builder.PlaceholderFormat(sq.Dollar).Suffix("returning id").ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return 0, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	row, err := r.client.DB().QueryContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return 0, ErrNoConnection
		}
		if errors.As(err, &ErrDuplicate) {
			log.Error("this user already exists", err)
			return 0, ErrAlreadyExists
		}
		log.Error("query execution error", err)
		return 0, ErrQuery
	}
	defer row.Close()

	var id int64
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		log.Error("failed to scan returning id", err)
		return 0, ErrPgxScan
	}

	return id, nil
}
