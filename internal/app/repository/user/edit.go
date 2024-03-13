package user

import (
	"booking-schedule/internal/app/model"
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) EditUser(ctx context.Context, user *model.UpdateUserInfo) error {
	const op = "bookings.repository.EditUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Update(t.UserTable).
		Set(t.UpdatedAt, time.Now()).
		Where(sq.Eq{t.ID: user.ID})

	if user.Name.Valid {
		builder = builder.Set(t.Name, user.Name.String)
	}

	if user.Nickname.Valid {
		builder = builder.Set(t.TelegramNickname, user.Nickname.String)
	}

	if user.Password.Valid {
		builder = builder.Set(t.Password, user.Password.String)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	result, err := r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		if errors.As(err, &ErrDuplicate) {
			log.Error("this user already exists", err)
			return ErrAlreadyExists
		}
		log.Error("query execution error", err)
		return ErrQuery
	}

	if result.RowsAffected() == 0 {
		log.Error("unsuccessful update", ErrNoRowsAffected)
		return ErrNotFound
	}

	return nil
}
