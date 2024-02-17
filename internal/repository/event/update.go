package event

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	t "event-schedule/internal/repository/table"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error {
	const op = "events.repository.UpdateEvent"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Update(t.EventTable).
		Set(t.UpdatedAt, time.Now().UTC()).
		Set("end_date", mod.EndDate).
		Set("start_date", mod.StartDate).
		Set("suite_id", mod.SuiteID).
		Where(sq.Eq{"id": mod.EventID}).
		PlaceholderFormat(sq.Dollar)

	if mod.NotifyAt.Valid {
		builder = builder.Set("notify_at", mod.NotifyAt.Time)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			r.log.Error("event with this id not found", err)
			return ErrNotFound
		}
		r.log.Error("query execution error", err)
		return ErrQuery
	}

	return nil
}
