package event

import (
	"context"
	"errors"
	t "event-schedule/internal/app/repository/table"
	"event-schedule/internal/pkg/db"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
)

func (r *repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	const op = "events.repository.DeleteEvent"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Delete(t.EventTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	result, err := r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		r.log.Error("query execution error", err)
		return ErrQuery
	}

	if result.RowsAffected() == 0 {
		r.log.Error("unsuccessful delete", ErrNoRowsAffected)
		return ErrNotFound
	}

	return nil
}
