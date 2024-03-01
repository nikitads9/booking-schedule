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

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Delete(t.EventTable).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     "events.repository.DeleteEvent",
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
