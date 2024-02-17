package event

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	t "event-schedule/internal/repository/table"
	"log/slog"
	"time"

	"github.com/go-chi/chi/middleware"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

func (r *repository) AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error) {
	const op = "events.repository.AddEvent"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Insert(t.EventTable).
		Columns(t.OwnerID, t.SuiteID, t.StartDate, t.EndDate, t.CreatedAt).
		Values(mod.UserID, mod.SuiteID, mod.StartDate, mod.EndDate, time.Now()).
		Suffix("returning id").
		PlaceholderFormat(sq.Dollar)

	if mod.NotifyAt.Valid {
		builder = builder.Columns("notify_at", t.NotifyAt).
			Values(mod.NotifyAt.Time)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return uuid.Nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	row, err := r.client.DB().QueryContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return uuid.Nil, ErrNoConnection
		}
		r.log.Error("query execution error", err)
		return uuid.Nil, ErrQuery
	}

	var id uuid.UUID
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		r.log.Error("failed to scan pgx row", err)
		return uuid.Nil, ErrPgxScan

	}

	return id, nil
}
