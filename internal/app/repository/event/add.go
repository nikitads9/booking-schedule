package event

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	t "event-schedule/internal/app/repository/table"
	"event-schedule/internal/pkg/db"
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

	var builder sq.InsertBuilder

	newID, err := uuid.NewV4()
	if err != nil {
		r.log.Error("failed to generate uuid", err)
		return uuid.Nil, ErrUuid
	}

	if mod.NotifyAt.Valid {
		builder = sq.Insert(t.EventTable).
			Columns(t.ID, t.OwnerID, t.SuiteID, t.StartDate, t.EndDate, t.CreatedAt, t.NotifyAt).
			Values(newID, mod.UserID, mod.SuiteID, mod.StartDate, mod.EndDate, time.Now(), mod.NotifyAt.Time)
	} else {
		builder = sq.Insert(t.EventTable).
			Columns(t.ID, t.OwnerID, t.SuiteID, t.StartDate, t.EndDate, t.CreatedAt).
			Values(newID, mod.UserID, mod.SuiteID, mod.StartDate, mod.EndDate, time.Now())
	}

	query, args, err := builder.Suffix("returning id").
		PlaceholderFormat(sq.Dollar).ToSql()
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
