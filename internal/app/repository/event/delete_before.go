package event

import (
	"context"
	"errors"
	t "event-schedule/internal/app/repository/table"
	"event-schedule/internal/pkg/db"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) DeleteEventsBeforeDate(ctx context.Context, date time.Time) error {
	const op = "events.repository.DeleteEventsBeforeDate"

	log := r.log.With(
		slog.String("op", op),
	)

	builder := sq.Delete(t.EventTable).
		Where(sq.Lt{t.EndDate: date}).
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

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		log.Error("query execution error", err)
		return ErrQuery
	}

	return nil
}
