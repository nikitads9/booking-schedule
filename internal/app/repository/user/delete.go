package user

import (
	"context"
	"errors"
	"event-schedule/internal/pkg/db"
	"log/slog"

	t "event-schedule/internal/app/repository/table"

	"github.com/go-chi/chi/middleware"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) DeleteUser(ctx context.Context, userID int64) error {
	const op = "users.repository.DeleteUser"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Delete(t.UserTable).
		Where(sq.Eq{t.ID: userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
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
		log.Error("query execution error", err)
		return ErrQuery
	}

	if result.RowsAffected() == 0 {
		log.Error("unsuccessful delete", ErrNoRowsAffected)
		return ErrNotFound
	}

	return nil

}
