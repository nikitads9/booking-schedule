package user

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"

	t "booking-schedule/internal/app/repository/table"

	"github.com/go-chi/chi/middleware"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	const op = "users.repository.GetUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Select("*").
		From(t.UserTable).
		Where(sq.Eq{t.ID: userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res = new(model.User)
	err = r.client.DB().GetContext(ctx, res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("booking with this id not found", err)
			return nil, ErrNotFound
		}
		log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
